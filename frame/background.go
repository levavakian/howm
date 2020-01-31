package frame

import (
	"log"
	"bytes"
	"image"
	"howm/ext"
	"github.com/disintegration/imaging"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/xgbutil/gopher"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/ewmh"
)

func GenerateBackgrounds(ctx *Context) error {
	for _, v := range(ctx.Backgrounds) {
		v.Unmap()
		v.Destroy()
	}
	ctx.Backgrounds = make(map[xproto.Window]*xwindow.Window)

	img, err := func()(image.Image, error){
		img, err := imaging.Open(ctx.Config.BackgroundImagePath)
		ext.Logerr(err)
		if err == nil {
			return img, err
		}

		dimg, _, err := image.Decode(bytes.NewBuffer(gopher.GopherPng()))
		ext.Logerr(err)
		return dimg, err
	}()
	if err != nil {
		return err
	}

	for _, screen := range ctx.Screens {
		rimg := imaging.Fill(img, int(screen.W), int(screen.H), imaging.Center, imaging.Lanczos)
		ximg := xgraphics.NewConvert(ctx.X, rimg)
		win := DisplayBackground(ximg, int(screen.X), int(screen.Y))
		if win != nil {
			ctx.Backgrounds[win.Id] = win
		}
	}
	return nil
}

// Modified from github.com/BurntSushi/xgbutil/xgraphics XShowExtra
func DisplayBackground(im *xgraphics.Image, x, y int) *xwindow.Window {
	quit := true
	name := "Background Image Window"
	w, h := im.Rect.Dx(), im.Rect.Dy()

	win, err := xwindow.Generate(im.X)
	if err != nil {
		log.Printf("Could not generate new window id: %s", err)
		return nil
	}

	// Create a very simple window with dimensions equal to the image.
	win.Create(im.X.RootWin(), x, y, w, h, 0)

	// Make this window close gracefully.
	win.WMGracefulClose(func(w *xwindow.Window) {
		xevent.Detach(w.X, w.Id)
		keybind.Detach(w.X, w.Id)
		mousebind.Detach(w.X, w.Id)
		w.Destroy()

		if quit {
			xevent.Quit(w.X)
		}
	})

	// Set WM_STATE so it is interpreted as a top-level window.
	err = icccm.WmStateSet(im.X, win.Id, &icccm.WmState{
		State: icccm.StateNormal,
	})
	if err != nil { // not a fatal error
		log.Printf("Could not set WM_STATE: %s", err)
	}

	// Set WM_NORMAL_HINTS so the window can't be resized.
	err = icccm.WmNormalHintsSet(im.X, win.Id, &icccm.NormalHints{
		Flags:     icccm.SizeHintPMinSize | icccm.SizeHintPMaxSize,
		MinWidth:  uint(w),
		MinHeight: uint(h),
		MaxWidth:  uint(w),
		MaxHeight: uint(h),
	})
	if err != nil { // not a fatal error
		log.Printf("Could not set WM_NORMAL_HINTS: %s", err)
	}

	// Set _NET_WM_NAME so it looks nice.
	err = ewmh.WmNameSet(im.X, win.Id, name)
	if err != nil { // not a fatal error
		log.Printf("Could not set _NET_WM_NAME: %s", err)
	}

	// Paint our image before mapping.
	im.XSurfaceSet(win.Id)
	im.XDraw()
	im.XPaint(win.Id)

	// Now we can map, since we've set all our properties.
	// (The initial map is when the window manager starts managing.)
	win.Map()
	win.Stack(xproto.StackModeBelow)

	return win
}