package ui

import (
	"matthewhope/tcp-chat/client/core"

	"github.com/awesome-gocui/gocui"
)

const (
	maxNameLength   = 20
	maxBufferLength = 128
)

var (
	client      *core.Client
	name        string
	inputEditor gocui.EditorFunc
)

func init() {
	client = core.New()
	inputEditor = gocui.EditorFunc(setInputEditor)
}
