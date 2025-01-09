package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeInfoTab() *fyne.Container {
	a := fyne.CurrentApp()
	versionText := widget.NewRichTextWithText("Version " + a.Metadata().Version)

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
	return container.NewStack(
		container.NewVScroll(
			container.NewVBox(
				versionText,
				infoText,
			),
		),
	)
}
