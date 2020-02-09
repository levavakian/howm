package root

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/levavakian/rowm/frame"
)

type AnchorTo struct {
	Anchor frame.AnchorType
	Screen frame.Rect
}

func RegisterTaskbarHooks(ctx *frame.Context) error {
	var err error
	// Toggle taskbar
	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		if ctx.Locked {
			return
		}

		updates := make(map[*frame.Container]AnchorTo)
		for c, _ := range ctx.Containers {
			anchor := func() AnchorTo {
				screen, _, index := ctx.GetScreenForShape(c.Shape)
				if index != 0 {
					return AnchorTo{frame.NONE, frame.Rect{}}
				}
				if opt := frame.AnchorMatch(ctx, screen, c.Shape); opt == frame.NONE {
					return AnchorTo{frame.NONE, frame.Rect{}}
				} else {
					return AnchorTo{opt, screen}
				}
			}()
			if anchor.Anchor != frame.NONE {
				updates[c] = anchor
			}
		}

		ctx.Taskbar.Hidden = !ctx.Taskbar.Hidden
		ctx.Taskbar.UpdateMapping(ctx)

		for c, anchor := range updates {
			c.MoveResizeShape(ctx, frame.AnchorShape(ctx, anchor.Screen, anchor.Anchor))
		}
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.ToggleTaskbar, true)
	if err != nil {
		return err
	}

	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		if ctx.Locked {
			return
		}

		ctx.Taskbar.Scroller.SlideLeft(ctx)
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.TaskbarSlideLeft, true)
	if err != nil {
		return err
	}

	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		if ctx.Locked {
			return
		}

		ctx.Taskbar.Scroller.SlideRight(ctx)
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.TaskbarSlideRight, true)
	if err != nil {
		return err
	}

	return err
}
