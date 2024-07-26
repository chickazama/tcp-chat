package ui

import (
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

func Run() {
	// cmd := exec.Command("clear")
	// cmd.Stdout = os.Stdout
	// err := cmd.Run()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	printGreeting()
	readName()
	// Set up GUI
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer g.Close()
	g.SetManagerFunc(setLayout)
	err = setKeyBindings(g)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Spawn a go-routine to update the view
	go updateOutputView(g)
	initalOut := fmt.Sprintf("%s\n", name)
	client.Send <- []byte(initalOut)
	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Fatal(err.Error())
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
