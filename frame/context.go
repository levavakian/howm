package frame

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/wingo/misc"
	"github.com/BurntSushi/wingo/prompt"
	"github.com/BurntSushi/xgb/xinerama"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xcursor"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"howm/ext"
	"howm/sideloop"
	"log"
	"time"
	"os/exec"
	"os/user"
)

var NoFont = xgraphics.MustFont(xgraphics.ParseFont(
	bytes.NewBuffer(misc.DataFile("write-your-password-with-this-font.ttf"))))

type AttachTarget struct {
	Target *Frame
	Type   PartitionType
}

type Yank struct {
	Container *Container
	Window    xproto.Window
}

type Context struct {
	X                      *xgbutil.XUtil
	AttachPoint            *AttachTarget
	Yanked                 *Yank
	Tracked                map[xproto.Window]*Frame
	UnmapCounter           map[xproto.Window]int
	Containers             map[*Container]struct{}
	Cursors                map[int]xproto.Cursor
	DummyIcon              *xgraphics.Image
	Backgrounds            map[xproto.Window]*xwindow.Window
	Config                 Config
	Screens                []Rect
	LastKnownFocused       xproto.Window
	LastKnownFocusedScreen int
	SplitPrompt            *prompt.Input
	Locked                 bool
	LockPrompt             *prompt.Input
	Taskbar                *Taskbar
	LastLockChange         time.Time
	Injector               *sideloop.Injector
}

func NewContext(x *xgbutil.XUtil, inj *sideloop.Injector) (*Context, error) {
	conf := DefaultConfig()

	var err error
	c := &Context{
		X:            x,
		Tracked:      make(map[xproto.Window]*Frame),
		UnmapCounter: make(map[xproto.Window]int),
		Cursors:      make(map[int]xproto.Cursor),
		Containers:   make(map[*Container]struct{}),
		Config:       conf,
		DummyIcon:    xgraphics.New(x, conf.TaskbarElementShape.ToImageRect()),
		LastLockChange: time.Now(),
		Injector:      inj,
	}
	c.UpdateScreens()
	c.Taskbar = NewTaskbar(c)
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

func (ctx *Context) GenerateLockPrompt() {
	theme := *prompt.DefaultInputTheme
	theme.Font = NoFont
	ctx.LockPrompt = prompt.NewInput(ctx.X, &theme, prompt.DefaultInputConfig)

	canc := func(inp *prompt.Input) {
		ctx.RaiseLock()
	}

	resp := func(inp *prompt.Input, text string) {
		usr, err := user.Current()
		if err != nil {
			log.Println(err)
			return
		}
		err = exec.Command(ctx.Config.Shell, "-c", fmt.Sprintf("echo %s | sudo -u %s -S", usr.Name, text)).Run()
		if err != nil {
			ctx.SetLocked(false)
		} else {
			log.Println(err)
		}

		ctx.LockPrompt.Destroy()
		ctx.LockPrompt = nil
		if ctx.Locked {
			ctx.RaiseLock()
		} else {
			ctx.LowerLock()
		}
	}
	ctx.LockPrompt.Show(ctx.Screens[0].ToXRect(), "", resp, canc)
}

func (ctx *Context) RaiseLock() {
	if !ctx.Locked {
		return
	}

	for _, bg := range ctx.Backgrounds {
		bg.Stack(xproto.StackModeAbove)
	}

	ctx.GenerateLockPrompt()
}

func (ctx *Context) LowerLock() {
	if ctx.Locked {
		return
	}

	for _, bg := range ctx.Backgrounds {
		bg.Stack(xproto.StackModeBelow)
	}
}

func (ctx *Context) SetLocked(state bool) {
	if state == ctx.Locked {
		return
	}
	ctx.Locked = state
	if ctx.Locked {
		cmd := exec.Command("bash", "-c", ctx.Config.SuspendCommand)
		err := cmd.Start()
		if err != nil {
			log.Println(err)
		}
		go func() {
			cmd.Wait()
		}()
	}
}

func (ctx *Context) DetectScreensChange() (bool, []Rect, error) {
	var Xin []xinerama.ScreenInfo
	if xin, err := xinerama.QueryScreens(ctx.X.Conn()).Reply(); err != nil {
		return false, nil, err
	} else {
		Xin = xin.ScreenInfo
	}

	screens := make([]Rect, 0, len(Xin))
	for _, xi := range Xin {
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
	for c, _ := range ctx.Containers {
		topShape := TopShape(ctx, c.Shape)
		if screen, overlap, _ := ctx.GetScreenForShape(topShape); topShape.Area() > overlap {
			c.MoveResizeShape(ctx, ctx.DefaultShapeForScreen(screen))
		}
	}

	if ctx.Taskbar != nil {
		ctx.Taskbar.MoveResize(ctx)
	}
	ctx.RaiseLock()
}

func (c *Context) Get(w xproto.Window) *Frame {
	f, _ := c.Tracked[w]
	return f
}

func (ctx *Context) MinShape() Rect {
	return Rect{
		X: 0,
		Y: 0,
		W: 5 * ctx.Config.ElemSize,
		H: 5 * ctx.Config.ElemSize,
	}
}

func (ctx *Context) GetScreenForShape(shape Rect) (Rect, int, int) {
	max_overlap := 0
	max_i := 0
	screen := ctx.Screens[0]
	for i, s := range ctx.Screens {
		overlap := AreaOfIntersection(shape, s)
		if overlap > max_overlap {
			max_overlap = overlap
			max_i = i
			screen = s
		}
	}
	return screen, max_overlap, max_i
}

func (ctx *Context) DefaultShapeForScreen(screen Rect) Rect {
	osize := ctx.Config.ElemSize * 2
	offset := 0
	for {
		s := Rect{
			X: screen.X + int(ctx.Config.DefaultShapeRatio.X*float64(screen.W)) + offset*osize,
			Y: screen.Y + int(ctx.Config.DefaultShapeRatio.Y*float64(screen.H)) + offset*osize,
			W: int(ctx.Config.DefaultShapeRatio.W * float64(screen.W)),
			H: int(ctx.Config.DefaultShapeRatio.H * float64(screen.H)),
		}
		tshape := TopShape(ctx, s)
		tscreen, overlap, _ := ctx.GetScreenForShape(tshape)
		if tscreen != screen || overlap < tshape.Area() {
			break
		}

		clear := true
		for c, _ := range ctx.Containers {
			if c.Shape == s {
				clear = false
				break
			}
		}
		if clear {
			return s
		}
		offset++
	}
	return Rect{
		X: screen.X + int(ctx.Config.DefaultShapeRatio.X*float64(screen.W)),
		Y: screen.Y + int(ctx.Config.DefaultShapeRatio.Y*float64(screen.H)),
		W: int(ctx.Config.DefaultShapeRatio.W * float64(screen.W)),
		H: int(ctx.Config.DefaultShapeRatio.H * float64(screen.H)),
	}
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
