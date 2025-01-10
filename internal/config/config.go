package config

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"github.com/emersion/go-autostart"
)

const (
	TNFS_ROOT_PATH_KEY   = "tnfsRootPath"
	ALLOW_BACKGROUND_KEY = "allowBackground"
)

type Config struct {
	ExePath         string
	TnfsRootPath    string
	Hostname        string
	AllowBackground bool
	StartAtLogin    bool
	autostartApp    *autostart.App
}

func (c *Config) SetRootPath(newPath string) {
	c.TnfsRootPath = newPath
	fyne.CurrentApp().Preferences().SetString(TNFS_ROOT_PATH_KEY, newPath)
}

func (c *Config) SetAllowBackground(newVal bool) {
	c.AllowBackground = newVal
	fyne.CurrentApp().Preferences().SetBool(ALLOW_BACKGROUND_KEY, newVal)
}

func (c *Config) SetStartAtLogin(enable bool) {
	var err error
	if enable {
		err = c.autostartApp.Enable()
	} else {
		err = c.autostartApp.Disable()
	}
	if err == nil {
		c.StartAtLogin = enable
	}
}

func LoadConfig() (*Config, error) {
	exePath, err := locateTnfsdExecutable()
	if err != nil {
		return &Config{}, err
	}

	prefs := fyne.CurrentApp().Preferences()
	autostartApp, _ := makeAutostartApp()

	cfg := &Config{
		ExePath:         exePath,
		Hostname:        getHostnameOrIP(),
		TnfsRootPath:    getRootPath(prefs),
		AllowBackground: prefs.BoolWithFallback(ALLOW_BACKGROUND_KEY, false),
		StartAtLogin:    autostartApp.IsEnabled(),

		autostartApp: autostartApp,
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

func getRootPath(prefs fyne.Preferences) string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		dirname = "."
	}
	return prefs.StringWithFallback(TNFS_ROOT_PATH_KEY, dirname)
}

func makeAutostartApp() (*autostart.App, error) {
	a := fyne.CurrentApp()
	currentExePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	aa := &autostart.App{
		Name:        a.Metadata().ID,
		DisplayName: a.Metadata().Name,
		Exec:        []string{currentExePath, "autorun"},
	}
	return aa, nil
}
