package app

import (
	"testing"
)

func TestParseUserId(t *testing.T) {
	output := "drwxr-x--x    2 u0_a21   u0_a21        4096 Nov  7 16:49 /data/data/com.zoosk.zoosk"

	expected := "u0_a21"
	uid, err := parseUserIdFromLs(output)

	if err != nil {
		t.Errorf("Error during parsing, %v", err)
		return
	}
	if uid != expected {
		t.Errorf("Incorrect id parsed, expected %v, actual %v", expected, uid)
		return
	}

	// test malformed string
	malformed := "hello"
	_, err = parseUserIdFromLs(malformed)

	if err == nil {
		t.Errorf("Error not returned")
		return
	}
}

func TestParseSizeFromDf(t *testing.T) {
	output := "Filesystem           1K-blocks      Used Available Use% Mounted on\n"
	output += "/dev/block/sdb3        5160576   1208208   3952368  23% /data"

	expected := int64(3952368000)
	size, err := parseSizeFromDf(output)

	if err != nil {
		t.Errorf("Error during parsing, %v", err)
		return
	}
	if size != expected {
		t.Errorf("Incorrect size parsed, expected %v, actual %v", expected, size)
		return
	}

	// test malformed string
	malformed := "hello"
	_, err = parseSizeFromDf(malformed)

	if err == nil {
		t.Errorf("Error not returned")
		return
	}
}

func TestParsePidFromPs(t *testing.T) {
	output := "First line is header\n"
	output += "u0_a51    950   141   541312 31080 ffffffff b7528f37 S com.google.android.gms.wearable\n"
	output += "u0_a68    1950  141   529136 39824 ffffffff b7528f37 S com.google.android.apps.plus"

	expected := 1950
	bundleId := "com.google.android.apps.plus"
	size, err := parsePidFromPs(output, bundleId)

	if err != nil {
		t.Errorf("Error during parsing, %v", err)
		return
	}
	if size != expected {
		t.Errorf("Incorrect pid parsed, expected %v, actual %v", expected, size)
		return
	}

	// test malformed string
	malformed := "hello\nworld"
	_, err = parsePidFromPs(malformed, bundleId)

	if err == nil {
		t.Errorf("Error not returned")
		return
	}
}

func TestParsePackageList(t *testing.T) {
	bundle := "com.voxel.sushishooter"
	shouldFail := parsePackageFromListPackage("", bundle)
	if shouldFail {
		t.Errorf("Should've failed to parse")
		return
	}

	succeed44 := parsePackageFromListPackage("package:/data/app/com.voxel.sushishooter-1.apk=com.voxel.sushishooter", bundle)
	if !succeed44 {
		t.Errorf("Should've suceeded")
		return
	}

	succeed4x := parsePackageFromListPackage("package:com.voxel.sushishooter", bundle)
	if !succeed4x {
		t.Errorf("Should've suceeded")
		return
	}
}
