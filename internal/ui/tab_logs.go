package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/fujiNetWIFI/tnfs-gui/internal/tnfs"
)

func makeLogsTab(ui *UI) *fyne.Container {
	logText := widget.NewLabel("")
	logText.Wrapping = fyne.TextWrapBreak
	ui.On(tnfs.Log, func(e tnfs.Event) {
		lines := strings.Split(logText.Text, "\n")
		lines = append(lines, e.Data)
		if len(lines) > 100 {
			lines = lines[len(lines)-100:]
		}
		logText.SetText(strings.Join(lines, "\n"))
	})
	return container.NewStack(container.NewVScroll(logText))
}
