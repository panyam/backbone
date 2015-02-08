package app

type ControllerInterface interface {
	CacheData(data *UserData) error
	Get(id int64) *App
	Save(app *App) error
	Install(id int64, version int64) (*App, error)
	PreInstallation(id int64, version int64) (*App, error)
	DownloadApp(app *App, dir string) (string, error)
	InstallApp(app *App, apkPath string) error
	Uninstall(app *App) error
	InstalledApps() []*App
	InstallAppExtensions(app *App) error
	SizeUsed() int64
	AppDiskSize(app *App) int64
	EnsureFreeSpace(size int64) (bool, error)
	SyncStatus(app *App, saveApp bool)
}

type Adder interface {
	Add(a, b int) int
}

type ResourceHandlerInterface interface {
	Download(url string, target string) error
	Upload(source string, url string) error
}

type HostInterface interface {
	// storage available for apps
	StorageAvailable() (int64, error)

	// data management
	ClearData(app *App) error
	RestoreData(app *App, dataFile string) error
	BackupData(app *App) (string, error)

	AppExtensionsDir(bundle_id string) string

	// app lifecycle
	Install(app *App, path string) error
	Uninstall(app *App) error
	IsInstalled(app *App) (bool, error)
	SetEnabled(app *App, enabled bool) error
	StartApp(app *App, params AppParams) error
	StopApp(app *App) error
	AppPid(app *App) (int, error)
	PausePid(pid int) error
	ResumePid(pid int) error

	// device settings
	SetInternetEnabled(enabled bool) error
	RemoveFiles(files ...string) error
	SetAndroidId(androidId string) error
}
