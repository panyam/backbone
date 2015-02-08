package app

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
)

const SERVICE_ACTIVITY = "com.voxel.vxandroservice/.MainActivity"
const ACTIVITY_FLAGS = "0x04000000"

type Host struct {
}

func (h *Host) AppExtensionsDir(bundle_id string) string {
	return path.Join("/mnt/shell/emulated/obb", bundle_id)
}

func (h *Host) IsInstalled(app *App) (bool, error) {
	output, err := h.runDeviceShell("pm", "list", "package", app.BundleId)
	if err != nil {
		return false, err
	}

	return parsePackageFromListPackage(output, app.BundleId), nil
}

func (h *Host) ClearData(app *App) error {
	uid, err := h.appUserId(app)
	if err != nil {
		return err
	}
	_, err = h.runDeviceShell("find", h.dataDir(app), "-user", uid,
		"-mindepth", "1", "-exec", "rm", "-rf", "{}", "\\;")

	h.runDeviceShell("rm", "-rf", "/sdcard/Android/data/*", "/sdcard/.*")

	return err
}

func (h *Host) SetEnabled(app *App, enabled bool) error {
	action := "enable"
	if !enabled {
		action = "disable"
	}
	_, err := h.runDeviceShell("pm", action, app.BundleId)
	return err
}

func (h *Host) StartApp(app *App, params AppParams) error {
	var fullName string
	dotIndex := strings.Index(app.MainActivity, ".")
	if dotIndex == 0 {
		fullName = fmt.Sprintf("%s%s", app.BundleId, app.MainActivity)
	} else if dotIndex < 0 {
		fullName = fmt.Sprintf("%s.%s", app.BundleId, app.MainActivity)
	} else {
		fullName = app.MainActivity
	}

	action := "com.voxel.vxandroservice.SET_PORTRAIT"
	if params.Orientation != 0 && params.Orientation != 180 {
		action = "com.voxel.vxandroservice.SET_LANDSCAPE"
	}

	_, err := h.runDeviceShell("am", "start", "-a", action, "-n",
		SERVICE_ACTIVITY,
		"--es", "com.voxel.vxandroservice.LAUNCH_PKG_EXTRA", app.BundleId,
		"--es", "com.voxel.vxandroservice.LAUNCH_ACT_EXTRA", fullName,
		"-f", ACTIVITY_FLAGS)
	return err
}

func (h *Host) StopApp(app *App) error {
	_, err := h.runDeviceShell("am", "force-stop", app.BundleId)
	return err
}

func (h *Host) AppPid(app *App) (int, error) {
	output, err := h.runDeviceShell("ps")
	if err != nil {
		return 0, err
	}

	return parsePidFromPs(output, app.BundleId)
}

func (h *Host) PausePid(pid int) error {
	_, err := h.runDeviceShell("kill", "-s", "SIGSTOP", strconv.Itoa(pid))
	return err
}

func (h *Host) ResumePid(pid int) error {
	_, err := h.runDeviceShell("kill", "-s", "SIGCONT", strconv.Itoa(pid))
	return err
}

func (h *Host) StorageAvailable() (int64, error) {
	output, err := h.runDeviceShell("busybox", "df", "-k", "/data")
	if err != nil {
		return 0, err
	}

	return parseSizeFromDf(output)
}

func (h *Host) SetInternetEnabled(enabled bool) error {
	//_, err := h.runDeviceShell("")
	return nil
}

func (h *Host) RemoveFiles(files ...string) error {
	args := make([]string, len(files)+2)
	args[0] = "rm"
	args[1] = "-rf"
	copy(args[2:], files)
	_, err := h.runDeviceShell(args...)
	return err
}

func (h *Host) SetAndroidId(androidId string) error {
	_, err := h.runDeviceShell("sqlite3",
		"/data/data/com.android.providers.settings/databases/settings.db",
		fmt.Sprintf(`"update secure set value='%s' where name = 'android_id'"`, androidId),
	)
	return err
}

func (h *Host) appUserId(app *App) (string, error) {
	output, err := h.runDeviceShell("busybox", "ls", "-ld", h.dataDir(app))
	if err != nil {
		return "", err
	}
	return parseUserIdFromLs(output)
}

func (h *Host) dataDir(app *App) string {
	return fmt.Sprintf("/data/data/%s", app.BundleId)
}

// format: drwxr-x--x    2 u0_a21   u0_a21        4096 Nov  7 16:49 /data/data/com.zoosk.zoosk
func parseUserIdFromLs(output string) (string, error) {
	parts := regexp.MustCompile("\\s+").Split(output, -1)
	if len(parts) < 3 {
		return "", fmt.Errorf("Could not parse output: %v", output)
	}
	return parts[2], nil
}

// Filesystem           1K-blocks      Used Available Use% Mounted on
// /dev/block/sdb3        5160576   1208208   3952368  23% /data
func parseSizeFromDf(output string) (int64, error) {
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("Too few lines.  Expected atleast 2, Got %d - %v", len(lines), output)
	}
	parts := regexp.MustCompile("\\s+").Split(lines[1], -1)
	if len(parts) < 4 {
		return 0, fmt.Errorf("Len(parts): %d, Incorrect output: %v", len(parts), lines[1])
	}
	var freespace int64
	_, err := fmt.Sscanf(parts[3], "%d", &freespace)
	freespace = freespace * 1000
	return freespace, err
}

// u0_a51    950   141   541312 31080 ffffffff b7528f37 S com.google.android.gms.wearable
// u0_a68    1950  141   529136 39824 ffffffff b7528f37 S com.google.android.apps.plus
func parsePidFromPs(output string, bundleId string) (int, error) {
	lines := strings.Split(output, "\n")
	for idx, line := range lines {
		if idx == 0 {
			continue
		}
		parts := regexp.MustCompile("\\s+").Split(line, -1)
		if len(parts) < 9 {
			return 0, fmt.Errorf("Couldn't parse pid from %s", line)
		}
		if parts[8] == bundleId {
			var pid int
			_, err := fmt.Sscanf(parts[1], "%d", &pid)
			if err != nil {
				return 0, err
			}
			return pid, nil
		}
	}

	// process isn't running
	return 0, nil
}

// 4.4: package:/data/app/com.voxel.sushishooter-1.apk=com.voxel.sushishooter
// 4.x: package:com.voxel.sushishooter
func parsePackageFromListPackage(output string, bundleId string) bool {
	output = strings.TrimSpace(output)
	if strings.Index(output, "=") > 0 {
		// 4.4 format
		index := strings.Index(output, fmt.Sprintf("=%s", bundleId))
		if index == len(output)-len(bundleId)-1 {
			return true
		}
	} else {
		index := strings.Index(output, fmt.Sprintf(":%s", bundleId))
		if index == len(output)-len(bundleId)-1 {
			return true
		}
	}

	log.Printf("package not installed, output: %s", output)
	return false
}
