package frame

import (
	"log"
	"howm/ext"
	// "github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xcursor"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	// "github.com/BurntSushi/xgbutil/xprop"
)

type PartitionType int
const (
	HORIZONTAL PartitionType = iota
	VERTICAL
)

type Rect struct {
	X int
	Y int
	W int
	H int
}

type Decoration struct {
	Shape Rect
	Window *xwindow.Window
}

type ContainerDecorations struct {
	Close, Grab Decoration
	Top, Left, Bottom, Right Decoration
	TopRight, TopLeft, BottomRight, BottomLeft Decoration
}

type DragOrigin struct {
	ContainerX int
	ContainerY int
	MouseX int
	MouseY int
}

type Container struct {
	Shape Rect
	Root *Frame
	DragContext DragOrigin
	Decorations ContainerDecorations
}

type Partition struct {
	Ratio float64
	Type PartitionType
	Decoration Decoration
}

type Frame struct {
	Shape Rect
	Window *xwindow.Window
	Container *Container
	Parent, ChildA, ChildB *Frame
	Separator Partition
}

type AttachTarget struct {
	Target *Frame
	Type PartitionType
}

type Config struct {
	ButtonClose string
	ButtonDrag string
	ElemSize int
	CloseCursor int
	DefaultShape Rect
}

type Context struct {
	X *xgbutil.XUtil
	AttachPoint *AttachTarget
	Tracked map[xproto.Window]*Frame
	Cursors map[int]xproto.Cursor
	Config Config
}

func (f *Frame) Traverse(fun func(*Frame)) {
	fun(f)
	if f.ChildA != nil {
		f.ChildA.Traverse(fun)
	}
	if f.ChildB != nil {
		f.ChildB.Traverse(fun)
	}
}

func (f *Frame) Root() *Frame {
	z := f
	for {
		if z.Parent != nil {
			z = z.Parent
		} else {
			return z
		}
	}
}

func (cd *ContainerDecorations) ForEach(f func(*Decoration)) {
	f(&cd.Close)
	f(&cd.Grab)
	f(&cd.Top)
	f(&cd.BottomRight)
}

func (d *Decoration) MoveResize(r Rect) {
	d.Window.MoveResize(r.X, r.Y, r.W, r.H)
}

func (cd *ContainerDecorations) MoveResize(ctx *Context, cShape Rect) {
	cd.Close.MoveResize(CloseShape(ctx, cShape))
	cd.Grab.MoveResize(GrabShape(ctx, cShape))
	cd.Top.MoveResize(TopShape(ctx, cShape))
	cd.BottomRight.MoveResize(BottomRightShape(ctx, cShape))
}

func (cd *ContainerDecorations) Map() {
	cd.ForEach(func(d *Decoration){
		d.Window.Map()
	})
}

func (f *Frame) Map() {
	f.Traverse(
		func(f *Frame){
			f.Window.Map()
		},
	)
}

func (f *Frame) IsLeaf() bool {
	return f.ChildA == nil && f.ChildB == nil
}

func (f *Frame) IsRoot() bool {
	return f.Parent == nil
}

func (f *Frame) MoveResize(ctx *Context) {
	f.Traverse(func(ft *Frame){
		f.Shape = f.CalcShape(ctx)
		if ft.IsLeaf() {
			ft.Window.MoveResize(ft.Shape.X, ft.Shape.Y, ft.Shape.W, ft.Shape.H)
		}
	})
}

func (c *Container) Map() {
	c.Root.Map()
	c.Decorations.Map()
}

func (c *Container) MoveResize(ctx *Context, x, y, w, h int) {
	c.Shape = Rect{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
	c.Root.MoveResize(ctx)
	c.Decorations.MoveResize(ctx, c.Shape)
}

func NewContext(x *xgbutil.XUtil) (*Context, error) {
	conf := Config{
		ButtonClose: "1",
		ButtonDrag: "1",
		ElemSize: 10,
		CloseCursor: xcursor.Dot,
		DefaultShape: Rect {
			X: 200,
			Y: 200,
			W: 800,
			H: 200,
		},
	}
	var err error
	c := &Context{
		X: x,
		Tracked: make(map[xproto.Window]*Frame),
		Cursors: make(map[int]xproto.Cursor),
		Config: conf,
	}

	for i := xcursor.XCursor; i < xcursor.XTerm; i++ {
		curs, err := xcursor.CreateCursor(x, uint16(i))
		if err != nil {
			break
		}
		c.Cursors[i] = curs
	}
	return c, err
}

func (ctx *Context) MinShape() Rect {
	return Rect{
		X: 0,
		Y: 0,
		W: 5*ctx.Config.ElemSize,
		H: 5*ctx.Config.ElemSize,
	}
}

func NewContainer(ctx *Context, ev xevent.MapRequestEvent) *Container {
	window := ev.Window
	if existing := ctx.Get(window); existing != nil {
		log.Println("NewContainer:", window, "already tracked")
		return existing.Container
	}

	root := &Frame{
		Shape: RootShape(ctx, ctx.Config.DefaultShape),
		Window: xwindow.New(ctx.X, window),
	}
	root.Window.MoveResize(root.Shape.X, root.Shape.Y, root.Shape.W, root.Shape.H)
	if err := ext.MapChecked(root.Window); err != nil {
		log.Println("NewContainer:", window, "could not be mapped")
		return nil
	}

	c := &Container{
		Shape: ctx.Config.DefaultShape,
		Root: root,
	}
	root.Container = c

	// Create Decorations
	var err error
	c.Decorations.Close, err = CreateDecoration(
		ctx,
		CloseShape(ctx, c.Shape),
		0xff0000,
		uint32(ctx.Cursors[xcursor.Dot]),
	)
	c.Decorations.Grab, err = CreateDecoration(
		ctx,
		GrabShape(ctx, c.Shape),
		0x335555,
		0,
	)
	c.Decorations.Top, err = CreateDecoration(
		ctx,
		TopShape(ctx, c.Shape),
		0x777777,
		0,
	)
	c.Decorations.BottomRight, err = CreateDecoration(
		ctx,
		BottomRightShape(ctx, c.Shape),
		0x00ff00,
		0,
	)

	// Add hooks
	err = c.AddCloseHook(ctx)
	c.AddBottomRightHook(ctx)
	c.AddGrabHook(ctx)
	AddWindowHook(ctx, window)

	if err != nil {
		log.Println("NewContainer: failed to create a decoration", err)
		return nil
	}

	// Yay
	c.Map()
	ctx.Tracked[window] = c.Root
	return c
}

func RootShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X + context.Config.ElemSize,
		Y: cShape.Y + 2*context.Config.ElemSize,
		W: cShape.W - 2*context.Config.ElemSize,
		H: cShape.H - 3*context.Config.ElemSize,
	}
}

func ContainerShapeFromRoot(ctx *Context, fShape Rect) Rect {
	return Rect{
		X: fShape.X - ctx.Config.ElemSize,
		Y: fShape.Y - 2*ctx.Config.ElemSize,
		W: fShape.W + 2*ctx.Config.ElemSize,
		H: fShape.H + 3*ctx.Config.ElemSize,
	}
}

func (f *Frame) CalcShape(ctx *Context) Rect {
	if f.IsRoot() {
		return RootShape(ctx, f.Container.Shape)
	}

	isChildA := (f.Parent.ChildA == f)

	if isChildA {
		if f.Parent.Separator.Type == HORIZONTAL {
			return Rect{
				X: f.Parent.Shape.X,
				Y: f.Parent.Shape.Y,
				W: int(float64(f.Parent.Shape.W) * f.Parent.Separator.Ratio) - ctx.Config.ElemSize,
				H: f.Parent.Shape.H,
			}
		} else {
			return Rect{
				X: f.Parent.Shape.X,
				Y: f.Parent.Shape.Y,
				W: f.Parent.Shape.W,
				H: int(float64(f.Parent.Shape.H) * f.Parent.Separator.Ratio) - ctx.Config.ElemSize,
			}
		}
	} else {
		if f.Parent.Separator.Type == HORIZONTAL {
			wr := int(float64(f.Parent.Shape.W) * f.Parent.Separator.Ratio)
			return Rect{
				X: wr,
				Y: f.Parent.Shape.Y,
				W: f.Parent.Shape.W - wr,
				H: f.Parent.Shape.H,
			}
		} else {
			hr := int(float64(f.Parent.Shape.H) * f.Parent.Separator.Ratio)
			return Rect{
				X: f.Parent.Shape.X,
				Y: hr,
				W: f.Parent.Shape.W,
				H: f.Parent.Shape.H - hr,
			}
		}
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

func GrabShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X,
		Y: cShape.Y + context.Config.ElemSize,
		W: cShape.W - context.Config.ElemSize,
		H: context.Config.ElemSize,
	}
}

func TopShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X + context.Config.ElemSize,
		Y: cShape.Y,
		W: cShape.W - 2*context.Config.ElemSize,
		H: context.Config.ElemSize,
	}
}

func CloseShape(context *Context, cShape Rect) Rect {
	return Rect{
		X: cShape.X + cShape.W - context.Config.ElemSize,
		Y: cShape.Y + context.Config.ElemSize,
		W: context.Config.ElemSize,
		H: context.Config.ElemSize,
	}
}

func AddWindowHook(ctx *Context, window xproto.Window) {
	xevent.ConfigureRequestFun(
		func(X *xgbutil.XUtil, ev xevent.ConfigureRequestEvent) {
			log.Println(ev)
			f := ctx.Get(window)
			if f != nil && f.IsRoot() && f.IsLeaf() {
				fShape := f.Shape
				if xproto.ConfigWindowX&ev.ValueMask > 0 {
					fShape.X = int(ev.X)
				}
				if xproto.ConfigWindowY&ev.ValueMask > 0 {
					fShape.Y = int(ev.Y)
				}
				if xproto.ConfigWindowWidth&ev.ValueMask > 0 {
					fShape.W = int(ev.Width)
				}
				if xproto.ConfigWindowHeight&ev.ValueMask > 0 {
					fShape.H = int(ev.Height)
				}
				cShape := ContainerShapeFromRoot(ctx, fShape)
				f.Container.MoveResize(ctx, cShape.X, cShape.Y, cShape.W, cShape.H)
			}
	}).Connect(ctx.X, window)
}

func (c *Container) AddGrabHook(ctx *Context) {
	mousebind.Drag(
		ctx.X, c.Decorations.Grab.Window.Id, c.Decorations.Grab.Window.Id, ctx.Config.ButtonDrag, true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			c.DragContext.ContainerX = c.Shape.X
			c.DragContext.ContainerY = c.Shape.Y
			c.DragContext.MouseX = rX
			c.DragContext.MouseY = rY
			// f.ToTop()
			return true, ctx.Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			dX := rX - c.DragContext.MouseX
			dY := rY - c.DragContext.MouseY
			c.MoveResize(ctx, c.DragContext.ContainerX + dX, c.DragContext.ContainerY + dY, c.Shape.W, c.Shape.H)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
		},
	)
}

func (ct *Container) AddCloseHook(c *Context) error {
	return mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
			log.Println("Close hook called")
		}).Connect(c.X, ct.Decorations.Close.Window.Id, c.Config.ButtonClose, false, true)
}

func (c *Container) AddBottomRightHook(ctx *Context) {
	mousebind.Drag(
		ctx.X, c.Decorations.BottomRight.Window.Id, c.Decorations.BottomRight.Window.Id, ctx.Config.ButtonDrag, true,
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
			c.DragContext.ContainerX = c.Shape.X
			c.DragContext.ContainerY = c.Shape.Y
			c.DragContext.MouseX = rX
			c.DragContext.MouseY = rY
			return true, ctx.Cursors[xcursor.Circle]
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
			w := ext.IMax(rX - c.DragContext.ContainerX, ctx.MinShape().W)
			h := ext.IMax(rY - c.DragContext.ContainerY, ctx.MinShape().H)
			c.MoveResize(ctx, c.DragContext.ContainerX, c.DragContext.ContainerY, w, h)
		},
		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
		},
	)
}

func CreateDecoration(c *Context, shape Rect, color uint32, cursor uint32) (Decoration, error) {
	w, err := xwindow.Generate(c.X)
	if err != nil {
		log.Println("CreateDecoration: failed to create xwindow")
		return Decoration{}, err
	}
	if cursor == 0 {
		w.CreateChecked(c.X.RootWin(), shape.X, shape.Y, shape.W, shape.H, xproto.CwBackPixel, color)
	} else {
		w.CreateChecked(c.X.RootWin(), shape.X, shape.Y, shape.W, shape.H, xproto.CwBackPixel | xproto.CwCursor, color, cursor)
	}

	return Decoration{
		Shape: Rect{
			X: shape.X,
			Y: shape.Y,
			W: shape.W,
			H: shape.H,
		},
		Window: w,
	}, nil
}

func (c *Context) Get(w xproto.Window) *Frame {
	f, _ := c.Tracked[w]
	return f
}

// OLD..........................................

// func (f *Frame) MoveResize(X* xgbutil.XUtil, x, y, w, h int) {
// 	f.X = x
// 	f.Y = y
// 	f.W = w
// 	f.H = h
// 	if f.Parent == nil {
// 		x = x + ElemSize
// 		y = y + ElemSize
// 		w = w - ElemSize
// 		h = h - ElemSize
// 		f.Close.MoveResize(f.X + f.W - ElemSize, f.Y, ElemSize, ElemSize)
// 		f.BR.MoveResize(f.X + f.W, f.Y + f.H, ElemSize, ElemSize)
// 		f.T.MoveResize(f.X + ElemSize, f.Y, f.W - 2*ElemSize, ElemSize)
// 	}

// 	if f.ChildA == nil && f.ChildB == nil {
// 		f.Window.MoveResize(x, y, w, h)
// 	}

// 	if f.ChildA != nil {
// 		if f.VerticalPartition {
// 			f.ChildA.MoveResize(X, x, y, w, h / 2)
// 		} else {
// 			f.ChildA.MoveResize(X, x, y, w / 2, h)
// 		}
// 	}

// 	if f.ChildB != nil {
// 		if f.VerticalPartition {
// 			f.ChildA.MoveResize(X, x, y + h / 2 + 1, w, h / 2 - 1)
// 		} else {
// 			f.ChildA.MoveResize(X, x + w / 2 + 1, y, w / 2 - 1, h)
// 		}
// 	}
// }

// func (f *Frame) Unmap() {
// 	f.Window.Unmap()
// 	f.Close.Unmap()
// 	f.T.Unmap()
// 	f.BR.Unmap()
// }

// func (f *Frame) Map() {
// 	f.Window.Map()
// 	f.Close.Map()
// 	f.T.Map()
// 	f.BR.Map()
// }

// func (f *Frame) Remap() {
// 	f.Unmap()
// 	f.Map()
// }

// func (f *Frame) Traverse(fun func(*Frame)) {
// 	fun(f)
// 	if f.ChildA != nil {
// 		f.ChildA.Traverse(fun)
// 	}
// 	if f.ChildB != nil {
// 		f.ChildB.Traverse(fun)
// 	}
// }

// func (f *Frame) Root() *Frame {
// 	z := f
// 	for {
// 		if z.Parent != nil {
// 			z = z.Parent
// 		} else {
// 			return z
// 		}
// 	}
// }

// func (f *Frame) ToTop() {
// 	f.Root().Traverse(func(fr *Frame){
// 		fr.Close.Stack(xproto.StackModeAbove)
// 		fr.T.Stack(xproto.StackModeAbove)
// 		fr.BR.Stack(xproto.StackModeAbove)
// 		fr.Window.Stack(xproto.StackModeAbove)
// 	})
// 	f.Window.Focus()
// }

// func NewHorizontal(X *xgbutil.XUtil, f *Frame, ev xevent.MapRequestEvent) *Frame {
// 	log.Println("Horizontal", f)
// 	nw := ev.Window
// 	if existing, ok := TrackedFrames[nw]; ok {
// 		log.Println("Already created frame for", nw)
// 		return existing
// 	}
// 	if f.ChildA == nil && f.ChildB == nil {
// 		// Make child and gift him our window
// 		f.ChildA = &Frame{
// 			X: f.X + ElemSize,
// 			Y: f.Y + ElemSize,
// 			W: f.W / 2 - ElemSize,
// 			H: f.H / 2 - ElemSize,
// 			Parent: f,
// 			Window: f.Window,
// 		}

// 		// Make the horizontal partition

// 		// Make the new child
// 		w := xwindow.New(X, nw)
// 		if err := ext.MapChecked(f.Window); err != nil {
// 			log.Fatal(err)
// 		}
// 		f.ChildB = &Frame{
// 			X: f.X + ElemSize + f.W / 2,
// 			Y: f.Y + ElemSize,
// 			W: f.W / 2,
// 			H: f.H,
// 			Parent: f,
// 			Window: w,
// 		}

// 		f.MoveResize(X, f.X, f.Y, f.W, f.H)
// 	}
// 	return f.ChildB
// }

// func NewVertical(X *xgbutil.XUtil, f *Frame, ev xevent.MapRequestEvent) *Frame  {
// 	log.Println("Vertical", f)
// 	return nil
// }

// func New(X *xgbutil.XUtil, ev xevent.MapRequestEvent) *Frame {
// 	nw := ev.Window
// 	log.Println(ev)
// 	if existing, ok := TrackedFrames[nw]; ok {
// 		log.Println("Already created frame for", nw)
// 		return existing
// 	}

// 	f := &Frame{
// 		X: 200,
// 		Y: 200,
// 		W: 800,
// 		H: 200,
// 	}

// 	f.Window = xwindow.New(X, nw)
// 	if err := ext.MapChecked(f.Window); err != nil {
// 		log.Fatal(err)
// 	}
// 	f.Window.MoveResize(f.X + ElemSize, f.Y + ElemSize, f.W - ElemSize, f.H - ElemSize)
// 	wId := f.Window.Id

// 	// Close button
// 	closew, err := xwindow.Generate(X)
//     if err != nil {
// 		log.Fatal(err)
// 	}
// 	closew.CreateChecked(X.RootWin(), f.X, f.Y, f.W, f.H, xproto.CwBackPixel | xproto.CwCursor, 0xff0000, uint32(Cursors[xcursor.Dot]))
// 	f.Close = closew

// 	// Move Bar
// 	t, err := xwindow.Generate(X)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	t.CreateChecked(X.RootWin(), f.X, f.Y, f.W, f.H, xproto.CwBackPixel, 0x777777)
// 	f.T = t

// 	xevent.DestroyNotifyFun(
// 		func(X *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
// 			f := Get(wId)
// 			log.Println("Destroy", ev)
// 			f.Close.Destroy()
// 			f.BR.Destroy()
// 			f.T.Destroy()
// 			f.Window.Destroy()
// 			delete(TrackedFrames, f.Window.Id)
// 		}).Connect(X, f.Window.Id)

// 	xevent.UnmapNotifyFun(
// 		func(X *xgbutil.XUtil, ev xevent.UnmapNotifyEvent) {
// 			f := Get(wId)
// 			log.Println("Unmap", ev)
// 			f.Close.Unmap()
// 			f.BR.Unmap()
// 			f.T.Unmap()
// 			f.Window.Unmap()
// 		}).Connect(X, f.Window.Id)
	
// 	xevent.ConfigureRequestFun(
// 		func(X *xgbutil.XUtil, ev xevent.ConfigureRequestEvent) {
// 			f := Get(wId)
// 			log.Println("Id:", f.Window.Id, ev)
// 			if f.Parent == nil {
// 				f.MoveResize(X, int(ev.X), int(ev.Y), int(ev.Width), int(ev.Height))
// 			}
// 		}).Connect(X, f.Window.Id)
	
// 	err = mousebind.ButtonPressFun(
// 		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
// 			log.Println("Clicked!", f.Close.Id)
// 			wm_protocols, err := xprop.Atm(X, "WM_PROTOCOLS")
// 			if err != nil {
// 				log.Println("xprop wm protocols failed:", err)
// 				return
// 			}
// 			wm_del_win, err := xprop.Atm(X, "WM_DELETE_WINDOW")
// 			if err != nil {
// 				log.Println("xprop delte win failed:", err)
// 				return
// 			}
// 			cm, err := xevent.NewClientMessage(32, f.Window.Id, wm_protocols, int(wm_del_win))
// 			if err != nil {
// 				log.Println("new client message failed", err)
// 				return
// 			}
// 			err = xproto.SendEventChecked(X.Conn(), false, f.Window.Id, 0, string(cm.Bytes())).Check()
// 			if err != nil {
// 				log.Println("Could not send WM_DELETE_WINDOW ClientMessage because:", err)
// 			}
// 		}).Connect(X, closew.Id, "1", false, true)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	err = mousebind.ButtonPressFun(
// 		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
// 			f := Get(wId)
// 			log.Println("Focusing window!", f.Window.Id)
// 			f.ToTop()
// 			xproto.AllowEvents(X.Conn(), xproto.AllowReplayPointer, 0)
// 		}).Connect(X, f.Window.Id, "1", true, true)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	mousebind.Drag(
// 		X, f.T.Id, f.T.Id, "1", true,
// 		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
// 			log.Println("Drag T")
// 			f.OX = f.X
// 			f.OY = f.Y
// 			f.RX = rX
// 			f.RY = rY
// 			f.ToTop()
// 			return true, Cursors[xcursor.Circle]
// 		},
// 		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
// 			dX := rX - f.RX
// 			dY := rY - f.RY
// 			f.MoveResize(X, f.OX + dX, f.OY + dY, f.W, f.H)
// 		},
// 		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
// 		},
// 	)

// 	// BR Resize
// 	br, err := xwindow.Generate(X)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	br.CreateChecked(X.RootWin(), f.X, f.Y, f.W, f.H, xproto.CwBackPixel, 0x00ff00)
// 	f.BR = br

// 	mousebind.Drag(
// 		X, f.BR.Id, f.BR.Id, "1", true,
// 		func(X *xgbutil.XUtil, rX, rY, eX, eY int) (bool, xproto.Cursor) {
// 			f.OX = f.X
// 			f.OY = f.Y
// 			f.RX = rX
// 			f.RY = rY
// 			f.ToTop()
// 			return true, Cursors[xcursor.Circle]
// 		},
// 		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
// 			w := ext.IMax(rX - f.X, ElemSize)
// 			h := ext.IMax(rY - f.Y, ElemSize)
// 			f.MoveResize(X, f.X, f.Y, w, h)
// 		},
// 		func(X *xgbutil.XUtil, rX, rY, eX, eY int) {
// 		},
// 	)

// 	f.MoveResize(X, f.X, f.Y, f.W, f.H)

// 	TrackedFrames[nw] = f
// 	defer f.Map()
// 	return f
// }