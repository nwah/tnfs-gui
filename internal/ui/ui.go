package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/fujiNetWIFI/tnfs-gui/internal/config"
	"github.com/fujiNetWIFI/tnfs-gui/internal/tnfs"
)

type UI struct {
	MainWindow fyne.Window
	ServerTab  *fyne.Container
	InfoTab    *fyne.Container
	LogsTab    *fyne.Container

	subscribers map[tnfs.EventType][]func(tnfs.Event)
	cfg         *config.Config
}

func (ui *UI) On(t tnfs.EventType, f func(tnfs.Event)) {
	ui.subscribers[t] = append(ui.subscribers[t], f)
}

func (ui *UI) Listen(eventch chan tnfs.Event) {
	go func(ch chan tnfs.Event) {
		for e := range ch {
			for _, f := range ui.subscribers[e.Type] {
				f(e)
			}
		}
	}(eventch)
}

func (ui *UI) ShowMain() {
	ui.MainWindow.ShowAndRun()
	ui.MainWindow.SetMaster()
}

func NewUI(cfg *config.Config, server *tnfs.Server, ch chan tnfs.Event) *UI {
	ui := &UI{
		cfg:         cfg,
		subscribers: make(map[tnfs.EventType][]func(tnfs.Event)),
	}

	ui.ServerTab = makeServerTab(ui, server)
	ui.LogsTab = makeLogsTab(ui)
	ui.InfoTab = makeInfoTab()
	ui.MainWindow = makeMainWindow(ui)

	ui.Listen(ch)

	return ui
}

func makeMainWindow(ui *UI) fyne.Window {
	a := fyne.CurrentApp()
	w := a.NewWindow("TNFS Server Manager")
	w.Resize(fyne.NewSize(420, 280))
	w.SetFixedSize(true)

	w.SetContent(
		container.NewVBox(
			container.NewAppTabs(
				container.NewTabItem("Server", ui.ServerTab),
				container.NewTabItem("Log", ui.LogsTab),
				container.NewTabItem("Info", ui.InfoTab),
			),
		),
	)

	return w
}

func ShowExeNotFound() {
	a := fyne.CurrentApp()
	w := a.NewWindow("TNFS Server Manager")
	w.Resize(fyne.NewSize(420, 140))
	w.SetFixedSize(true)
	button := widget.NewButton("Close", func() {
		a.Quit()
	})
	text := widget.NewLabel("Cannot find tnfsd executable. Please check that it is in the same folder as this program and try again.")
	text.Wrapping = fyne.TextWrapWord
	w.SetContent(container.NewVBox(text, button))
	w.ShowAndRun()
}
