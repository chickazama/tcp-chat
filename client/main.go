package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/awesome-gocui/gocui"
)

const (
	network = "tcp4"
	addr    = "127.0.0.1:49000"
)

var (
	wg         sync.WaitGroup
	viewEditor = gocui.EditorFunc(simpleEditor)
	conn       net.Conn
	send       chan []byte
	receive    chan []byte
)

func main() {
	send = make(chan []byte)
	receive = make(chan []byte)
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(layout)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := g.SetKeybinding("input", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	conn, err = net.Dial(network, addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()
	// wg.Add(2)
	go read()
	go updateView(g)
	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
	// wg.Wait()
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if _, err := g.SetView("output", 0, 1, maxX-1, 2*maxY/3-2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}
	if v, err := g.SetView("input", 0, 2*maxY/3, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Editable = true
		v.Editor = viewEditor
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func updateView(g *gocui.Gui) {
	defer wg.Done()
	for {
		select {
		case buf := <-receive:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("output")
				if err != nil {
					return err
				}
				fmt.Fprintf(v, "%s", buf)
				return nil
			})
		case buf := <-send:
			n, err := conn.Write(buf)
			if err != nil {
				log.Fatalf("Bytes Written: %d: %s\n", n, err.Error())
			}
		}
	}
}

func simpleEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if ch != 0 && mod == 0 {
		v.EditWrite(ch)
		return
	}
	switch key {
	case gocui.KeySpace:
		v.EditWrite(' ')
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		v.EditDelete(true)
	case gocui.KeyDelete:
		v.EditDelete(false)
	case gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case gocui.KeyEnter:
		v.EditWrite('\n')
		send <- []byte(v.Buffer())
		v.Clear()
	case gocui.KeyArrowDown:
		v.MoveCursor(0, 1)
	case gocui.KeyArrowUp:
		v.MoveCursor(0, -1)
	case gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0)
	case gocui.KeyArrowRight:
		v.MoveCursor(1, 0)
	case gocui.KeyTab:
		v.EditWrite('\t')
	case gocui.KeyEsc:
		// If not here the esc key will act like the gocui.KeySpace
	default:
		v.EditWrite(ch)
	}
}

func read() {
	defer wg.Done()
	br := bufio.NewReader(conn)
	for {
		buf, err := br.ReadBytes('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		receive <- buf
		br.Reset(conn)
	}
}
