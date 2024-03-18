package client

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/rocket049/gocui"
	"github.com/zhaommmmomo/zim/common/sdk"
	"log"
)

func InitCui() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = false
	g.ASCII = false

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("main", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, updateInput); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyPgup, gocui.ModNone, viewUpScroll); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("main", gocui.KeyPgdn, gocui.ModNone, viewDownScroll); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if err := tips(g, 1, 1, maxX-1, 3); err != nil {
		return err
	}
	if err := viewOutput(g, 1, 4, maxX-1, maxY-4); err != nil {
		return err
	}
	if err := viewInput(g, 1, maxY-3, maxX-1, maxY-1); err != nil {
		return err
	}
	return nil
}

func tips(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView("tips", x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Wrap = false
		v.Overwrite = true
		v.Title = "Tips"
		fmt.Fprint(v, color.FgGreen.Text("welcome to zim!"))
	}
	return nil
}

func viewOutput(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView("out", x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Wrap = true
		v.Overwrite = false
		v.Autoscroll = true
		v.SelBgColor = gocui.ColorRed
		v.Title = "Messages"
	}
	return nil
}

func viewInput(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView("main", x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Editable = true
		v.Wrap = true
		v.Overwrite = false
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

func updateInput(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	if cv != nil && err == nil {
		var p = cv.ReadEditor()
		if p != nil {
			v.Autoscroll = true
			v.Write([]byte("你:"))
			v.Write(append(p, '\n'))
			sdk.Send(sdk.Message{Content: p})
		}
	}
	l := len(cv.Buffer())
	cv.MoveCursor(0-l, 0, true)
	cv.Clear()
	return nil
}

func updateOutput(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	if cv != nil && err == nil {
		var p = cv.ReadEditor()
		if p != nil {
			v.Autoscroll = true
			v.Write([]byte("你:"))
			v.Write(append(p, '\n'))
			sdk.Send(sdk.Message{Content: p})
		}
	}
	l := len(cv.Buffer())
	cv.MoveCursor(0-l, 0, true)
	cv.Clear()
	return nil
}

func viewUpScroll(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	if cv != nil && err == nil {
		x, y := v.Origin()
		v.Autoscroll = false
		v.SetOrigin(x, y-1)
	}
	return nil
}

func viewDownScroll(g *gocui.Gui, cv *gocui.View) error {
	v, err := g.View("out")
	if cv != nil && err == nil {
		x, y := v.Origin()
		v.Autoscroll = false
		v.SetOrigin(x, y+1)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
