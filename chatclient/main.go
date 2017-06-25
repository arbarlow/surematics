package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jroimartin/gocui"
)

var (
	viewArr    = []string{"input", "side"}
	active     = 0
	serverConn net.Conn
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		panic(err)
	}

	serverConn = conn

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, sendMessage); err != nil {
		log.Panicln(err)
	}

	go func(g *gocui.Gui) {
		r := bufio.NewReader(conn)
		var err error
		var s string
		for err == nil {
			s, err = r.ReadString('\r')
			if err == nil {
				fmt.Printf("%q", s)
			}
			g.Execute(func(g *gocui.Gui) error {
				m, err := g.View("main")
				if err != nil {
					return err
				}

				m.Write([]byte(s))
				return nil
			})
		}
	}(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}

func layout(g *gocui.Gui) error {
	w, h := g.Size()

	v, err := g.SetView("input", 0, h-4, w-1, h-1)

	if err != gocui.ErrUnknownView {
		return err
	}
	v.Title = "type a message..."
	v.Editable = true
	v.Wrap = true

	if _, err = setCurrentViewOnTop(g, "input"); err != nil {

		return err
	}

	v, err = g.SetView("side", 1, 0, int(0.2*float32(w)), h-5)
	if err != gocui.ErrUnknownView {
		return err
	}
	v.Title = "people online"

	if _, err := g.SetView("main", int(0.2*float32(w))+1, 0, w-1, h-5); err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

func sendMessage(g *gocui.Gui, v *gocui.View) error {
	_, err := io.Copy(serverConn, v)
	v.Rewind()
	v.Clear()
	v.SetCursor(0, 0)
	return err
}
