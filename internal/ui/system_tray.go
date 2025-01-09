package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/fujiNetWIFI/tnfs-gui/internal/tnfs"
)

func loadSystemTrayIcon() *theme.ThemedResource {
	icon, err := fyne.LoadResourceFromPath("assets/TrayIcon.svg")
	if err != nil {
		icon = theme.MediaRecordIcon()
	}
	return theme.NewThemedResource(icon)
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
		fyne.NewMenuItem("Show TNFS Server Manager", func() {
			ui.MainWindow.Show()
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
