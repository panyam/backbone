package app

import (
	"cloudvm/gocommon/model"
	"time"
)

type AppFile struct {
	Id          int64
	AppId       int64 // Foreign key to App
	Name        string
	Url         string
	Size        int64
	Md5         string
	State       uint32
	WrittenSize int64
}

type App struct {
	Id            int64
	BundleId      string
	MainActivity  string
	Url           string
	Version       int64
	Orientation   int
	State         uint32
	Size          int64
	Synced        bool
	persisted     bool
	InstalledDate time.Time
	LastUsed      time.Time
	RunCount      uint32
	Files         []AppFile
}

type AppState uint32

const (
	NotInstalled = iota
	Installing   = iota
	Installed    = iota
)

type AppParams struct {
	AndroidId   string
	Orientation int
}

func (a *AppFile) Update(info *model.AppFile) {
	a.Id = info.Id
	a.AppId = info.AppId
	a.Name = info.Name
	a.Url = info.Url
	a.Size = info.Size
	a.Md5 = info.Md5
	a.State = NotInstalled
	a.WrittenSize = 0
}

func (a *App) Update(info *model.AppInfoResponse) {
	a.Id = info.App.Id
	a.BundleId = info.App.BundleId
	a.MainActivity = info.App.MainActivity
	a.Version = info.App.Version
	a.Orientation = info.App.Orientation
	a.Size = int64(info.App.Bundle.Size)
	a.Url = info.App.Bundle.Url

	// process the app files
	a.Files = nil
	if info.App.Files != nil {
		a.Files = make([]AppFile, 0, len(info.App.Files))
		for _, file := range info.App.Files {
			appFile := AppFile{}
			appFile.Update(&file)
			a.Files = append(a.Files, appFile)
		}
	}
}

func (a *App) IsInstalled(version int64) bool {
	return a.State == Installed && a.Version == version
}
