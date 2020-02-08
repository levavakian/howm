package frame

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xwindow"
	"log"
)

type PartitionType int

const (
	HORIZONTAL PartitionType = iota
	VERTICAL
)

type Partition struct {
	Ratio      float64
	Type       PartitionType
	Decoration Decoration
}

type Decoration struct {
	Window *xwindow.Window
}

type ContainerDecorations struct {
	Hidden                                     bool
	Close, Maximize, Minimize, Grab            Decoration
	Top, Left, Bottom, Right                   Decoration
	TopRight, TopLeft, BottomRight, BottomLeft Decoration
}

func CreateDecoration(c *Context, shape Rect, color uint32, cursor uint32) (Decoration, error) {
	w, err := xwindow.Generate(c.X)
	if err != nil {
		log.Println("CreateDecoration: failed to create xwindow")
		return Decoration{}, err
	}
	if cursor == 0 {
		err := w.CreateChecked(c.X.RootWin(), shape.X, shape.Y, shape.W, shape.H, xproto.CwBackPixel, color)
		if err != nil {
			log.Println(err)
		}
	} else {
		err := w.CreateChecked(c.X.RootWin(), shape.X, shape.Y, shape.W, shape.H, xproto.CwBackPixel|xproto.CwCursor, color, cursor)
		if err != nil {
			log.Println(err)
		}
	}

	return Decoration{
		Window: w,
	}, nil
}

func (d *Decoration) MoveResize(r Rect) {
	d.Window.MoveResize(r.X, r.Y, r.W, r.H)
}

func (cd *ContainerDecorations) ForEach(f func(*Decoration)) {
	f(&cd.Close)
	f(&cd.Minimize)
	f(&cd.Maximize)
	f(&cd.Grab)
	f(&cd.Top)
	f(&cd.Bottom)
	f(&cd.Left)
	f(&cd.Right)
	f(&cd.BottomRight)
	f(&cd.BottomLeft)
	f(&cd.TopRight)
	f(&cd.TopLeft)
}

func (cd *ContainerDecorations) Destroy(ctx *Context) {
	cd.ForEach(func(d *Decoration) {
		d.Window.Unmap()
		d.Window.Destroy()
	})
}

func (cd *ContainerDecorations) MoveResize(ctx *Context, cShape Rect) {
	cd.Close.MoveResize(CloseShape(ctx, cShape))
	cd.Grab.MoveResize(GrabShape(ctx, cShape))
	cd.Top.MoveResize(TopShape(ctx, cShape))
	cd.Bottom.MoveResize(BottomShape(ctx, cShape))
	cd.Left.MoveResize(LeftShape(ctx, cShape))
	cd.Right.MoveResize(RightShape(ctx, cShape))
	cd.BottomRight.MoveResize(BottomRightShape(ctx, cShape))
	cd.BottomLeft.MoveResize(BottomLeftShape(ctx, cShape))
	cd.TopRight.MoveResize(TopRightShape(ctx, cShape))
	cd.TopLeft.MoveResize(TopLeftShape(ctx, cShape))
	cd.Maximize.MoveResize(MaximizeShape(ctx, cShape))
	cd.Minimize.MoveResize(MinimizeShape(ctx, cShape))
}

func (cd *ContainerDecorations) Map() {
	cd.ForEach(func(d *Decoration) {
		d.Window.Map()
	})
}

func (cd *ContainerDecorations) Unmap() {
	cd.ForEach(func(d *Decoration) {
		d.Window.Unmap()
	})
}
