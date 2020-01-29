package frame

import (
	"log"
	"howm/ext"
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
	ScreenInfos []xinerama.ScreenInfo
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
		ScreenInfos: Xin,
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

func (ctx *Context) GetFocusedFrame() *Frame {
	focus, err := xproto.GetInputFocus(ctx.X.Conn()).Reply()
	if err != nil {
		log.Println(err)
		return nil
	}

	found, ok := ctx.Tracked[focus.Focus]
	if ok {
		return found
	}

	parent, err := xwindow.New(ctx.X, focus.Focus).Parent()
	ext.Logerr(err)
	if err == nil {
		found, ok = ctx.Tracked[parent.Id]
		if ok {
			return found
		}
	}

	return nil
}