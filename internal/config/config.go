package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
)

const (
	TNFS_ROOT_PATH_KEY = "tnfsRootPath"
)

type Config struct {
	ExePath      string
	TnfsRootPath string
}

func (c *Config) UpdateRootPath(newPath string) {
	a := fyne.CurrentApp()
	c.TnfsRootPath = newPath
	a.Preferences().SetString(TNFS_ROOT_PATH_KEY, newPath)

}

func LoadConfig() (*Config, error) {
	rootPath := loadDefaultRootPath()
	exePath, err := locateTnfsdExecutable()
	if err != nil {
		return &Config{}, err
	}
	cfg := &Config{ExePath: exePath, TnfsRootPath: rootPath}
	return cfg, nil
}

func locateTnfsdExecutable() (string, error) {
	dir := "."
	exeName := "tnfsd"
	if runtime.GOOS == "windows" {
		exeName = "tnfsd.exe"
	}

	currentExePath, _ := os.Executable()
	if currentExePath != "" {
		dir = filepath.Dir(currentExePath)
	}

	exePath := filepath.Join(dir, exeName)

	fmt.Println(currentExePath)
	fmt.Println(exePath)

	_, err := exec.LookPath(exePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return exePath, nil
}

func loadDefaultRootPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		dirname = "."
	}
	prefs := fyne.CurrentApp().Preferences()
	return prefs.StringWithFallback(TNFS_ROOT_PATH_KEY, dirname)
}
