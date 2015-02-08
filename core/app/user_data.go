package app

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

type UserData struct {
	Id         int64
	AppId      int64
	Url        string
	ModifiedAt int64
}

func (u *UserData) IsCached(cacheRoot string) bool {
	path := u.DataPath(cacheRoot)

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// check timestamp against modified at
	return u.ModifiedAt == info.ModTime().Unix()
}

func (u *UserData) DataPath(cacheRoot string) string {
	return path.Join(cacheRoot, strconv.FormatInt(u.AppId, 10),
		fmt.Sprintf("%d.tgz", u.Id))
}

func (u *UserData) Cache(cacheRoot string, file string) error {
	target := u.DataPath(cacheRoot)
	dir := path.Dir(target)
	os.MkdirAll(dir, 0700)

	// move file to the right location
	err := os.Rename(file, target)
	if err != nil {
		return err
	}

	time := time.Unix(u.ModifiedAt, 0)
	return os.Chtimes(target, time, time)
}
