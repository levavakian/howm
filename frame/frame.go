package frame

import (
	// "github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xcursor"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"log"
	"howm/basics"
)

const (
	ElemSize = 10
)

type Frame struct {
	X      int
	Y      int
	W      int
	H      int
	Window *xwindow.Window
	Close, TL, TR, BL, BR, T, L, B, R *xwindow.Window

	OX, OY, RX, RY int
}

var trackedFrames = make(map[xproto.Window]*Frame)
var Cursors = make(map[int]xproto.Cursor)

func InitializeCursors(X *xgbutil.XUtil) error {
	for i := xcursor.XCursor; i <= xcursor.XTerm; i++ {
		curs, err := xcursor.CreateCursor(X, uint16(i))
		if err != nil {
			return err
		}
		Cursors[i] = curs
	}
	return nil
}

func MapChecked(w *xwindow.Window) error {
	if w == nil {
		return nil
	}
	return xproto.MapWindowChecked(w.X.Conn(), w.Id).Check()
}

func (f *Frame) MoveResize(X* xgbutil.XUtil, x, y, w, h int) {
	f.X = x
	f.Y = y
	f.W = w
	f.H = h
	f.Window.MoveResize(f.X + ElemSize, f.Y + ElemSize, f.W - ElemSize, f.H - ElemSize)
	f.Close.MoveResize(f.X + f.W - ElemSize, f.Y, ElemSize, ElemSize)
	f.BR.MoveResize(f.X + f.W, f.Y + f.H, ElemSize, ElemSize)
	f.T.MoveResize(f.X + ElemSize, f.Y, f.W - 2*ElemSize, ElemSize)
}

func (f *Frame) Unmap() {
	f.Window.Unmap()
	f.Close.Unmap()
	f.T.Unmap()
	f.BR.Unmap()
}

func (f *Frame) Map() {
	f.Window.Map()
	f.Close.Map()
	f.T.Map()
	f.BR.Map()
}

func (f *Frame) Remap() {
	f.Unmap()
	f.Map()
}

func (f *Frame) ToTop() {
	f.Close.Stack(xproto.StackModeAbove)
	f.T.Stack(xproto.StackModeAbove)
	f.BR.Stack(xproto.StackModeAbove)
	f.Window.Stack(xproto.StackModeAbove)
	f.Window.Focus()
}

func New(X *xgbutil.XUtil, nw xproto.Window) *Frame {
	log.Println("Request for new frame for", nw)
	if existing, ok := trackedFrames[nw]; ok {
		log.Println("Already created frame for", nw)
		return existing
	}

	f := &Frame{
		X: 200,
		Y: 200,
		W: 800,
		H: 200,
	}

	f.Window = xwindow.New(X, nw)
	if err := MapChecked(f.Window); err != nil {
		log.Fatal(err)
	}
	f.Window.MoveResize(f.X + ElemSize, f.Y + ElemSize, f.W - ElemSize, f.H - ElemSize)

	// Close button
	closew, err := xwindow.Generate(X)
    if err != nil {
		log.Fatal(err)
	}
	closew.CreateChecked(X.RootWin(), f.X, f.Y, f.W, f.H, xproto.CwBackPixel | xproto.CwCursor, 0xff0000, uint32(Cursors[xcursor.Dot]))
	f.Close = closew

	err = mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
			log.Println("Clicked!", f.Close.Id)
			wm_protocols, err := xprop.Atm(X, "WM_PROTOCOLS")
			if err != nil {
				log.Println("xprop wm protocols failed:", err)
				return
			}
			wm_del_win, err := xprop.Atm(X, "WM_DELETE_WINDOW")
			if err != nil {
				log.Println("xprop delte win failed:", err)
				return
			}
			cm, err := xevent.NewClientMessage(32, f.Window.Id, wm_protocols, int(wm_del_win))
			if err != nil {
				log.Println("new client message failed", err)
				return
			}
			err = xproto.SendEventChecked(X.Conn(), false, f.Window.Id, 0, string(cm.Bytes())).Check()
			if err != nil {
				log.Println("Could not send WM_DELETE_WINDOW ClientMessage because:", err)
			}
		}).Connect(X, closew.Id, "1", false, true)
	if err != nil {
		log.Fatal(err)
	}

	xevent.DestroyNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
			log.Println("Destroy", ev)
			f.Close.Destroy()
			f.BR.Destroy()
			f.T.Destroy()
			f.Window.Destroy()
			delete(trackedFrames, f.Window.Id)
		}).Connect(X, f.Window.Id)

	xevent.UnmapNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.UnmapNotifyEvent) {
			log.Println("Unmap", ev)
			f.Close.Unmap()
			f.BR.Unmap()
			f.T.Unmap()
			f.Window.Unmap()
		}).Connect(X, f.Window.Id)

	// Move Bar
	t, err := xwindow.Generate(X)
	if err != nil {
		log.Fatal(err)
	}
	t.CreateChecked(X.RootWin(), f.X, f.Y, f.W, f.H, xproto.CwBackPixel, 0x777777)
	f.T = t

	mousebind.Drag(
		X, f.T.Id, f.T.Id, "1", true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			log.Println("Drag T")
			f.OX = f.X
			f.OY = f.Y
			f.RX = rX
			f.RY = rY
			f.ToTop()
			return true, Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			dX := rX - f.RX
			dY := rY - f.RY
			f.MoveResize(X, f.OX + dX, f.OY + dY, f.W, f.H)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
		},
	)

	// BR Resize
	br, err := xwindow.Generate(X)
	if err != nil {
		log.Fatal(err)
	}
	br.CreateChecked(X.RootWin(), f.X, f.Y, f.W, f.H, xproto.CwBackPixel, 0x00ff00)
	f.BR = br

	mousebind.Drag(
		X, f.BR.Id, f.BR.Id, "1", true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			log.Println("Drag BR")
			f.OX = f.X
			f.OY = f.Y
			f.RX = rX
			f.RY = rY
			f.ToTop()
			return true, Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			w := basics.IMax(rX - f.X, ElemSize)
			h := basics.IMax(rY - f.Y, ElemSize)
			f.MoveResize(X, f.X, f.Y, w, h)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
		},
	)

	f.MoveResize(X, f.X, f.Y, f.W, f.H)

	trackedFrames[nw] = f
	defer f.Map()
	return f
}