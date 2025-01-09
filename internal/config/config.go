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
	TNFS_ROOT_PATH_KEY   = "tnfsRootPath"
	ALLOW_BACKGROUND_KEY = "allowBackground"
	START_AT_LOGIN_KEY   = "startAtLogin"
)

type Config struct {
	ExePath         string
	TnfsRootPath    string
	Hostname        string
	AllowBackground bool
	StartAtLogin    bool
}

func (c *Config) SetRootPath(newPath string) {
	c.TnfsRootPath = newPath
	fyne.CurrentApp().Preferences().SetString(TNFS_ROOT_PATH_KEY, newPath)
}

func (c *Config) SetAllowBackground(newVal bool) {
	c.AllowBackground = newVal
	fyne.CurrentApp().Preferences().SetBool(ALLOW_BACKGROUND_KEY, newVal)
}

func (c *Config) SetStartAtLogin(newVal bool) {
	c.StartAtLogin = newVal
	fyne.CurrentApp().Preferences().SetBool(START_AT_LOGIN_KEY, newVal)
}

func LoadConfig() (*Config, error) {
	exePath, err := locateTnfsdExecutable()
	if err != nil {
		return &Config{}, err
	}

	prefs := fyne.CurrentApp().Preferences()
	cfg := &Config{
		ExePath:         exePath,
		Hostname:        getHostnameOrIP(),
		TnfsRootPath:    getRootPath(prefs),
		AllowBackground: prefs.BoolWithFallback(ALLOW_BACKGROUND_KEY, false),
		StartAtLogin:    prefs.BoolWithFallback(START_AT_LOGIN_KEY, false),
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
	exePath = "bin/tnfsd-bsd"

	fmt.Println(currentExePath)
	fmt.Println(exePath)

	_, err := exec.LookPath(exePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return exePath, nil
}

func getRootPath(prefs fyne.Preferences) string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		dirname = "."
	}
	return prefs.StringWithFallback(TNFS_ROOT_PATH_KEY, dirname)
}
