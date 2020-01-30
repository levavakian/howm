package frame

import (
	"github.com/BurntSushi/xgbutil/xcursor"
)

type Config struct {
	ButtonClose string
	ButtonDrag string
	ButtonClick string
	SplitVertical string
	SplitHorizontal string
	RunCmd string
	Shutdown string
	CloseFrame string
	ToggleExpandFrame string
	ToggleExternalDecorator string
	ResetSize string
	WindowUp string
	WindowDown string
	WindowLeft string
	WindowRight string
	ElemSize int
	CloseCursor int
	DefaultShapeRatio Rectf
	SeparatorColor uint32
	GrabColor uint32
	CloseColor uint32
	ResizeColor uint32
	InternalPadding int
	BackgroundImagePath string
	BuiltinCommands map[string]string
}

func DefaultConfig() Config {
	return Config{
		ButtonClose: "1",
		ButtonDrag: "1",
		ButtonClick: "1",
		SplitVertical: "Mod4-r",
		SplitHorizontal: "Mod4-e",
		RunCmd: "Mod4-c",
		Shutdown: "Mod4-BackSpace",
		CloseFrame: "Mod4-d",
		ToggleExpandFrame: "Mod4-f",
		ToggleExternalDecorator: "Mod4-h",
		WindowUp: "Mod4-up",
		WindowDown: "Mod4-down",
		WindowLeft: "Mod4-left",
		WindowRight: "Mod4-right",
		ResetSize: "Mod4-Shift-up",
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
		BuiltinCommands: map[string]string{
			"Mod4-t": "x-terminal-emulator",
			"Mod4-w": "google-chrome",
		},
	}
}