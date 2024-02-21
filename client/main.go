package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/awesome-gocui/gocui"
)

const (
	queueSize       = 8
	network         = "tcp4"
	addr            = "127.0.0.1:49000"
	maxNameLength   = 20
	maxBufferLength = 128
)

var (
	name        string
	inputEditor = gocui.EditorFunc(inputEditorFunc)
	conn        net.Conn
	send        chan []byte
	receive     chan []byte
)

func init() {
	send = make(chan []byte, queueSize)
	receive = make(chan []byte, queueSize)
	var err error
	conn, err = net.Dial(network, addr)
	if err != nil {
		log.Fatal(err.Error())
	}
	go read()
}

func main() {
	defer conn.Close()
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
	printGreeting()
	fmt.Printf("\nPlease enter your name: ")
	br := bufio.NewReader(os.Stdin)
	for {
		buf, err := br.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		if len(buf) > 1 {
			if len(buf) > maxNameLength {
				buf = buf[:maxNameLength]
			}
			name = fmt.Sprintf("%s: ", buf[:len(buf)-1])
			break
		}
		fmt.Printf("Name cannot be empty. Please enter your name: ")
	}
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(layout)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(g, -1)
			return nil
		}); err != nil {
		log.Panicln(err.Error())
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			scrollView(g, 1)
			return nil
		}); err != nil {
		log.Panicln(err.Error())
	}
	go updateView(g)
	n := fmt.Sprintf("%s\n", name)
	send <- []byte(n)
	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("output", 0, 1, maxX-1, maxY-4, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Autoscroll = true
	}
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

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func updateView(g *gocui.Gui) {
	for {
		select {
		case buf := <-receive:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("output")
				if err != nil {
					return err
				}
				v.Autoscroll = true
				fmt.Fprintf(v, "%s", buf)
				return nil
			})
		case buf := <-send:
			n, err := conn.Write(buf)
			if err != nil {
				log.Printf("Bytes Written: %d: %s\n", n, err.Error())
				// return
			}
		}
	}
}

func inputEditorFunc(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if ch != 0 && mod == 0 {
		v.EditWrite(ch)
		return
	}
	switch key {
	case gocui.KeyCtrlD:
		// DO NOTHING
	case gocui.KeySpace:
		v.EditWrite(' ')
	case gocui.KeyBackspace, gocui.KeyBackspace2:
		x, _ := v.Cursor()
		if x > len(name) {
			v.EditDelete(true)
		}
	case gocui.KeyDelete:
		x, _ := v.Cursor()
		if x > len(name) {
			v.EditDelete(false)
		}
	case gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case gocui.KeyEnter:
		if len(v.Buffer()) > len(name) {
			v.EditWrite('\n')
			out := []byte(v.Buffer())
			if len(out) > maxBufferLength {
				out = out[:maxBufferLength]
				out = append(out, '\n')
			}
			send <- out
			v.Clear()
			for _, r := range name {
				v.EditWrite(r)
			}
		}
	case gocui.KeyArrowDown:
		v.MoveCursor(0, 1)
	case gocui.KeyArrowUp:
		v.MoveCursor(0, -1)
	case gocui.KeyArrowLeft:
		x, _ := v.Cursor()
		if x > len(name) {
			v.MoveCursor(-1, 0)
		}
	case gocui.KeyArrowRight:
		v.MoveCursor(1, 0)
	case gocui.KeyTab:
		v.EditWrite(' ')
	case gocui.KeyEsc:
		// If not here the esc key will act like the gocui.KeySpace
	default:
		v.EditWrite(ch)
	}
}

func read() {
	br := bufio.NewReader(conn)
	for {
		buf, err := br.ReadBytes(0)
		if err != nil {
			// log.Println(err.Error())
			return
		}
		buf[len(buf)-1] = '\n'
		receive <- buf
		br.Reset(conn)
	}

}

func scrollView(g *gocui.Gui, dy int) {
	v, err := g.View("output")
	if err != nil {
		log.Fatal(err.Error())
	}
	// Get the size and position of the view.
	_, y := v.Size()
	ox, oy := v.Origin()
	// If we're at the bottom...
	if oy+dy > strings.Count(v.ViewBuffer(), "\n")-y-1 {
		// Set autoscroll to normal again.
		v.Autoscroll = true
	} else {
		// Set autoscroll to false and scroll.
		v.Autoscroll = false
		v.SetOrigin(ox, oy+dy)
	}
}

func printGreeting() {
	fmt.Println("Welcome to TCP-Chat!")
	fmt.Println("         _nnnn_")
	fmt.Println("        dGGGGMMb")
	fmt.Println("       @p~qp~~qMb")
	fmt.Println("       M|@||@) M|")
	fmt.Println("       @,----.JM|")
	fmt.Println("      JS^\\__/  qKL")
	fmt.Println("     dZP        qKRb")
	fmt.Println("    dZP          qKKb")
	fmt.Println("   fZP            SMMb")
	fmt.Println("   HZM            MMMM")
	fmt.Println("   FqM            MMMM")
	fmt.Println(" __| \".        |\\dS\"qML")
	fmt.Println(" |    `.       | `' \\Zq")
	fmt.Println("_)      \\.___.,|     .'")
	fmt.Println("\\____   )MMMMMP|   .'")
	fmt.Println("     `-'       `--'")
}
