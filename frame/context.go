package frame

import (
	"log"
	"howm/ext"
	"github.com/BurntSushi/wingo/prompt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xcursor"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Context struct {
	X *xgbutil.XUtil
	AttachPoint *AttachTarget
	Tracked map[xproto.Window]*Frame
	Containers map[*Container]struct{}
	Cursors map[int]xproto.Cursor
	Backgrounds map[xproto.Window]*xwindow.Window
	Config Config
	Screens []Rect
	LastKnownFocused xproto.Window
	LastKnownFocusedScreen int
	SplitPrompt *prompt.Input
}

func NewContext(x *xgbutil.XUtil) (*Context, error) {
	conf := DefaultConfig()

	var err error
	c := &Context{
		X: x,
		Tracked: make(map[xproto.Window]*Frame),
		Cursors: make(map[int]xproto.Cursor),
		Containers: make(map[*Container]struct{}),
		Config: conf,
	}
	c.UpdateScreens()
	if err != nil {
		log.Fatal(err)
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

func (ctx *Context) DetectScreensChange() (bool, []Rect, error) {
	var Xin []xinerama.ScreenInfo
	if xin, err := xinerama.QueryScreens(ctx.X.Conn()).Reply(); err != nil {
		return false, nil, err
	} else {
		Xin = xin.ScreenInfo
	}

	screens := make([]Rect, 0, len(Xin))
	for _, xi := range(Xin) {
		screens = append(screens, Rect{
			X: int(xi.XOrg),
			Y: int(xi.YOrg),
			W: int(xi.Width),
			H: int(xi.Height),
		})
	}

	if len(screens) != len(ctx.Screens) {
		return true, screens, nil
	}

	for i := 0; i < len(screens); i++ {
		if screens[i] != ctx.Screens[i] {
			return true, screens, nil
		}
	}
	return false, ctx.Screens, nil
}

func (ctx *Context) UpdateScreens() {
	changed, screens, err := ctx.DetectScreensChange()
	ext.Logerr(err)
	if err != nil || !changed {
		return
	}
    log.Println("found", len(screens), "screen(s)", screens)

    ctx.Screens = screens
  	GenerateBackgrounds(ctx)
    for c, _ := range(ctx.Containers) {
	  topShape := TopShape(ctx, c.Shape)
      if screen, overlap, _ := ctx.GetScreenForShape(topShape); topShape.Area() > overlap {
        c.MoveResizeShape(ctx, ctx.DefaultShapeForScreen(screen))
      }
	}
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

func (ctx *Context) GetScreenForShape(shape Rect) (Rect, int, int) {
	max_overlap := 0
	max_i := 0
	screen := ctx.Screens[0]
	for i, s := range(ctx.Screens) {
		overlap := AreaOfIntersection(shape, s)
		if overlap > max_overlap {
			max_overlap = overlap
			max_i = i
			screen = s
		}
	}
	return screen, max_overlap, max_i
}

func (ctx *Context) LastFocusedScreen() Rect {
	if len(ctx.Screens) <= ctx.LastKnownFocusedScreen {
		ctx.LastKnownFocusedScreen = 0
	}
	return ctx.Screens[ctx.LastKnownFocusedScreen]
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