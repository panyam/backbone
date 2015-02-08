package app

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"cloudvm/gocommon/model"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type AppSuite struct{}

var _ = Suite(&AppSuite{})

func MakeTestAppResponse(appId int64, appName string, apkSize int64, fileSizes ...int64) *model.AppInfoResponse {
	app := MakeTestApp(appId, appName, apkSize, fileSizes...)
	out := &model.AppInfoResponse{
		App: model.App{
			Id:       appId,
			Name:     appName,
			BundleId: app.BundleId,
		},
	}
	if app.Files != nil {
		for _, file := range app.Files {
			appFile := model.AppFile{
				Id:    file.Id,
				AppId: appId,
				Name:  file.Name,
				Url:   file.Url,
				Size:  file.Size,
				Md5:   file.Md5,
			}
			out.App.Files = append(out.App.Files, appFile)
		}
	}
	return out
}

func MakeTestApp(appId int64, appName string, apkSize int64, fileSizes ...int64) *App {
	bundleId := fmt.Sprintf("com.voxel.test.%s", strings.ToLower(appName))
	app := App{
		Id:            appId,
		BundleId:      bundleId,
		Size:          apkSize,
		InstalledDate: time.Now(),
		LastUsed:      time.Now(),
	}
	app.Files = nil
	for index, fileSize := range fileSizes {
		appFile := AppFile{
			Id:    nextAppFileId,
			AppId: appId,
			Size:  fileSize,
			Name:  fmt.Sprintf("%s_file_%d.obb", bundleId, index),
			Url:   fmt.Sprintf("http://voxel.com/files/%s/file_%d.apk", bundleId, index),
		}
		app.Files = append(app.Files, appFile)
		nextAppFileId++
	}
	return &app
}

func (s *AppSuite) TestUpdateApp(c *C) {
	appInfoResponse := MakeTestAppResponse(1, "App1", 50000, 10000, 20000, 30000)
	app := App{}
	app.Update(appInfoResponse)
	c.Assert(app.Id, Equals, appInfoResponse.App.Id)
	c.Assert(len(app.Files), Equals, 3)
}
