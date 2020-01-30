package frame

import (
	"log"
	"github.com/BurntSushi/wingo/prompt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xcursor"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/xgb/xinerama"
)

type Context struct {
	X *xgbutil.XUtil
	AttachPoint *AttachTarget
	Tracked map[xproto.Window]*Frame
	Cursors map[int]xproto.Cursor
	Config Config
	Screens []Rect
	LastKnownFocused xproto.Window
	SplitPrompt *prompt.Input
}

func NewContext(x *xgbutil.XUtil) (*Context, error) {
	conf := DefaultConfig()

	var Xin []xinerama.ScreenInfo
	if err := xinerama.Init(x.Conn()); err != nil {
		log.Fatal(err)
	}
	if xin, err := xinerama.QueryScreens(x.Conn()).Reply(); err != nil {
		log.Fatal(err)
	} else {
		Xin = xin.ScreenInfo
	}

	var err error
	c := &Context{
		X: x,
		Tracked: make(map[xproto.Window]*Frame),
		Cursors: make(map[int]xproto.Cursor),
		Config: conf,
	}
	for _, xi := range(Xin) {
		c.Screens = append(c.Screens, Rect{
			X: int(xi.XOrg),
			Y: int(xi.YOrg),
			W: int(xi.Width),
			H: int(xi.Height),
		})
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

func (c *Context) Get(w xproto.Window) *Frame {
	f, _ := c.Tracked[w]
	return f
}

func (ctx *Context) MinShape() Rect {
	return Rect{
		X: 0,
		Y: 0,
		W: 5*ctx.Config.ElemSize,
		H: 5*ctx.Config.ElemSize,
	}
}

func (ctx *Context) DefaultShapeForScreen(screen Rect) Rect {
	return Rect{
		X: screen.X + int(ctx.Config.DefaultShapeRatio.X * float64(screen.W)),
		Y: screen.Y + int(ctx.Config.DefaultShapeRatio.Y * float64(screen.H)),
		W: int(ctx.Config.DefaultShapeRatio.W * float64(screen.W)),
		H: int(ctx.Config.DefaultShapeRatio.H * float64(screen.H)),
	}
}

func (ctx *Context) GetScreenForShape(shape Rect) (Rect, int) {
	max_overlap := 0
	screen := ctx.Screens[0]
	for _, s := range(ctx.Screens) {
		overlap := AreaOfIntersection(shape, s)
		if overlap > max_overlap {
			max_overlap = overlap
			screen = s
		}
	}
	return screen, max_overlap
}

func (ctx *Context) GetFocusedFrame() *Frame {
	// Try getting the input focus directly
	focus, err := xproto.GetInputFocus(ctx.X.Conn()).Reply()
	if err != nil {
		return nil
	}

	found, ok := ctx.Tracked[focus.Focus]
	if ok {
		return found
	}

	// Try the parent as well
	parent, err := xwindow.New(ctx.X, focus.Focus).Parent()
	if err == nil {
		found, ok = ctx.Tracked[parent.Id]
		if ok {
			return found
		}
	}

	// Try our last known focus
	found, ok = ctx.Tracked[ctx.LastKnownFocused]
	if ok {
		return found
	}

	return nil
}