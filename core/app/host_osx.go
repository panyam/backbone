// +build darwin

package app

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Handles the manual work of interacting with apps on the device

func (h *Host) Install(app *App, path string) error {
	log.Printf("Installing %v\n", app.BundleId)
	_, err := h.runAdb("install", "-r", path)
	if err != nil {
		log.Printf("install command ran, error: %v\n", err)
	}
	return err
}

func (h *Host) Uninstall(app *App) error {
	log.Printf("Uninstalling %v\n", app.BundleId)
	_, err := h.runAdb("uninstall", app.BundleId)
	if err != nil {
		return err
	}
	err = h.RemoveFiles(fmt.Sprintf("/data/app-lib/%s-1", app.BundleId),
		fmt.Sprintf("/data/app-lib/%s-2", app.BundleId))
	return nil
}

func (h *Host) RestoreData(app *App, dataFile string) error {
	return nil
}

func (h *Host) BackupData(app *App) (string, error) {
	return "", nil
}

func (h *Host) runDeviceShell(args ...string) (string, error) {
	fullArgs := make([]string, len(args)+1)
	fullArgs[0] = "shell"
	copy(fullArgs[1:], args)
	return h.runAdb(fullArgs...)
}

func (h *Host) runAdb(args ...string) (string, error) {
	log.Printf("Running %v\n", args)
	cmd := exec.Command("/Applications/Genymotion.app/Contents/MacOS/player.app/Contents/MacOS/tools/adb", args...)
	combined, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	output := string(combined)

	if strings.Index(strings.ToLower(output), "error") == 0 {
		return "", errors.New(fmt.Sprintf("Command returned error: %s", output))
	}

	return output, nil
}
