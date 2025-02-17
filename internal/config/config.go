package config

import (
	"net"
	"os"

	"fyne.io/fyne/v2"
	"github.com/emersion/go-autostart"
)

const (
	TNFS_ROOT_PATH_KEY   = "tnfsRootPath"
	READ_ONLY_KEY        = "readonly"
	ALLOW_BACKGROUND_KEY = "allowBackground"
)

type Config struct {
	TnfsRootPath    string
	ReadOnly        bool
	Hostname        string
	AllowBackground bool
	StartAtLogin    bool
	autostartApp    *autostart.App
}

func (c *Config) SetRootPath(newPath string) {
	c.TnfsRootPath = newPath
	fyne.CurrentApp().Preferences().SetString(TNFS_ROOT_PATH_KEY, newPath)
}

func (c *Config) SetReadOnly(readonly bool) {
	c.ReadOnly = readonly
	fyne.CurrentApp().Preferences().SetBool(READ_ONLY_KEY, readonly)
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
	prefs := fyne.CurrentApp().Preferences()
	autostartApp, _ := makeAutostartApp()

	cfg := &Config{
		Hostname:        getHostnameOrIP(),
		TnfsRootPath:    getRootPath(prefs),
		ReadOnly:        prefs.BoolWithFallback(READ_ONLY_KEY, false),
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
