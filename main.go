package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/fujiNetWIFI/tnfs-gui/internal/config"
	"github.com/fujiNetWIFI/tnfs-gui/internal/tnfs"
	"github.com/fujiNetWIFI/tnfs-gui/internal/ui"
)

func main() {
	icon, _ := fyne.LoadResourceFromPath("Icon.png")
	a := app.New()
	a.SetIcon(icon)

	cfg, err := config.LoadConfig()
	if err != nil {
		ui.ShowExeNotFound()
	}

	events := make(chan tnfs.Event)
	server := tnfs.NewServer(cfg, events)

	u := ui.NewUI(cfg, server, events)
	u.ShowMain()
}
