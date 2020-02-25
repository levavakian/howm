package root

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/levavakian/rowm/ext"
	"github.com/levavakian/rowm/frame"
	"time"
	"log"
)

func RegisterBaseHooks(ctx *frame.Context) error {
	var err error

	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		xevent.Quit(ctx.X)
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.Shutdown, true)
	if err != nil {
		return err
	}

	err = mousebind.ButtonPressFun(func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
		if ctx.Locked {
			return
		}
		ext.Focus(xwindow.New(ctx.X, ctx.X.RootWin()))
		xproto.AllowEvents(ctx.X.Conn(), xproto.AllowReplayPointer, 0)
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.ButtonClick, false, false)
	if err != nil {
		return err
	}

	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, ev xevent.KeyReleaseEvent) {
		ctx.SetLocked(true)
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.Lock, true)
	if err != nil {
		return err
	}

	focusNext := func(X *xgbutil.XUtil, ev xevent.KeyReleaseEvent, reverse bool){
		if ctx.Locked {
			return
		}
		ffoc := ctx.GetFocusedFrame()
		if ffoc == nil || ffoc.IsOrphan() {
			return
		}
		nfoc := ffoc.FindNextLeaf(func(nf *frame.Frame)bool{ return nf.IsLeaf() && nf.Mapped }, reverse, ffoc.Container.ActiveRoot())
		if nfoc != nil {
			nfoc.Container.Raise(ctx)
			nfoc.Focus(ctx)

			// Get the center and show a little square to mark where the thing is
			if ctx.FocusMarker != nil {
				ctx.FocusMarker.Unmap()
				ctx.FocusMarker.Destroy()
			}

			fshape := nfoc.CalcShape(ctx)
			dshape := frame.Rect{
				X: fshape.X + (fshape.W / 2) - (ctx.Config.ElemSize / 2),
				Y: fshape.Y + (fshape.H / 2) - (ctx.Config.ElemSize / 2),
				W: ctx.Config.ElemSize,
				H: ctx.Config.ElemSize,
			}
			decoration, err := frame.CreateDecoration(ctx, dshape, ctx.Config.FocusColor, 0)
			decoration.Window.Map()
			if err != nil {
				log.Println(err)
				return
			}
			ctx.FocusMarker = decoration.Window
			go func(){
				time.Sleep(ctx.Config.FocusMarkerTime)
				ctx.Injector.Do(func() {
						if ctx.FocusMarker != decoration.Window {
							return
						}
						decoration.Window.Unmap()
						decoration.Window.Destroy()
					},
				)
			}()
		}
	}
	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, ev xevent.KeyReleaseEvent) {
		focusNext(X, ev, false)
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.FocusNext, true)
	if err != nil {
		return err
	}
	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, ev xevent.KeyReleaseEvent) {
		focusNext(X, ev, true)
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.FocusPrev, true)
	if err != nil {
		return err
	}

	xevent.MapRequestFun(func(X *xgbutil.XUtil, ev xevent.MapRequestEvent) {
		frame.NewWindow(ctx, ev.Window)
		ctx.RaiseLock()
	}).Connect(ctx.X, ctx.X.RootWin())
	return nil
}
