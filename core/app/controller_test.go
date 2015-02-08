package app

import (
	"errors"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"cloudvm/gocommon/mocks"
	"cloudvm/gocommon/model"
	"cloudvm/solo/data"
	"cloudvm/solo/qa"
	"code.google.com/p/gomock/gomock"
	. "gopkg.in/check.v1"
)

var (
	realController    *Controller
	stubController    *StubController
	controller        ControllerInterface
	mockApiInterface  *mocks.MockApiInterface
	realResHandler    ResourceHandlerInterface
	mockResHandler    *MockResourceHandlerInterface
	mockHostInterface *MockHostInterface
	mockController    *gomock.Controller
	nextAppFileId     int64
)

type ControllerSuite struct{}

var _ = Suite(&ControllerSuite{})

func waitForFlagChange(flag *bool) {
	for i := 0; i < 100 && !*flag; i++ {
		time.Sleep(10) // so that the SyncStatus go routine is also called
	}
}

func installDummyFile(size int64) func(string, string) {
	return func(url, path string) {
		// log.Printf("Downloading %s -> %s.  %d bytes", url, path, size)
		file, _ := os.Create(path)
		file.Truncate(size)
		os.Chtimes(path, time.Now(), time.Now())
		file.Close()
	}
}

func SetupTest(tr gomock.TestReporter, stub bool) {
	// setup code before a test
	nextAppFileId = 1

	// drop existing table
	mockController = gomock.NewController(tr)
	mockApiInterface = mocks.NewMockApiInterface(mockController)
	data.Context.ApiClient = mockApiInterface
	data.Context.Config.DbFile = "test.db"
	data.Context.OpenDatabase()
	data.Context.Db.DropTableIfExists(&App{})
	data.Context.Db.DropTableIfExists(&AppFile{})

	var err error
	realController, err = NewController(data.Context)
	if err != nil {
		log.Printf("%v", err)
		panic("Could not create Controller")
	} else {
		// inject mocks
		mockHostInterface = NewMockHostInterface(mockController)
		realController.host = mockHostInterface

		realResHandler = realController.ResHandler
		mockResHandler = NewMockResourceHandlerInterface(mockController)
		realController.ResHandler = mockResHandler

		stubController = &StubController{*realController, nil, nil, nil, nil}
		stubController.Cls = stubController

		if stub {
			controller = stubController
		} else {
			controller = realController
		}
	}
}

func (s *ControllerSuite) SetupTestApp(appId int64, appName string, appState uint32, apkSize int64, fileSizes ...int64) *App {
	app := MakeTestApp(appId, appName, apkSize, fileSizes...)
	app.State = appState
	controller.Save(app)
	return controller.Get(appId)
}

func (s *ControllerSuite) TestNewController(c *C) {
	SetupTest(c, false)
	file, err := os.Open(realController.tempDir)
	c.Assert(file, Not(Equals), nil)
	c.Assert(err, Equals, nil)

	fileinfo, err := file.Stat()
	c.Assert(fileinfo, Not(Equals), nil)
	c.Assert(fileinfo.IsDir(), Equals, true)
	c.Assert(fileinfo.Mode().Perm(), Equals, os.FileMode(0755))
	c.Assert(err, Equals, nil)

	// ensure no apps installed
	c.Assert(len(controller.InstalledApps()), Equals, 0)
}

func (s *ControllerSuite) TestCreateSchema(c *C) {
	SetupTest(c, false)
	dbPath := path.Join(data.Context.DataRoot, data.Context.Config.DbFile)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		c.Error("DB not created")
		return
	}
}

func (s *ControllerSuite) TestCacheData(c *C) {
	SetupTest(c, false)
	realController.ResHandler = realResHandler
	userData := UserData{
		Id:    1,
		AppId: 3,
		Url:   "https://s3-us-west-2.amazonaws.com/cloudvm-public/test.vxl",
	}
	dataRoot := path.Join(data.Context.DataRoot, "user_data")

	// remove user_data folder first
	app_dir := path.Join(dataRoot, strconv.FormatInt(userData.AppId, 10))
	err := os.RemoveAll(app_dir)
	if err != nil {
		log.Println("RemoveDir Err: ", err)
	}

	c.Assert(userData.IsCached(dataRoot), Equals, false)

	err = controller.CacheData(&userData)
	c.Assert(err, Equals, nil)
	c.Assert(userData.IsCached(dataRoot), Equals, true)
}

func (s *ControllerSuite) TestGet_NotFound(c *C) {
	SetupTest(c, false)
	// test app not existing
	app := controller.Get(10)
	c.Assert(app, Equals, (*App)(nil))
}

func (s *ControllerSuite) TestSave_Simple(c *C) {
	SetupTest(c, false)
	// Save an app and ensure it exists
	app := MakeTestApp(1, "App1", 50000)
	controller.Save(app)
	saved_app := controller.Get(1)

	c.Assert(saved_app, DeepEquals, app)
}

func (s *ControllerSuite) TestSave_WithObb(c *C) {
	SetupTest(c, false)
	// Save an app and ensure it exists
	app := MakeTestApp(1, "App1", 50000, 10000, 20000, 30000, 40000, 50000)
	controller.Save(app)

	saved_app := controller.Get(1)
	c.Assert(saved_app, DeepEquals, app)
}

// Test Pre installation

func (s *ControllerSuite) TestPreInstall_AppMarkedAsInstalledButNotInstalled(c *C) {
	SetupTest(c, false)
	appId := int64(1)
	app := s.SetupTestApp(appId, "App1", Installed, 50000)
	app = controller.Get(appId)

	appInfo := model.AppInfoResponse{
		App: model.App{
			Id:     appId,
			Bundle: model.AppBundle{Url: "http://www.voxel.com/testapk"},
		},
	}
	gomock.InOrder(
		mockApiInterface.EXPECT().AppInfo(appId).Return(&appInfo, nil),
		mockHostInterface.EXPECT().IsInstalled(gomock.Any()).Return(false, nil),
	)

	installed_app, err := controller.PreInstallation(appId, 0)
	installedApps := controller.InstalledApps()
	c.Assert(len(installedApps), Equals, 0)
	c.Assert(err, Equals, nil)
	c.Assert(installed_app.Id, Equals, app.Id)
	c.Assert(installed_app.Synced, Equals, false)
	c.Assert(installed_app.State, Equals, uint32(Installing))
}

func (s *ControllerSuite) TestPreInstall_NoAppInfo(c *C) {
	SetupTest(c, false)
	// crete and save a simple app
	appId := int64(1)

	mockApiInterface.EXPECT().AppInfo(int64(1)).Return(nil, nil)

	_, err := controller.Install(appId, 0)
	c.Assert(err.Error(), Equals, "No app found")
}

func (s *ControllerSuite) TestPreInstall_NoAppBundleUrl(c *C) {
	SetupTest(c, false)
	// create and save a simple app
	appId := int64(1)
	appInfo := model.AppInfoResponse{
		App: model.App{
			Bundle: model.AppBundle{Url: ""},
		},
	}
	mockApiInterface.EXPECT().AppInfo(appId).Return(&appInfo, nil)

	_, err := controller.PreInstallation(appId, 0)
	c.Assert(err.Error(), Equals, "No bundle URL returned")
}

func (s *ControllerSuite) TestPreInstall_Success(c *C) {
	SetupTest(c, false)
	// create and save a simple app
	appId := int64(1)
	appInfo := model.AppInfoResponse{
		App: model.App{
			Bundle: model.AppBundle{Url: "http://www.voxel.com/testapk"},
		},
	}
	mockApiInterface.EXPECT().AppInfo(appId).Return(&appInfo, nil)

	app, err := controller.PreInstallation(appId, 0)
	c.Assert(err, Equals, nil)
	c.Assert(app.Url, Equals, appInfo.App.Bundle.Url)
	c.Assert(app.Id, Equals, appId)
	c.Assert(app.State, Equals, uint32(Installing))
	c.Assert(app.Synced, Equals, false)
}

// Test apk Downloading

func (s *ControllerSuite) TestDownloadApp_NotEnoughDiskSpace(c *C) {
	SetupTest(c, false)
	app := MakeTestApp(1, "App1", 50000)
	app.Url = "http://www.voxel.com/testapk"
	controller.Save(app)
	app = controller.Get(1)

	update_status_called := false
	updateAppStatus := func(appId int64, deleted bool, version int64) error {
		update_status_called = true
		return nil
	}

	mockApiInterface.EXPECT().UpdateAppStatus(gomock.Any(), gomock.Any(), gomock.Any()).Do(updateAppStatus).Return(nil).AnyTimes()
	mockHostInterface.EXPECT().StorageAvailable().Return(int64(0), nil).AnyTimes()
	mockHostInterface.EXPECT().AppExtensionsDir(app.BundleId).Return(path.Join("/tmp/obb", app.BundleId))
	mockHostInterface.EXPECT().Uninstall(gomock.Any()).Return(nil).AnyTimes()

	_, err := controller.DownloadApp(app, "/tmp/downloads")
	c.Assert(err.Error(), qa.HasPrefix, "Insufficient disk space")
}

func (s *ControllerSuite) TestDownloadApp_Success(c *C) {
	SetupTest(c, false)
	app := MakeTestApp(1, "App1", 50000)
	app.Url = "http://www.voxel.com/testapk"
	controller.Save(app)
	app = controller.Get(1)

	mockHostInterface.EXPECT().StorageAvailable().Return(int64(100000), nil).AnyTimes()
	mockResHandler.EXPECT().Download(app.Url, gomock.Any()).Return(nil).AnyTimes()

	// TODO: do other checks like no apps having been uninstalled etc

	apkPath, err := controller.DownloadApp(app, "/tmp/downloads")
	c.Assert(err, Equals, nil)
	c.Assert(apkPath, qa.HasSuffix, "/app.apk")
}

func (s *ControllerSuite) TestInstallApp_HostInstallError(c *C) {
	SetupTest(c, false)
	app := s.SetupTestApp(1, "App1", Installing, 50000)

	mockHostInterface.EXPECT().Install(app, gomock.Any()).Return(errors.New("SomeError"))

	err := controller.InstallApp(app, "some_apk_path")

	app = controller.Get(1)
	c.Assert(err.Error(), Equals, "SomeError")
	// state should be the old state
	c.Assert(app.State, Equals, uint32(Installing))
}

func (s *ControllerSuite) TestInstallApp_HostIsInstalledReturnsError(c *C) {
	SetupTest(c, false)
	app := s.SetupTestApp(1, "App1", Installing, 50000)

	gomock.InOrder(
		mockHostInterface.EXPECT().Install(app, gomock.Any()).Return(nil),
		mockHostInterface.EXPECT().IsInstalled(app).Return(false, errors.New("IsInstalledError")),
	)

	err := controller.InstallApp(app, "some_apk_path")

	app = controller.Get(1)
	c.Assert(err.Error(), Equals, "IsInstalledError")
	// state should be the old state
	c.Assert(app.State, Equals, uint32(Installing))
}

func (s *ControllerSuite) TestInstallApp_SuccessWithNoExtensions(c *C) {
	SetupTest(c, false)
	app := s.SetupTestApp(1, "App1", Installing, 50000)

	update_status_called := false
	updateAppStatus := func(appId int64, deleted bool, version int64) error {
		update_status_called = true
		return nil
	}

	gomock.InOrder(
		mockHostInterface.EXPECT().Install(app, gomock.Any()).Return(nil),
		mockHostInterface.EXPECT().IsInstalled(app).Return(true, nil),
		mockHostInterface.EXPECT().SetEnabled(app, false).Return(nil),
		mockHostInterface.EXPECT().AppExtensionsDir(app.BundleId).Return(path.Join("/tmp/obb", app.BundleId)),
		mockApiInterface.EXPECT().UpdateAppStatus(gomock.Any(), gomock.Any(), gomock.Any()).Do(updateAppStatus).Return(nil).AnyTimes(),
	)

	err := controller.InstallApp(app, "some_apk_path")

	waitForFlagChange(&update_status_called)
	app = controller.Get(1)
	c.Assert(err, Equals, nil)
	c.Assert(app.State, Equals, uint32(Installed))
	c.Assert(app.Synced, Equals, true)

	installedApps := controller.InstalledApps()
	c.Assert(len(installedApps), Equals, 1)
	c.Assert(installedApps[0].Id, Equals, app.Id)
	c.Assert(controller.SizeUsed(), Equals, app.Size)
}

// Test installation of extensions

func (s *ControllerSuite) TestInstallAppExtensions_Success(c *C) {
	SetupTest(c, false)
	app := s.SetupTestApp(1, "App1", 666, 50000, 10000, 20000, 30000)

	// remove test obb folder for now
	os.RemoveAll("/tmp/obb")

	app = controller.Get(1)
	obb_path := path.Join("/tmp/obb", app.BundleId)
	gomock.InOrder(
		mockHostInterface.EXPECT().AppExtensionsDir(app.BundleId).Return(obb_path),
		mockResHandler.EXPECT().Download(gomock.Any(), gomock.Any()).Do(installDummyFile(10000)).Return(nil),
		mockResHandler.EXPECT().Download(gomock.Any(), gomock.Any()).Do(installDummyFile(20000)).Return(nil),
		mockResHandler.EXPECT().Download(gomock.Any(), gomock.Any()).Do(installDummyFile(30000)).Return(nil),
	)

	file, err := os.Open(obb_path)
	c.Assert(file, Equals, (*os.File)(nil))

	err = controller.InstallAppExtensions(app)
	c.Assert(err, Equals, nil)
	c.Assert(app.State, Equals, uint32(666)) // state should NOT be changed

	file, err = os.Open(obb_path)
	c.Assert(file, Not(Equals), nil)
	c.Assert(err, Equals, nil)

	fileinfo, err := file.Stat()
	c.Assert(fileinfo, Not(Equals), nil)
	c.Assert(fileinfo.IsDir(), Equals, true)
	c.Assert(fileinfo.Mode().Perm(), Equals, os.FileMode(0755))
	c.Assert(err, Equals, nil)

	// check obb sizes
	for _, appfile := range app.Files {
		filepath := path.Join(obb_path, appfile.Name)
		file, err = os.Open(filepath)
		fileinfo, err = file.Stat()
		c.Assert(fileinfo.Size(), Equals, appfile.Size)
	}
}

func (s *ControllerSuite) TestInstallApp_SuccessWithExtensions(c *C) {
	SetupTest(c, false)
	app := s.SetupTestApp(1, "App1", Installing, 50000, 10000, 20000, 30000)

	update_status_called := false
	updateAppStatus := func(appId int64, deleted bool, version int64) error {
		update_status_called = true
		return nil
	}

	gomock.InOrder(
		mockHostInterface.EXPECT().Install(app, gomock.Any()).Return(nil),
		mockHostInterface.EXPECT().IsInstalled(app).Return(true, nil),
		mockHostInterface.EXPECT().SetEnabled(app, false).Return(nil),
		mockHostInterface.EXPECT().AppExtensionsDir(app.BundleId).Return(path.Join("/tmp/obb", app.BundleId)),
		mockResHandler.EXPECT().Download(gomock.Any(), gomock.Any()).Do(installDummyFile(10000)).Return(nil),
		mockResHandler.EXPECT().Download(gomock.Any(), gomock.Any()).Do(installDummyFile(20000)).Return(nil),
		mockResHandler.EXPECT().Download(gomock.Any(), gomock.Any()).Do(installDummyFile(30000)).Return(nil),
		mockApiInterface.EXPECT().UpdateAppStatus(gomock.Any(), gomock.Any(), gomock.Any()).Do(updateAppStatus).Return(nil).AnyTimes(),
	)

	err := controller.InstallApp(app, "some_apk_path")

	waitForFlagChange(&update_status_called)
	app = controller.Get(1)
	c.Assert(err, Equals, nil)
	c.Assert(app.State, Equals, uint32(Installed))
	c.Assert(app.Synced, Equals, true)

	installedApps := controller.InstalledApps()
	c.Assert(len(installedApps), Equals, 1)
	c.Assert(installedApps[0].Id, Equals, app.Id)
	c.Assert(controller.SizeUsed(), Equals, app.Size+60000)
}

// Test the overall "Install" method

func (s *ControllerSuite) TestInstall_AppNotYetInstalled(c *C) {
	SetupTest(c, true)
	app := s.SetupTestApp(1, "App1", NotInstalled, 50000, 10000, 20000, 30000)

	stubController.PreInstallationStub = func(controller *Controller, id int64, version int64) (*App, error) {
		app.State = Installing
		return app, nil
	}

	download_called := false
	stubController.DownloadAppStub = func(controller *Controller, app *App, dir string) (string, error) {
		download_called = true
		return "test", nil
	}

	install_called := false
	stubController.InstallAppStub = func(controller *Controller, app *App, apkpath string) error {
		install_called = true
		return nil
	}

	uninstalled := false
	stubController.UninstallStub = func(controller *Controller, app *App) error {
		uninstalled = true
		return nil
	}

	// app entry was created but should be removed on download failed
	// get it in a state where PreInstallation succeeded
	_, err := controller.Install(app.Id, 0)
	c.Assert(err, Equals, nil)
	c.Assert(download_called, Equals, true)
	c.Assert(install_called, Equals, true)
	c.Assert(uninstalled, Equals, false)
}

func (s *ControllerSuite) SetupTestInstallFailureCommon(c *C, download_error error, install_error error) *App {
	SetupTest(c, true)
	app := MakeTestApp(1, "App1", 50000, 10000, 20000, 30000)
	controller.Save(app)
	app = controller.Get(1)

	stubController.PreInstallationStub = func(controller *Controller, id int64, version int64) (*App, error) {
		return app, nil
	}

	stubController.DownloadAppStub = func(controller *Controller, app *App, dir string) (string, error) {
		return "", download_error
	}

	stubController.InstallAppStub = func(controller *Controller, app *App, apkpath string) error {
		return install_error
	}
	return app
}

func (s *ControllerSuite) TestInstall_DownloadAppStepFailed(c *C) {
	app := s.SetupTestInstallFailureCommon(c, errors.New("SomeError"), nil)
	app.State = Installing

	uninstalled := false
	stubController.UninstallStub = func(controller *Controller, app *App) error {
		uninstalled = true
		return nil
	}

	// app entry was created but should be removed on download failed
	// get it in a state where PreInstallation succeeded
	_, err := controller.Install(app.Id, 0)
	c.Assert(err, Not(Equals), nil)
	c.Assert(uninstalled, Equals, true)
}

func (s *ControllerSuite) TestInstall_InstallAppStepFailed(c *C) {
	app := s.SetupTestInstallFailureCommon(c, nil, errors.New("SomeInstallError"))
	app.State = Installing

	uninstalled := false
	stubController.UninstallStub = func(controller *Controller, app *App) error {
		uninstalled = true
		return nil
	}

	// app entry was created but should be removed on download failed
	// get it in a state where PreInstallation succeeded
	_, err := controller.Install(app.Id, 0)
	c.Assert(err, Not(Equals), nil)
	c.Assert(uninstalled, Equals, true)
}

// Finally test uninstallation
func (s *ControllerSuite) TestUninstall(c *C) {
	SetupTest(c, false)
	appId := int64(1)
	app := s.SetupTestApp(appId, "App1", Installed, 50000, 10000, 20000, 30000)

	installedApps := controller.InstalledApps()
	c.Assert(len(installedApps), Equals, 1)
	c.Assert(installedApps[0].Id, Equals, app.Id)
	c.Assert(controller.SizeUsed(), Equals, app.Size+60000)

	update_status_called := false
	updateAppStatus := func(appId int64, deleted bool, version int64) error {
		update_status_called = true
		return nil
	}

	gomock.InOrder(
		mockHostInterface.EXPECT().Uninstall(gomock.Any()).Return(nil).AnyTimes(),
		mockHostInterface.EXPECT().AppExtensionsDir(app.BundleId).Return(path.Join("/tmp/obb", app.BundleId)),
		mockApiInterface.EXPECT().UpdateAppStatus(gomock.Any(), gomock.Any(), gomock.Any()).Do(updateAppStatus).Return(nil).AnyTimes(),
	)

	err := controller.Uninstall(app)
	c.Assert(err, Equals, nil)

	waitForFlagChange(&update_status_called)
	app = controller.Get(appId)
	installedApps = controller.InstalledApps()
	c.Assert(app, Equals, (*App)(nil))
	c.Assert(len(installedApps), Equals, 0)
	c.Assert(controller.SizeUsed(), Equals, int64(0))
}

func (s *ControllerSuite) TestSyncStatus(c *C) {
	SetupTest(c, false)
	app := s.SetupTestApp(1, "App1", Installed, 50000, 10000, 20000, 30000)

	update_status_called := false
	updateAppStatus := func(appId int64, deleted bool, version int64) error {
		update_status_called = true
		return nil
	}

	mockApiInterface.EXPECT().UpdateAppStatus(gomock.Any(), gomock.Any(), gomock.Any()).Do(updateAppStatus).Return(nil).AnyTimes()

	app.Synced = false
	update_status_called = false
	controller.SyncStatus(app, true)
	c.Assert(app.Synced, Equals, true)
	c.Assert(update_status_called, Equals, true)

	app.Synced = false
	update_status_called = false
	controller.SyncStatus(app, false)
	c.Assert(app.Synced, Equals, false)
	c.Assert(update_status_called, Equals, true)
}
