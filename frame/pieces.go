package frame

import (
	"howm/ext"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xcursor"
	"github.com/BurntSushi/xgbutil/mousebind"
)

func (c *Container) AddGrabHook(ctx *Context) {
	mousebind.Drag(
		ctx.X, c.Decorations.Grab.Window.Id, c.Decorations.Grab.Window.Id, ctx.Config.ButtonDrag, true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			c.DragContext = GenerateDragContext(ctx, c, nil, rX, rY)
			c.RaiseFindFocus(ctx)
			return true, ctx.Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			dX := rX - c.DragContext.MouseX
			dY := rY - c.DragContext.MouseY
			c.MoveResize(ctx, c.DragContext.Container.X + dX, c.DragContext.Container.Y + dY, c.Shape.W, c.Shape.H)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			c.RaiseFindFocus(ctx)
		},
	)
}

func (c *Container) AddCloseHook(ctx *Context) error {
	return mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
			c.Root.Close(ctx)
		}).Connect(ctx.X, c.Decorations.Close.Window.Id, ctx.Config.ButtonClose, false, true)
}

func (c *Container) AddTopHook(ctx *Context) {
	mousebind.Drag(
		ctx.X, c.Decorations.Top.Window.Id, c.Decorations.Top.Window.Id, ctx.Config.ButtonDrag, true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			c.DragContext = GenerateDragContext(ctx, c, nil, rX, rY)
			c.RaiseFindFocus(ctx)
			return true, ctx.Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			origYEnd := c.DragContext.Container.Y + c.DragContext.Container.H
			h := ext.IMax(origYEnd - rY, ctx.MinShape().H)
			y := origYEnd - h
			c.MoveResize(ctx, c.DragContext.Container.X, y, c.DragContext.Container.W, h)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			c.RaiseFindFocus(ctx)
		},
	)
}

func (c *Container) AddBottomHook(ctx *Context) {
	mousebind.Drag(
		ctx.X, c.Decorations.Bottom.Window.Id, c.Decorations.Bottom.Window.Id, ctx.Config.ButtonDrag, true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			c.DragContext = GenerateDragContext(ctx, c, nil, rX, rY)
			c.RaiseFindFocus(ctx)
			return true, ctx.Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			h := ext.IMax(rY - c.DragContext.Container.Y, ctx.MinShape().H)
			c.MoveResize(ctx, c.DragContext.Container.X, c.DragContext.Container.Y, c.DragContext.Container.W, h)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			c.RaiseFindFocus(ctx)
		},
	)
}

func (c *Container) AddRightHook(ctx *Context) {
	mousebind.Drag(
		ctx.X, c.Decorations.Right.Window.Id, c.Decorations.Right.Window.Id, ctx.Config.ButtonDrag, true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			c.DragContext = GenerateDragContext(ctx, c, nil, rX, rY)
			c.RaiseFindFocus(ctx)
			return true, ctx.Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			w := ext.IMax(rX - c.DragContext.Container.X, ctx.MinShape().W)
			c.MoveResize(ctx, c.DragContext.Container.X, c.DragContext.Container.Y, w, c.DragContext.Container.H)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			c.RaiseFindFocus(ctx)
		},
	)
}

func (c *Container) AddLeftHook(ctx *Context) {
	mousebind.Drag(
		ctx.X, c.Decorations.Left.Window.Id, c.Decorations.Left.Window.Id, ctx.Config.ButtonDrag, true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			c.DragContext = GenerateDragContext(ctx, c, nil, rX, rY)
			c.RaiseFindFocus(ctx)
			return true, ctx.Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			origXEnd := c.DragContext.Container.X + c.DragContext.Container.W
			w := ext.IMax(origXEnd - rX, ctx.MinShape().W)
			x := origXEnd - w
			c.MoveResize(ctx, x, c.DragContext.Container.Y, w, c.DragContext.Container.H)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			c.RaiseFindFocus(ctx)
		},
	)
}

func (c *Container) AddBottomRightHook(ctx *Context) {
	mousebind.Drag(
		ctx.X, c.Decorations.BottomRight.Window.Id, c.Decorations.BottomRight.Window.Id, ctx.Config.ButtonDrag, true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			c.DragContext = GenerateDragContext(ctx, c, nil, rX, rY)
			c.RaiseFindFocus(ctx)
			return true, ctx.Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			w := ext.IMax(rX - c.DragContext.Container.X, ctx.MinShape().W)
			h := ext.IMax(rY - c.DragContext.Container.Y, ctx.MinShape().H)
			c.MoveResize(ctx, c.DragContext.Container.X, c.DragContext.Container.Y, w, h)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			c.RaiseFindFocus(ctx)
		},
	)
}

func (c *Container) AddBottomLeftHook(ctx *Context) {
	mousebind.Drag(
		ctx.X, c.Decorations.BottomLeft.Window.Id, c.Decorations.BottomLeft.Window.Id, ctx.Config.ButtonDrag, true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			c.DragContext = GenerateDragContext(ctx, c, nil, rX, rY)
			c.RaiseFindFocus(ctx)
			return true, ctx.Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			origXEnd := c.DragContext.Container.X + c.DragContext.Container.W
			w := ext.IMax(origXEnd - rX, ctx.MinShape().W)
			x := origXEnd - w
			h := ext.IMax(rY - c.DragContext.Container.Y, ctx.MinShape().H)
			c.MoveResize(ctx, x, c.DragContext.Container.Y, w, h)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			c.RaiseFindFocus(ctx)
		},
	)
}

func (c *Container) AddTopRightHook(ctx *Context) {
	mousebind.Drag(
		ctx.X, c.Decorations.TopRight.Window.Id, c.Decorations.TopRight.Window.Id, ctx.Config.ButtonDrag, true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			c.DragContext = GenerateDragContext(ctx, c, nil, rX, rY)
			c.RaiseFindFocus(ctx)
			return true, ctx.Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			origYEnd := c.DragContext.Container.Y + c.DragContext.Container.H
			w := ext.IMax(rX - c.DragContext.Container.X, ctx.MinShape().W)
			h := ext.IMax(origYEnd - rY, ctx.MinShape().H)
			y := origYEnd - h
			c.MoveResize(ctx, c.DragContext.Container.X, y, w, h)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			c.RaiseFindFocus(ctx)
		},
	)
}

func (c *Container) AddTopLeftHook(ctx *Context) {
	mousebind.Drag(
		ctx.X, c.Decorations.TopLeft.Window.Id, c.Decorations.TopLeft.Window.Id, ctx.Config.ButtonDrag, true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			c.DragContext = GenerateDragContext(ctx, c, nil, rX, rY)
			c.RaiseFindFocus(ctx)
			return true, ctx.Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			origYEnd := c.DragContext.Container.Y + c.DragContext.Container.H
			origXEnd := c.DragContext.Container.X + c.DragContext.Container.W
			w := ext.IMax(origXEnd - rX, ctx.MinShape().W)
			h := ext.IMax(origYEnd - rY, ctx.MinShape().H)
			y := origYEnd - h
			x := origXEnd - w
			c.MoveResize(ctx, x, y, w, h)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			c.RaiseFindFocus(ctx)
		},
	)
}

func TopRightShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X + cShape.W - context.Config.ElemSize,
		Y: cShape.Y + context.Config.ElemSize,
		W: context.Config.ElemSize,
		H: context.Config.ElemSize,
	}
}

func TopLeftShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X,
		Y: cShape.Y + context.Config.ElemSize,
		W: context.Config.ElemSize,
		H: context.Config.ElemSize,
	}
}

func BottomRightShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X + cShape.W - context.Config.ElemSize,
		Y: cShape.Y + cShape.H - context.Config.ElemSize,
		W: context.Config.ElemSize,
		H: context.Config.ElemSize,
	}
}

func BottomLeftShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X,
		Y: cShape.Y + cShape.H - context.Config.ElemSize,
		W: context.Config.ElemSize,
		H: context.Config.ElemSize,
	}
}

func TopShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X + context.Config.ElemSize,
		Y: cShape.Y + context.Config.ElemSize,
		W: cShape.W - 2*context.Config.ElemSize,
		H: context.Config.ElemSize,
	}
}

func BottomShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X + context.Config.ElemSize,
		Y: cShape.Y + cShape.H - context.Config.ElemSize,
		W: cShape.W - 2*context.Config.ElemSize,
		H: context.Config.ElemSize,
	}
}

func LeftShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X,
		Y: cShape.Y + 2*context.Config.ElemSize,
		W: context.Config.ElemSize,
		H: cShape.H - 3*context.Config.ElemSize,
	}
}

func RightShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X + cShape.W - context.Config.ElemSize,
		Y: cShape.Y + 2*context.Config.ElemSize,
		W: context.Config.ElemSize,
		H: cShape.H - 3*context.Config.ElemSize,
	}
}

func GrabShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X,
		Y: cShape.Y,
		W: cShape.W - context.Config.ElemSize,
		H: context.Config.ElemSize,
	}
}

func CloseShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X + cShape.W - context.Config.ElemSize,
		Y: cShape.Y,
		W: context.Config.ElemSize,
		H: context.Config.ElemSize,
	}
}