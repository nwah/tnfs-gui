package main

import (
	"os"

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

	if len(os.Args) > 1 && os.Args[1] == "autorun" {
		go server.Start()
		if cfg.AllowBackground {
			ui.HideFromDock()
			a.Run()
		} else {
			u.ShowMain()
		}
	} else {
		u.ShowMain()
	}
}
