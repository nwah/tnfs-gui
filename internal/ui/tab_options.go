package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeOptionsTab(ui *UI) *fyne.Container {
	readOnlyCheck := widget.NewCheck("Read-only mode", func(checked bool) {
		ui.cfg.SetReadOnly(checked)
	})
	allowBackgroundCheck := widget.NewCheck("Run in background", func(checked bool) {
		ui.cfg.SetAllowBackground(checked)
	})
	startAtLoginCheck := widget.NewCheck("Start automatically at login", func(checked bool) {
		ui.cfg.SetStartAtLogin(checked)
	})
	readOnlyCheck.SetChecked(ui.cfg.ReadOnly)
	allowBackgroundCheck.SetChecked(ui.cfg.AllowBackground)
	startAtLoginCheck.SetChecked(ui.cfg.StartAtLogin)

	return container.NewVBox(
		readOnlyCheck,
		allowBackgroundCheck,
		startAtLoginCheck,
	)
}
