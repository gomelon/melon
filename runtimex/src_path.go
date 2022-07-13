package runtimex

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// MainSrcPath get the main.go absolute path
func MainSrcPath() string {
	dir := MainSrcPathOnRun()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return MainSrcPathOnGoRun()
	}
	return dir
}

// MainSrcPathOnRun get the main.go absolute on run executable file
func MainSrcPathOnRun() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// MainSrcPathOnGoRun get the main.go absolute on go run
func MainSrcPathOnGoRun() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
