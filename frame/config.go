package frame

import (
	"github.com/BurntSushi/xgbutil/xcursor"
)

type Config struct {
	ButtonClose string
	ButtonDrag string
	ButtonClick string
	CloseFrame string
	ToggleExpandFrame string
	ToggleExternalDecorator string
	ElemSize int
	CloseCursor int
	DefaultShapeRatio Rectf
	SeparatorColor uint32
	GrabColor uint32
	CloseColor uint32
	ResizeColor uint32
	InternalPadding int
	BackgroundImagePath string
}

func DefaultConfig() Config {
	return Config{
		ButtonClose: "1",
		ButtonDrag: "1",
		ButtonClick: "1",
		CloseFrame: "Mod4-d",
		ToggleExpandFrame: "Mod4-f",
		ToggleExternalDecorator: "Mod4-h",
		ElemSize: 10,
		CloseCursor: xcursor.Dot,
		DefaultShapeRatio: Rectf {
			X: .05,
			Y: .05,
			W: .9,
			H: .9,
		},
		SeparatorColor: 0x777777,
		GrabColor: 0x339999,
		CloseColor: 0xff0000,
		ResizeColor: 0x00ff00,
		InternalPadding: 0,
		BackgroundImagePath: "/home/lev/.config/howm/bg.jpg",
	}
}