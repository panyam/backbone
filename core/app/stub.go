package app

import ()

// Create StubController so we can stub out behaviours of particular methods (should be a better way than this)
type StubController struct {
	Controller
	PreInstallationStub func(controller *Controller, id int64, version int64) (*App, error)
	DownloadAppStub     func(controller *Controller, app *App, dir string) (string, error)
	InstallAppStub      func(controller *Controller, app *App, apkpath string) error
	UninstallStub       func(controller *Controller, app *App) error
}

func (sc *StubController) PreInstallation(id int64, version int64) (*App, error) {
	if sc.PreInstallationStub != nil {
		return sc.PreInstallationStub(&sc.Controller, id, version)
	} else {
		return sc.Controller.PreInstallation(id, version)
	}
}

func (sc *StubController) DownloadApp(app *App, dir string) (string, error) {
	if sc.DownloadAppStub != nil {
		return sc.DownloadAppStub(&sc.Controller, app, dir)
	} else {
		return sc.Controller.DownloadApp(app, dir)
	}
}

func (sc *StubController) InstallApp(app *App, apkpath string) error {
	if sc.InstallAppStub != nil {
		return sc.InstallAppStub(&sc.Controller, app, apkpath)
	} else {
		return sc.Controller.InstallApp(app, apkpath)
	}
}

func (sc *StubController) Uninstall(app *App) error {
	if sc.UninstallStub != nil {
		return sc.UninstallStub(&sc.Controller, app)
	} else {
		return sc.Controller.Uninstall(app)
	}
}

func (sc *StubController) CacheData(data *UserData) error {
	return sc.Controller.CacheData(data)
}

func (sc *StubController) Get(id int64) *App {
	return sc.Controller.Get(id)
}

func (sc *StubController) Save(app *App) error {
	return sc.Controller.Save(app)
}

func (sc *StubController) Install(id int64, version int64) (*App, error) {
	return sc.Controller.Install(id, version)
}

func (sc *StubController) InstalledApps() []*App {
	return sc.Controller.InstalledApps()
}

func (sc *StubController) InstallAppExtensions(app *App) error {
	return sc.Controller.InstallAppExtensions(app)
}

func (sc *StubController) SizeUsed() int64 {
	return sc.Controller.SizeUsed()
}

func (sc *StubController) AppDiskSize(app *App) int64 {
	return sc.Controller.AppDiskSize(app)
}

func (sc *StubController) EnsureFreeSpace(size int64) (bool, error) {
	return sc.Controller.EnsureFreeSpace(size)
}

func (sc *StubController) SyncStatus(app *App, saveApp bool) {
	sc.Controller.SyncStatus(app, saveApp)
}
