package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

func updateOutputView(g *gocui.Gui) {
	for buf := range client.Receive {
		g.Update(func(g *gocui.Gui) error {
			v, err := g.View("output")
			if err != nil {
				return err
			}
			v.Autoscroll = true
			fmt.Fprintf(v, "%s", buf)
			return nil
		})
	}
}
