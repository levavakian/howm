package main

import (
	"howm/frame"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/keybind"
)

type AnchorTo struct {
	Anchor int
	Screen frame.Rect
}

func RegisterTaskbarHooks(ctx *frame.Context) error {
	var err error
	// Toggle taskbar
    err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent){
        if ctx.Locked {
            return
		}
		
		options := []int{
			frame.FULL,
			frame.TOP,
			frame.LEFT,
			frame.RIGHT,
			frame.BOTTOM,
		}
		updates := make(map[*frame.Container]AnchorTo)
		for c, _ := range(ctx.Containers) {
			anchor := func()AnchorTo{
				screen, _, index := ctx.GetScreenForShape(c.Shape)
				if index != 0 {
					return AnchorTo{frame.NONE, frame.Rect{}}
				}
				for _, opt := range(options) {
					if c.Shape == frame.AnchorShape(ctx, screen, opt) {
						return AnchorTo{opt, screen}
					}
				}
				return AnchorTo{frame.NONE, frame.Rect{}}
			}()
			if anchor.Anchor != frame.NONE {
				updates[c] = anchor
			}
		}

		ctx.Taskbar.Hidden = !ctx.Taskbar.Hidden
		ctx.Taskbar.UpdateMapping(ctx)

		for c, anchor := range(updates) {
			c.MoveResizeShape(ctx, frame.AnchorShape(ctx, anchor.Screen, anchor.Anchor))
		}
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.ToggleTaskbar, true)
	if err != nil {
		return err
	}

	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent){
        if ctx.Locked {
            return
		}
		
		ctx.Taskbar.Scroller.SlideLeft(ctx)
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.TaskbarSlideLeft, true)
	if err != nil {
		return err
	}

	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent){
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
