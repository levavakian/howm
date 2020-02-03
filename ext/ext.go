package ext

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xwindow"
	"log"
	"math"
)

func IMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func IClamp(n, min, max int) int {
	return IMax(IMin(n, max), min)
}

func Clamp(n, min, max float64) float64 {
	return math.Max(math.Min(n, max), min)
}

func MapChecked(w *xwindow.Window) error {
	if w == nil {
		return nil
	}
	return xproto.MapWindowChecked(w.X.Conn(), w.Id).Check()
}

func UnmapChecked(w *xwindow.Window) error {
	if w == nil {
		return nil
	}
	return xproto.UnmapWindowChecked(w.X.Conn(), w.Id).Check()
}

func Focus(w *xwindow.Window) {
	mode := byte(xproto.InputFocusNone)
	err := xproto.SetInputFocusChecked(w.X.Conn(), mode, w.Id, 0).Check()
	if err != nil {
		xgbutil.Logger.Println(err)
	}
}

func Logerr(err error) {
	if err != nil {
		log.Println(err)
	}
}
