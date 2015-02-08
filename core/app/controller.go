package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"cloudvm/gocommon/interfaces"
	"cloudvm/solo/data"

	"github.com/jinzhu/gorm"
)

type Controller struct {
	Cls          ControllerInterface // interface{}
	ResHandler   ResourceHandlerInterface
	db           *gorm.DB
	api          interfaces.ApiInterface
	host         HostInterface
	tempDir      string
	dataCacheDir string
}

func NewController(context *data.VoxelContext) (*Controller, error) {
	ResHandler, _ := NewResourceHandler()
	host := Host{}
	controller := Controller{
		db:           context.Db,
		api:          context.ApiClient,
		host:         &host,
		tempDir:      path.Join(os.TempDir(), "vxlapps"),
		dataCacheDir: path.Join(context.DataRoot, "user_data"),
		ResHandler:   ResHandler,
	}
	controller.Cls = &controller

	// create tmpdir
	os.RemoveAll(controller.tempDir)
	os.MkdirAll(controller.tempDir, 0755)
	os.Chmod(controller.tempDir, 0755)

	// disable google play for now as some existing games are crashing
	log.Println("Disabling google play services...")
	host.runDeviceShell("pm", "disable", "com.google.android.gms")
	host.runDeviceShell("pm", "disable", "com.google.android.gsf")

	//db.LogMode(true)
	// create tables if not already there
	if err := controller.db.AutoMigrate(&App{}).Error; err != nil {
		return nil, err
	}
	if err := controller.db.AutoMigrate(&AppFile{}).Error; err != nil {
		return nil, err
	}

	return &controller, nil
}

func (c *Controller) Self() ControllerInterface {
	return c.Cls // .(ControllerInterface)
}

func (c *Controller) InstalledApps() []*App {
	var apps []*App
	c.db.Where("state = ?", Installed).Find(&apps)
	for _, app := range apps {
		app.persisted = true
		c.db.Where("app_id = ?", app.Id).Find(&app.Files)
	}
	return apps
}

func (c *Controller) Get(id int64) *App {
	var app App
	r := c.db.First(&app, id)
	if r.RecordNotFound() {
		return nil
	}
	app.persisted = true

	// Fetch the files too if any
	c.db.Where("app_id = ?", id).Find(&app.Files)
	return &app
}

func (c *Controller) Save(app *App) error {
	if app.persisted {
		c.db.Save(app)
	} else {
		c.db.Create(app)
		app.persisted = true
	}

	// remove old AppFiles first
	c.db.Where("app_id = ?", app.Id).Delete(AppFile{})

	// save the app files
	for _, appfile := range app.Files {
		c.db.Create(&appfile)
	}
	return nil
}

func (c *Controller) SizeUsed() int64 {
	total := int64(0)

	apps := c.Self().InstalledApps()
	for _, app := range apps {
		total += app.Size
		for _, appfile := range app.Files {
			total += appfile.Size
		}
	}

	return total
}

func (c *Controller) StartApp(id int64, userData *UserData, params AppParams) error {
	app := c.Self().Get(id)
	if app == nil || app.State == NotInstalled {
		return fmt.Errorf("App %d is not installed", id)
	}

	err := c.host.SetEnabled(app, true)
	if err != nil {
		return err
	}

	if len(params.AndroidId) > 0 {
		err = c.host.SetAndroidId(params.AndroidId)
		if err != nil {
			log.Println("Could not set android id")
			return nil
		}
	}

	if userData != nil {
		log.Println("restoring data")
		err = c.RestoreData(app, userData)
		if err != nil {
			return err
		}
	}

	// increase app last used date and used count
	app.LastUsed = time.Now()
	app.RunCount += 1
	c.Self().Save(app)

	return c.host.StartApp(app, params)
}

// TODO: This is called ResumeAppInstance instead of ResumeApp precisely because of the
// need to avoid do calls to "ps".  But this looks like a very hacky way of doing this
// so have to see how to refactor this interface
func (c *Controller) ResumeAppInstance(pid int) error {
	return c.host.ResumePid(pid)
}

func (c *Controller) StopApp(id int64, pause bool) error {
	app := c.Self().Get(id)
	if app == nil || app.State == NotInstalled {
		return fmt.Errorf("App %d is not installed", id)
	}

	// resume the app first so it can be stopped
	pid, err := c.host.AppPid(app)
	if err != nil {
		return err
	}
	if pid != 0 {
		err = c.host.ResumePid(pid)
		if err != nil {
			return err
		}
	}

	if pause {
		err = c.host.PausePid(pid)
	} else {
		err = c.host.StopApp(app)
	}
	if err != nil {
		return err
	}
	if !pause {
		return c.host.SetEnabled(app, false)
	}
	return nil
}

func (c *Controller) AppPid(id int64) (int, error) {
	app := c.Self().Get(id)
	if app == nil || app.State == NotInstalled {
		return 0, fmt.Errorf("App %d is not installed", id)
	}

	return c.host.AppPid(app)
}

func (c *Controller) Install(id int64, version int64) (*App, error) {
	// First download app info
	// create or update app record with latest app info
	// mark app state as Installing
	app, err := c.Self().PreInstallation(id, version)

	// Begin the actual download and installation
	if app != nil && app.State == Installing {
		if err == nil {
			var apkPath, tmpDir string
			// write to tempfile
			tmpDir, err = ioutil.TempDir(c.tempDir, "download")
			if err == nil {
				defer func() { os.RemoveAll(tmpDir) }()
				apkPath, err = c.Self().DownloadApp(app, tmpDir)
				if err == nil {
					log.Println("Starting install")
					err = c.Self().InstallApp(app, apkPath)
				}
			}
		}

		// on error uninstall and remove it so it is not sitting around wasting space
		if err != nil {
			c.Self().Uninstall(app)
		}
	}
	return app, err
}

func (c *Controller) PreInstallation(id int64, version int64) (*App, error) {
	app := c.Self().Get(id)

	info, err := c.api.AppInfo(id)
	if err != nil {
		return nil, err
	} else if info == nil {
		return nil, errors.New("No app found")
	}

	if len(info.App.Bundle.Url) == 0 {
		return nil, errors.New("No bundle URL returned")
	}

	if app == nil {
		app = new(App)
	} else if app.IsInstalled(info.App.Version) {
		installed, _ := c.host.IsInstalled(app)
		if installed {
			// don't have to reinstall
			log.Println("already installed app %d, version %d", app.Id, info.App.Version)
			go c.Self().SyncStatus(app, true)
			return app, nil
		}
	}
	app.Update(info)
	app.Synced = false
	app.State = Installing
	c.Self().Save(app)
	return app, nil
}

func (c *Controller) DownloadApp(app *App, dir string) (string, error) {
	// clean up older apps if need be
	app_size := c.Self().AppDiskSize(app)
	freed, err := c.Self().EnsureFreeSpace(app_size)
	if err != nil {
		return "", err
	} else if !freed {
		return "", errors.New(fmt.Sprintf("Insufficient disk space.  Required: %d", app_size))
	}

	// ensure dir exists
	os.MkdirAll(dir, 0755)
	os.Chmod(dir, 0755)

	apkPath := path.Join(dir, "app.apk")
	err = c.ResHandler.Download(app.Url, apkPath)
	os.Chmod(apkPath, 0755)
	return apkPath, err
}

func (c *Controller) InstallApp(app *App, apkPath string) error {
	// perform install
	err := c.host.Install(app, apkPath)
	if err != nil {
		return err
	}

	installed, err := c.host.IsInstalled(app)
	if err != nil {
		return err
	} else if !installed {
		return fmt.Errorf("Install failed")
	}

	// disable package
	err = c.host.SetEnabled(app, false)

	app.State = Installed
	app.InstalledDate = time.Now()
	c.Self().Save(app)

	// now see if there are any extensions that also need to be downloaded
	err = c.Self().InstallAppExtensions(app)
	if err != nil {
		log.Println("App Extensions could not be installed: ", err)
		return err
	}

	go c.Self().SyncStatus(app, true)
	return err
}

func (c *Controller) Uninstall(app *App) error {
	if app == nil {
		return nil
	}

	c.host.Uninstall(app)

	// delete UserData cache and extensions for the app
	os.RemoveAll(c.appDataDir(app.Id))
	os.RemoveAll(c.host.AppExtensionsDir(app.BundleId))
	c.db.Where("app_id = ?", app.Id).Delete(AppFile{})
	c.db.Where("id = ?", app.Id).Delete(App{})
	/*
		app.Files = nil
		app.Synced = false
		app.State = NotInstalled
		app.persisted = true
		c.Self().Save(app)
	*/

	go c.Self().SyncStatus(app, false)
	return nil
}

func (c *Controller) InstallAppExtensions(app *App) error {
	obb_path := c.host.AppExtensionsDir(app.BundleId)
	err := os.MkdirAll(obb_path, 0755)
	if err != nil {
		return err
	}
	for _, appfile := range app.Files {
		targetPath := path.Join(obb_path, appfile.Name)
		err := c.ResHandler.Download(appfile.Url, targetPath)
		os.Chmod(targetPath, 0755)
		if err != nil {
			log.Println("Download error: ", err)
			return err
		}
	}
	return nil
}

func (c *Controller) AppDiskSize(app *App) int64 {
	out := app.Size
	for _, appFile := range app.Files {
		out += appFile.Size
	}
	return out
}

// Tries to free space required for a particular app
func (c *Controller) EnsureFreeSpace(size int64) (bool, error) {
	available, err := c.host.StorageAvailable()
	if err != nil {
		return false, err
	} else if available > size {
		return true, nil
	}

	var all_apps []*App
	c.db.Where("state in (?, ?)", NotInstalled, Installed).Order("last_used desc").Find(&all_apps)
	for _, app := range all_apps {
		c.db.Where("app_id = ?", app.Id).Find(&app.Files)
		// uninstall it and check for space
		c.Self().Uninstall(app)

		available, err = c.host.StorageAvailable()
		if err != nil {
			return false, err
		} else if available >= size {
			return true, nil
		}
	}
	if available >= size {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *Controller) ClearData(id int64) error {
	app := c.Self().Get(id)
	if app == nil {
		return nil
	}

	return c.host.ClearData(app)
}

func (c *Controller) SyncStatus(app *App, saveApp bool) {
	c.api.UpdateAppStatus(app.Id, app.State != Installed, app.Version)
	if saveApp {
		app.Synced = true
		c.Self().Save(app)
	}
}

func (c *Controller) CacheData(data *UserData) error {
	if len(data.Url) == 0 {
		return nil
	}
	if data.IsCached(c.dataCacheDir) {
		return nil
	}

	// download and cache first
	dir, err := ioutil.TempDir(c.tempDir, "user_data")
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(dir)
	}()

	path := path.Join(dir, "data")
	if err = c.ResHandler.Download(data.Url, path); err != nil {
		return err
	}

	// cache it
	data.Cache(c.dataCacheDir, path)
	return nil
}

func (c *Controller) RestoreData(app *App, data *UserData) error {
	if len(data.Url) == 0 {
		return nil
	}

	err := c.Self().CacheData(data)
	if err != nil {
		return err
	}
	err = c.host.RestoreData(app, data.DataPath(c.dataCacheDir))
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) appDataDir(appId int64) string {
	return path.Join(c.dataCacheDir, strconv.FormatInt(appId, 10))
}
