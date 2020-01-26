package ext

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/xgbutil"
)

func IMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MapChecked(w *xwindow.Window) error {
	if w == nil {
		return nil
	}
	return xproto.MapWindowChecked(w.X.Conn(), w.Id).Check()
}

func Focus(w *xwindow.Window) {
	mode := byte(xproto.InputFocusNone)
	err := xproto.SetInputFocusChecked(w.X.Conn(), mode, w.Id, 0).Check()
	if err != nil {
		xgbutil.Logger.Println(err)
	}
}