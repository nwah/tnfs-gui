package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/fujiNetWIFI/tnfs-gui/internal/tnfs"
)

func makeLogsTab(ui *UI) *fyne.Container {
	logText := widget.NewLabel("")
	logText.Wrapping = fyne.TextWrapBreak
	ui.On(tnfs.Log, func(e tnfs.Event) {
		logText.SetText(logText.Text + "\n" + e.Data)
	})
	return container.NewStack(container.NewVScroll(logText))
}
