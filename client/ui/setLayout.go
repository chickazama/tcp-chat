package ui

import (
	"errors"

	"github.com/awesome-gocui/gocui"
)

func setLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// Reserve majority of space for the output view
	if v, err := g.SetView("output", 0, 1, maxX-1, maxY-4, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Autoscroll = true
	}
	// Reserve remaining space for input view
	if v, err := g.SetView("input", 0, maxY-3, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Editable = true
		v.Editor = inputEditor
		for _, r := range name {
			v.EditWrite(r)
		}
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}
	return nil
}
