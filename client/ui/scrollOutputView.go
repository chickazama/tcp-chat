package ui

import (
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"
)

func scrollOutputView(g *gocui.Gui, dy int) {
	v, err := g.View("output")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, y := v.Size()
	ox, oy := v.Origin()
	if oy+dy > strings.Count(v.ViewBuffer(), "\n")-y-1 {
		v.Autoscroll = true
	} else {
		v.Autoscroll = false
		v.SetOrigin(ox, oy+dy)
	}
}
