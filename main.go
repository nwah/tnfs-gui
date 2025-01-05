package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

var Version string

const (
	APP_ID             = "org.fujinet.tnfsd.gui"
	VERSION            = "0.0.1"
	TNFS_ROOT_PATH_KEY = "tnfsRootPath"
)

var exePath, tnfsRootPath string
var server *TnfsServer
var ui *TnfsUi

type TnfsUi struct {
	main   fyne.Window
	server *fyne.Container
	info   *fyne.Container
	logs   *fyne.Container
}

var subscribers map[TnfsEventType][]func(TnfsEvent)

func subscribe(t TnfsEventType, f func(TnfsEvent)) {
	subscribers[t] = append(subscribers[t], f)
}

func listenToServerEvents(eventch chan TnfsEvent) {
	go func(ch chan TnfsEvent) {
		for e := range ch {
			for _, f := range subscribers[e.Type] {
				f(e)
			}
		}
	}(eventch)
}

func locateTnfsdExecutable() error {
	dir := "."
	exeName := "tnfsd"
	if runtime.GOOS == "windows" {
		exeName = "tnfsd.exe"
	}

	currentExePath, _ := os.Executable()
	if currentExePath != "" {
		dir = filepath.Dir(currentExePath)
	}

	exePath = filepath.Join(dir, exeName)
	_, err := exec.LookPath(exePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func loadDefaultRootPath() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		dirname = "."
	}
	prefs := fyne.CurrentApp().Preferences()
	tnfsRootPath = prefs.StringWithFallback(TNFS_ROOT_PATH_KEY, dirname)
}

func makeServerTab(ui *TnfsUi) *fyne.Container {
	a := fyne.CurrentApp()

	dirPickerLabel := widget.NewLabel("Server Root")
	dirPickerLabel.Importance = widget.HighImportance

	currentDirLabel := widget.NewLabel(tnfsRootPath)
	currentDirLabel.Wrapping = fyne.TextWrapBreak

	statusLabel := widget.NewLabel("Server Status")
	statusLabel.Importance = widget.HighImportance
	statusText := widget.NewLabel("🔴 Not running")

	dirPickerButton := widget.NewButton("Choose directory", func() {
		directory, err := dialog.Directory().Title("Folder to serve").Browse()
		if err != nil && err != dialog.ErrCancelled {
			currentDirLabel.SetText(err.Error())
		} else {
			tnfsRootPath = directory
			a.Preferences().SetString(TNFS_ROOT_PATH_KEY, tnfsRootPath)
			currentDirLabel.SetText(directory)
		}
	})

	startButton := widget.NewButton("Start Server", func() {
		go server.Start()
	})

	stopButton := widget.NewButton("Stop Server", func() {
		if server != nil {
			go server.Stop()
		}
	})
	stopButton.Hide()

	subscribe(statusChange, func(e TnfsEvent) {
		var msg string
		switch server.Status {
		// TODO: Don't use Emoji?
		case STOPPED:
			msg = "🔴 Not running"
			stopButton.Hide()
			startButton.Show()
			startButton.Enable()
		case FAILED:
			msg = "⚠️ Error"
			if server != nil && server.Err != nil {
				msg += ": " + server.Err.Error()
			}
			stopButton.Hide()
			startButton.Show()
			startButton.Enable()
		case STOPPING:
			msg = "🟠 Stopping..."
			stopButton.Disable()
		case STARTING:
			msg = "🟡 Starting..."
			startButton.Disable()
		case STARTED:
			msg = "🟢 Running"
			startButton.Hide()
			stopButton.Show()
			stopButton.Enable()
		}
		statusText.SetText(msg)
	})

	return container.NewVBox(
		dirPickerLabel,
		currentDirLabel,
		dirPickerButton,
		statusLabel,
		statusText,
		startButton,
		stopButton,
	)
}

func makeLogTab() *fyne.Container {
	logText := widget.NewLabel("")
	logText.Wrapping = fyne.TextWrapBreak
	subscribe(log, func(e TnfsEvent) {
		logText.SetText(logText.Text + "\n" + e.Data)
	})
	return container.NewStack(container.NewVScroll(logText))
}

func makeInfoTab() *fyne.Container {
	infoText := widget.NewRichTextFromMarkdown(`
## About
Use this program to start and stop a local [TNFS server](https://github.com/fujinetWIFI/tnfsd).

---

## FujiNet Users
1. Choose a directory on your computer where you have disk or
cassette images stored.
2. Start the server
3. Boot up your vintage computer, making your FujiNet is on
the same network as this computer.
4. On your vintage computer, add this computer's hostname or
IP address as a new host.
5. Select this computer and you will be able to browse, mount,
and boot images from this machine.

---

## Spectranet Users
This should work but is untested.

---

## Learn More
- Visit [fujinet.online](https://fujinet.online)
- Join the [FujiNet Discord](https://discord.gg/7MfFTvD)
	`)
	infoText.Wrapping = fyne.TextWrapBreak
	return container.NewStack(container.NewVScroll(infoText))
}

func makeMainWindow(ui *TnfsUi) fyne.Window {
	a := fyne.CurrentApp()
	w := a.NewWindow("TNFS Server - PRE-RELEASE VERSION")
	w.Resize(fyne.NewSize(420, 280))
	w.SetFixedSize(true)

	notice := widget.NewLabel("PRE-RELEASE VERSION - DO NOT DISTRIBUTE")
	notice.Importance = widget.DangerImportance

	w.SetContent(
		container.NewVBox(
			notice,
			container.NewAppTabs(
				container.NewTabItem("Server", ui.server),
				container.NewTabItem("Log", ui.logs),
				container.NewTabItem("Info", ui.info),
			),
		),
	)

	return w
}

func showExeNotFound() {
	a := fyne.CurrentApp()
	w := a.NewWindow("TNFS Server - PRERELEASE VERSION")
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

func main() {
	icon, _ := fyne.LoadResourceFromPath("Icon.png")
	a := app.NewWithID(APP_ID)
	a.SetIcon(icon)

	err := locateTnfsdExecutable()
	if err != nil {
		showExeNotFound()
		return
	}

	loadDefaultRootPath()

	subscribers = make(map[TnfsEventType][]func(TnfsEvent))
	ui = &TnfsUi{}

	ui.server = makeServerTab(ui)
	ui.logs = makeLogTab()
	ui.info = makeInfoTab()
	ui.main = makeMainWindow(ui)

	eventch := make(chan TnfsEvent)
	listenToServerEvents(eventch)

	server = NewTnfsServer(eventch)
	a.Lifecycle().SetOnStopped(func() {
		server.killSubprocess()
	})

	go server.findExistingProcess()

	ui.main.ShowAndRun()
	ui.main.SetMaster()
}
