package ui

import "github.com/awesome-gocui/gocui"

func setKeyBindings(g *gocui.Gui) error {
	// Handle CTRL+C & CTRL+D as quit commands
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, quit); err != nil {
		return err
	}
	// Handle Up & Down Arrows as scroll commands
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(g, -1)
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(g, 1)
			return nil
		}); err != nil {
		return err
	}
	return nil
}
