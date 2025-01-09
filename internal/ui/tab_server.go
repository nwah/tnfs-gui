package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/fujiNetWIFI/tnfs-gui/internal/tnfs"
	"github.com/sqweek/dialog"
)

func makeServerTab(ui *UI, server *tnfs.Server) *fyne.Container {
	dirPickerLabel := widget.NewLabel("Server Root Directory")
	dirPickerLabel.Importance = widget.HighImportance

	currentDirLabel := widget.NewLabel(ui.cfg.TnfsRootPath)
	currentDirLabel.Wrapping = fyne.TextWrapBreak

	statusLabel := widget.NewLabel("Server Status")
	statusLabel.Importance = widget.HighImportance
	statusIcon := widget.NewIcon(theme.NewColoredResource(theme.MediaRecordIcon(), theme.ColorNameError))
	statusText := widget.NewLabel("Not running")
	statusContent := container.NewHBox(statusIcon, statusText)

	dirPickerButton := widget.NewButton("Choose directory", func() {
		directory, err := dialog.Directory().Title("Folder to serve").Browse()
		if err != nil && err != dialog.ErrCancelled {
			currentDirLabel.SetText(err.Error())
		} else {
			ui.cfg.UpdateRootPath(directory)
			currentDirLabel.SetText(directory)
		}
	})

	startButton := widget.NewButton("Start Server", func() {
		go server.Start()
	})

	stopButton := widget.NewButton("Stop Server", func() {
		go server.Stop()
	})
	stopButton.Hide()

	ui.On(tnfs.StatusChange, func(e tnfs.Event) {
		var msg string
		var icon fyne.Resource

		switch server.Status {
		case tnfs.STOPPED:
			msg = "Not running"
			icon = theme.NewColoredResource(theme.MediaRecordIcon(), theme.ColorNameError)

			stopButton.Hide()
			startButton.Show()
			startButton.Enable()
		case tnfs.FAILED:
			msg = "Error"
			icon = theme.NewColoredResource(theme.WarningIcon(), theme.ColorNameWarning)

			if server != nil && server.Err != nil {
				msg += ": " + server.Err.Error()
			}
			stopButton.Hide()
			startButton.Show()
			startButton.Enable()
		case tnfs.STOPPING:
			msg = "Stopping..."
			icon = theme.NewColoredResource(theme.MediaRecordIcon(), theme.ColorNameWarning)

			stopButton.Disable()
		case tnfs.STARTING:
			msg = "Starting..."
			icon = theme.NewColoredResource(theme.MediaRecordIcon(), theme.ColorNameWarning)

			startButton.Disable()
		case tnfs.STARTED:
			msg = "Running on " + ui.cfg.Hostname
			icon = theme.NewColoredResource(theme.MediaRecordIcon(), theme.ColorNameSuccess)

			startButton.Hide()
			stopButton.Show()
			stopButton.Enable()
		}
		statusText.SetText(msg)
		statusIcon.SetResource(icon)
	})

	return container.NewVBox(
		dirPickerLabel,
		currentDirLabel,
		dirPickerButton,
		statusLabel,
		statusContent,
		startButton,
		stopButton,
	)
}
