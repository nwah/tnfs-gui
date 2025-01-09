package config

import (
	"fmt"
	"net"
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
	Hostname     string
}

func (c *Config) UpdateRootPath(newPath string) {
	a := fyne.CurrentApp()
	c.TnfsRootPath = newPath
	a.Preferences().SetString(TNFS_ROOT_PATH_KEY, newPath)

}

func LoadConfig() (*Config, error) {
	exePath, err := locateTnfsdExecutable()
	if err != nil {
		return &Config{}, err
	}
	cfg := &Config{
		ExePath:      exePath,
		TnfsRootPath: loadDefaultRootPath(),
		Hostname:     getHostnameOrIP(),
	}
	return cfg, nil
}

func getHostnameOrIP() string {
	host, err := os.Hostname()
	if err == nil && host != "" {
		return host
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
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
