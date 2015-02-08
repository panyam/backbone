// +build android

package app

// Handles the manual work of interacting with apps on the device

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func (h *Host) Install(app *App, path string) error {
	_, err := h.runDeviceShell("pm", "install", "-r", path)
	return err
}

func (h *Host) Uninstall(app *App) error {
	_, err := h.runDeviceShell("pm", "uninstall", app.BundleId)

	os.RemoveAll(h.AppExtensionsDir(app.BundleId))
	app.Files = nil

	return err
}

func (h *Host) RestoreData(app *App, dataFile string) error {
	userId, err := h.appUserId(app)
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "restore_data")
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(dir)
	}()

	_, err = h.runDeviceShell("tar", "xzf", dataFile, "-C", dir)
	if err != nil {
		return err
	}

	sdSource := path.Join(dir, "sdcard", app.BundleId)
	sdTarget := "/sdcard/Android/data/"
	if _, err = os.Stat(sdSource); err == nil {
		// copy to sdcard
		h.runDeviceShell("rm", "-rf", sdTarget+app.BundleId)
		h.runDeviceShell("cp", "-Rp", sdSource, sdTarget)
	}

	packageDir := h.dataDir(app)
	// prepare target for copy
	h.runDeviceShell("find", packageDir, "-user", userId, "-mindepth", "1", "-delete")
	// get rid of lib symlink
	h.runDeviceShell("rm", "-rf", fmt.Sprintf("%s/%s/lib", dir, app.BundleId))
	h.runDeviceShell("cp", "-R", fmt.Sprintf("%s/%s/*", dir, app.BundleId), packageDir)
	h.runDeviceShell("busybox", "chown", "-R", fmt.Sprintf("%s.", userId), packageDir)
	h.runDeviceShell("busybox", "chown", "-R", "system.",
		fmt.Sprintf("/data/app-lib/%s-*", app.BundleId),
		path.Join(packageDir, "lib"))

	return nil
}

func (h *Host) BackupData(app *App) (string, error) {
	return "", nil
}

func (h *Host) runDeviceShell(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("No command passed in")
	}
	fullArgs := make([]string, 2)
	fullArgs[0] = "-c"
	fullArgs[1] = strings.Join(args, " ")

	log.Printf("Running sh %v\n", fullArgs)

	cmd := exec.Command("sh", fullArgs...)
	combined, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	output := string(combined)
	return output, nil
}
