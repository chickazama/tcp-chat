package ui

import "github.com/awesome-gocui/gocui"

func setInputEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	if ch != 0 && mod == 0 {
		v.EditWrite(ch)
		return
	}
	switch key {
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
			client.Send <- out
			v.Clear()
			for _, r := range name {
				v.EditWrite(r)
			}
		}
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
