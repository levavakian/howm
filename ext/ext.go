package ext

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xwindow"
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
