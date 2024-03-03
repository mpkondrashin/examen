/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

page_intro.go

First installer page
*/
package main

import (
	"fmt"
	"sandboxer/pkg/globals"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	downgrade = "Downgrade"
)

type PageDowngrade struct {
	BasePage
	downgradeRadio *widget.RadioGroup
}

var _ Page = &PageDowngrade{}

func (p *PageDowngrade) Name() string {
	return "Downgrade"
}

func (p *PageDowngrade) Next(previousPage PageIndex) PageIndex {
	p.SavePrevious(previousPage)
	return pgInstallation
}

func (p *PageDowngrade) Content(win fyne.Window, installer *Installer) fyne.CanvasObject {
	titleLabel := widget.NewLabel(fmt.Sprintf(globals.AppName+" %s version is already installed. Downgrade to %s?",
		installer.config.Version, globals.Version))

	p.downgradeRadio = widget.NewRadioGroup([]string{abort, downgrade}, p.radioChanged)
	p.downgradeRadio.SetSelected(abort)
	return container.NewVBox(
		titleLabel,
		p.downgradeRadio,
	)
}

func (p *PageDowngrade) Run(win fyne.Window, installer *Installer) {

}

func (p *PageDowngrade) AquireData(installer *Installer) error {
	switch p.downgradeRadio.Selected {
	case downgrade:
		installer.config.Version = globals.Version
		return nil
	case abort:
		return ErrAbort
	}
	return nil
}

// Umbrella - segments print glue
func (p *PageDowngrade) radioChanged(s string) {
	//p.wiz.win.SetContent(p.wiz.Window())
}
