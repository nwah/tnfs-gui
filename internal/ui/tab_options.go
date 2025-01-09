package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeOptionsTab(ui *UI) *fyne.Container {
	allowBackgroundCheck := widget.NewCheck("Run in background", func(checked bool) {
		ui.cfg.SetAllowBackground(checked)
		if checked {
			HideFromDock()
		} else {
			ShowInDock()
		}
	})
	startAtLoginCheck := widget.NewCheck("Start automatically at login", func(checked bool) {
		ui.cfg.SetStartAtLogin(checked)
	})
	allowBackgroundCheck.SetChecked(ui.cfg.AllowBackground)
	startAtLoginCheck.SetChecked(ui.cfg.StartAtLogin)

	return container.NewVBox(
		allowBackgroundCheck,
		startAtLoginCheck,
	)
}
