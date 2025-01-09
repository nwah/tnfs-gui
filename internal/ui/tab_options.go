package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeOptionsTab(ui *UI) *fyne.Container {
	allowBackgroundCheck := widget.NewCheck("Allow TNFS server to keep running in background", func(checked bool) {
		ui.cfg.SetAllowBackground(checked)
	})
	startAtLoginCheck := widget.NewCheck("Start automatically at login", func(checked bool) {
		ui.cfg.SetAllowBackground(checked)
	})
	allowBackgroundCheck.SetChecked(ui.cfg.AllowBackground)
	startAtLoginCheck.SetChecked(ui.cfg.StartAtLogin)

	return container.NewVBox(
		allowBackgroundCheck,
		startAtLoginCheck,
	)
}
