package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"

	"github.com/fujiNetWIFI/tnfs-gui/internal/server"
	"github.com/fujiNetWIFI/tnfs-gui/internal/ui"
)

const (
	APP_ID             = "org.fujinet.tnfsd.gui"
	VERSION            = "0.0.1"
	TNFS_ROOT_PATH_KEY = "tnfsRootPath"
)

func main() {
	icon, _ := fyne.LoadResourceFromPath("Icon.png")
	a := app.NewWithID(APP_ID)
	a.SetIcon(icon)

	err := ui.locateTnfsdExecutable()
	if err != nil {
		ui.showExeNotFound()
		return
	}

	ui.loadDefaultRootPath()

	subscribers = make(map[server.TnfsEventType][]func(server.TnfsEvent))
	u := &ui.TnfsUi{}

	u.server = ui.makeServerTab()
	u.logs = ui.makeLogTab()
	u.info = ui.makeInfoTab()
	u.main = ui.makeMainWindow(u)

	eventch := make(chan server.TnfsEvent)
	ui.listenToServerEvents(eventch)

	s := server.NewTnfsServer(eventch)
	a.Lifecycle().SetOnStopped(func() {
		s.killSubprocess()
	})

	go s.findExistingProcess()

	u.main.ShowAndRun()
	u.main.SetMaster()
}
