package frame

import (
	"log"
	"time"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/wingo/prompt"
	"github.com/BurntSushi/wingo/text"
	"github.com/BurntSushi/wingo/render"
)


type Taskbar struct {
	Base Decoration
	TimeWin *xwindow.Window
	Hidden bool
	Elements map[*Container]*Element
}

type Element struct {
	Container *Container
	Window *xwindow.Window
	MinWin *xwindow.Window
	Shape Rect
}

func (e *Element) Raise(ctx *Context) {
	e.Window.Stack(xproto.StackModeAbove)
	e.MinWin.Stack(xproto.StackModeAbove)
}

func (e *Element) Lower(ctx *Context) {
	e.Window.Stack(xproto.StackModeBelow)
	e.MinWin.Stack(xproto.StackModeBelow)
}

func (e *Element) UpdateMapping(t *Taskbar) {
	if t.Hidden {
		e.MinWin.Unmap()
		e.Window.Unmap()
	} else {
		if e.Container.Hidden {
			e.MinWin.Map()
		} else {
			e.MinWin.Unmap()
		}
		e.Window.Map()
	}
}

func (e *Element) Destroy() {
	e.Window.Destroy()
	e.MinWin.Destroy()
}

func (e *Element) MoveResize(ctx *Context, shape Rect) {
	e.Shape = shape
	e.Window.MoveResize(shape.X, shape.Y, shape.W, shape.H)
	mwShape := MinWinShape(ctx, e.Shape)
	e.MinWin.MoveResize(mwShape.X, mwShape.Y, mwShape.W, mwShape.H)
}

func (t *Taskbar) UpdateContainer(ctx *Context, c *Container) {
	var err error
	elem, ok := t.Elements[c]
	if !ok {
		elem = &Element{
			Container: c,
			Shape: ElementShape(ctx, len(t.Elements), 0),
		}

		win, err := xwindow.Generate(ctx.X)
		if err != nil {
			log.Println(err)
			return
		}
		win.Create(ctx.X.RootWin(), elem.Shape.X, elem.Shape.Y, elem.Shape.W, elem.Shape.H, 0)
		elem.Window = win

		dec, err := CreateDecoration(ctx, MinWinShape(ctx, elem.Shape), ctx.Config.TaskbarMinMarkColor, 0)
		if err != nil {
			log.Println(err)
		}
		elem.MinWin = dec.Window

		elem.AddIconHooks(ctx)
		t.Elements[c] = elem
	}

	f := c.Root.Find(func(fr *Frame)bool{
		return fr.IsLeaf()
	})
	if f == nil {
		log.Println("could not find representative leaf")
		return
	}

	ximg, err := xgraphics.FindIcon(ctx.X, f.Window.Id, elem.Shape.W, elem.Shape.H)
	if err != nil {
		log.Println(err)
		return
	}

	ximg.XSurfaceSet(elem.Window.Id)
	ximg.XDraw()
	ximg.XPaint(elem.Window.Id)

	elem.UpdateMapping(t)
}

func (e *Element) AddIconHooks(ctx *Context) error {
	err := mousebind.ButtonPressFun(func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
		if ctx.Locked {
			return
		}

		e.Container.ChangeMinimizationState(ctx)
	}).Connect(ctx.X, e.Window.Id, ctx.Config.ButtonClick, false, true)
	if err != nil {
		return err
	}
	return err
}

func (t *Taskbar) RemoveContainer(ctx *Context, c *Container) {
	elem, ok := t.Elements[c]
	if !ok {
		return
	}
	elem.Window.Unmap()
	elem.MinWin.Unmap()
	elem.Destroy()

	delete(t.Elements, c)
	combWidth := ctx.Config.TaskbarElementShape.X *2 + ctx.Config.TaskbarElementShape.W
	for _, e := range(t.Elements) {
		if e.Shape.X >= elem.Shape.X {
			s := e.Shape
			s.X = s.X - combWidth
			e.MoveResize(ctx, s)
		}
	}
}

func ElementShape(ctx *Context, index int, offset int) Rect {
	tb := TaskbarShape(ctx)
	combWidth := ctx.Config.TaskbarElementShape.X *2 + ctx.Config.TaskbarElementShape.W
	return Rect{
		X: tb.X + (combWidth*index + ctx.Config.TaskbarElementShape.X) + offset,
		Y: tb.Y + ctx.Config.TaskbarElementShape.Y,
		W: ctx.Config.TaskbarElementShape.W,
		H: ctx.Config.TaskbarElementShape.H,
	}
}

func MinWinShape(ctx *Context, elementShape Rect) Rect {
	return Rect{
		X: elementShape.X,
		Y: elementShape.Y + elementShape.H,
		W: elementShape.W,
		H: ctx.Config.TaskbarMinMarkHeight,
	}
}

func NewTaskbar(ctx *Context) *Taskbar {
	t := &Taskbar{
		Elements: make(map[*Container]*Element),
	}
	var err error

	// Base background
	t.Base, err = CreateDecoration(ctx, TaskbarShape(ctx), ctx.Config.TaskbarBaseColor, 0)
	if err != nil {
		log.Fatal(err)
	}
	t.Base.Window.Map()

	// Time
	now := time.Now()
	s := TimeShape(ctx, now)
	win, err := xwindow.Generate(ctx.X)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	win.Create(ctx.X.RootWin(), s.X, s.Y, s.W, s.H, 0)
	win.Map()
	t.TimeWin = win

	// Initial render
	t.Update(ctx)
	return t
}

func (t *Taskbar) MoveResize(ctx *Context) {
	s := TaskbarShape(ctx)
	t.Base.Window.MoveResize(s.X, s.Y, s.W, s.H)
	st := TimeShape(ctx, time.Now())
	t.TimeWin.MoveResize(st.X, st.Y, st.W, st.H)
}

func (t *Taskbar) Update(ctx *Context) {
	now := time.Now()
	s := TimeShape(ctx, now)

	text.DrawText(
		t.TimeWin,
		prompt.DefaultInputTheme.Font,
		ctx.Config.TaskbarFontSize,
		render.NewColor(int(ctx.Config.TaskbarTextColor)),
		render.NewColor(int(ctx.Config.TaskbarBaseColor)),
		now.Format(ctx.Config.TaskbarTimeFormat),
	)
	t.TimeWin.MoveResize(s.X, s.Y, s.W, s.H)
}

func (t *Taskbar) Map() {
	t.Base.Window.Map()
	t.TimeWin.Map()
	for _, e := range(t.Elements) {
		e.UpdateMapping(t)
	}
}

func (t *Taskbar) Unmap() {
	t.Base.Window.Unmap()
	t.TimeWin.Unmap()
	for _, e := range(t.Elements) {
		e.UpdateMapping(t)
	}
}

func (t *Taskbar) Raise(ctx *Context) {
	t.Base.Window.Stack(xproto.StackModeAbove)
	t.TimeWin.Stack(xproto.StackModeAbove)
	for _, e := range(t.Elements) {
		e.Raise(ctx)
	}
}

func (t *Taskbar) Lower(ctx *Context) {
	t.Base.Window.Stack(xproto.StackModeBelow)
	t.TimeWin.Stack(xproto.StackModeBelow)
	for _, e := range(t.Elements) {
		e.Lower(ctx)
	}
}

func TaskbarShape(ctx *Context) Rect {
	if len(ctx.Screens) == 0 {
		return Rect{
			X: 0,
			Y: 0,
			W: ctx.Config.ElemSize,
			H: ctx.Config.ElemSize,
		}
	}
	return Rect{
		X: ctx.Screens[0].X,
		Y: ctx.Screens[0].Y + ctx.Screens[0].H - ctx.Config.TaskbarHeight,
		W: ctx.Screens[0].W,
		H: ctx.Config.TaskbarHeight,
	}
}

func TimeShape(ctx *Context, time time.Time) Rect {
	ew, eh := xgraphics.Extents(prompt.DefaultInputTheme.Font, ctx.Config.TaskbarFontSize, time.Format(ctx.Config.TaskbarTimeFormat))
	s := TaskbarShape(ctx)
	return Rect{
		X: s.X + s.W - ew - ctx.Config.TaskbarXPad,
		Y: s.Y + s.H - eh - ctx.Config.TaskbarYPad,
		W: ew,
		H: eh,
	}
}
