//go:generate fyne bundle --pkg ui -o bundled.go ../../assets/TrayIcon.svg

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/fujiNetWIFI/tnfs-gui/internal/tnfs"
)

func loadSystemTrayIcon() *theme.ThemedResource {
	return theme.NewThemedResource(resourceTrayIconPng)
}

func makeSystemMenu(ui *UI, server *tnfs.Server) *fyne.Menu {
	a := fyne.CurrentApp()
	desk, ok := a.(desktop.App)
	if !ok {
		return nil
	}

	startStop := fyne.NewMenuItem("Start Server", func() {
		if server.Status == tnfs.STOPPED {
			server.Start()
		} else if server.Status == tnfs.STARTED {
			server.Stop()
		}
	})

	m := fyne.NewMenu("TNFS Server Manager",
		fyne.NewMenuItem("Open TNFS Server Manager", func() {
			ui.MainWindow.Show()
			if ui.cfg.AllowBackground {
				ShowInDock()
			}
		}),
		fyne.NewMenuItemSeparator(),
		startStop,
	)

	ui.On(tnfs.StatusChange, func(e tnfs.Event) {
		switch server.Status {
		case tnfs.STARTED:
			startStop.Label = "Stop Server"
			startStop.Disabled = false
		case tnfs.STOPPED:
			startStop.Label = "Start Server"
			startStop.Disabled = false
		case tnfs.STARTING:
			startStop.Disabled = true
		case tnfs.STOPPING:
			startStop.Disabled = true
		}
		m.Refresh()
	})

	icon := loadSystemTrayIcon()
	desk.SetSystemTrayIcon(icon)
	desk.SetSystemTrayMenu(m)

	return m
}
