package frame

import (
	"fmt"
	"github.com/BurntSushi/xgbutil/xcursor"
	"log"
	"os/user"
	"path"
	"time"
)

type Config struct {
	Shell                     string
	Lock                      string
	TabByFrame                bool
	TabForward                string
	TabBackward               string
	ButtonDrag                string
	ButtonClick               string
	SplitVertical             string
	SplitHorizontal           string
	RunCmd                    string
	Shutdown                  string
	CloseFrame                string
	ToggleExpandFrame         string
	ToggleExternalDecorator   string
	ToggleTaskbar             string
	PopFrame                  string
	ResetSize                 string
	Minimize                  string
	WindowUp                  string
	WindowDown                string
	WindowLeft                string
	WindowRight               string
	VolumeUp                  string
	VolumeDown                string
	BrightnessUp              string
	BrightnessDown            string
	Backlight                 string
	VolumeMute                string
	ElemSize                  int
	CloseCursor               int
	DefaultShapeRatio         Rectf
	SeparatorColor            uint32
	GrabColor                 uint32
	CloseColor                uint32
	MaximizeColor             uint32
	MinimizeColor             uint32
	ResizeColor               uint32
	TaskbarBaseColor          uint32
	TaskbarTextColor          uint32
	InternalPadding           int
	BackgroundImagePath       string
	BuiltinCommands           map[string]string
	ScreenPoll                time.Duration
	TaskbarHeight             int
	TaskbarSlideWidth         int
	TaskbarSlideActiveColor   uint32
	TaskbarSlideInactiveColor uint32
	TaskbarFontSize           float64
	TaskbarTimeBaseColor      uint32
	TaskbarXPad               int
	TaskbarYPad               int
	TaskbarTimeFormat         string
	TaskbarBatFormat          string
	TaskbarElementShape       Rect
	TaskbarMinMaxHeight       int
	TaskbarMinMaxColor        uint32
	TaskbarSlideLeft          string
	TaskbarSlideRight         string
	CutSelectFrame            string
	CutSelectContainer        string
	CopySelectHorizontal      string
	CopySelectVertical        string
	SuspendCommand string
}

func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
		return ""
	}
	return usr.HomeDir
}

func DefaultConfig() Config {
	return Config{
		Shell:                   "/bin/bash",
		Lock:                    "Mod4-l",
		TabByFrame:              true,
		TabForward:              "Mod1-tab",
		TabBackward:             "Mod1-Shift-tab",
		ButtonDrag:              "1",
		ButtonClick:             "1",
		SplitVertical:           "Mod4-r",
		SplitHorizontal:         "Mod4-e",
		RunCmd:                  "Mod4-f",
		Shutdown:                "Mod4-BackSpace",
		CloseFrame:              "Mod4-d",
		ToggleExpandFrame:       "Mod4-x",
		ToggleExternalDecorator: "Mod4-h",
		ToggleTaskbar:           "Mod4-s",
		WindowUp:                "Mod4-up",
		WindowDown:              "Mod4-down",
		WindowLeft:              "Mod4-left",
		WindowRight:             "Mod4-right",
		PopFrame:                "Mod4-q",
		ResetSize:               "Mod4-Shift-up",
		Minimize:                "Mod4-Shift-down",
		VolumeUp:                "Mod4-F3",
		VolumeDown:              "Mod4-F2",
		VolumeMute:              "Mod4-F1",
		BrightnessUp:            "Mod4-F12",
		BrightnessDown:          "Mod4-F11",
		Backlight:               "intel_backlight",
		ElemSize:                10,
		CloseCursor:             xcursor.Dot,
		DefaultShapeRatio: Rectf{
			X: .05,
			Y: .05,
			W: .9,
			H: .9,
		},
		SeparatorColor:      0x777777,
		GrabColor:           0x339999,
		ResizeColor:         0x777777,
		TaskbarBaseColor:    0x222222,
		TaskbarTextColor:    0xbbbbbb,
		CloseColor:          0xff0000,
		MaximizeColor:       0x00ff00,
		MinimizeColor:       0xfdfd96,
		InternalPadding:     0,
		BackgroundImagePath: path.Join(HomeDir(), ".config/howm/bg.png"),
		ScreenPoll:          time.Second * 2,
		BuiltinCommands: map[string]string{
			"Mod4-t": "x-terminal-emulator",
			"Mod4-w": "google-chrome",
			"Mod4-p": "XDG_CURRENT_DESKTOP=GNOME gnome-control-center",
		},
		TaskbarHeight:        20,
		TaskbarFontSize:      12,
		TaskbarTimeBaseColor: 0x222222,
		TaskbarXPad:          5,
		TaskbarYPad:          5,
		TaskbarTimeFormat:    "2006 Mon Jan 02 - 15:04:05 (MST)",
		TaskbarBatFormat:     "%s%3d%%",
		TaskbarElementShape: Rect{
			X: 2,
			Y: 0,
			W: 16,
			H: 16,
		},
		TaskbarSlideWidth:         10,
		TaskbarSlideActiveColor:   0x666666,
		TaskbarSlideInactiveColor: 0x333333,
		TaskbarMinMaxHeight:       4,
		TaskbarMinMaxColor:        0x999999,
		TaskbarSlideLeft:          "Mod4-Shift-left",
		TaskbarSlideRight:         "Mod4-Shift-right",
		CutSelectFrame:            "Mod4-c",
		CutSelectContainer:        "Mod4-Shift-c",
		CopySelectHorizontal:      "Mod4-v",
		CopySelectVertical:        "Mod4-b",
		SuspendCommand:            "systemctl suspend",
	}
}

func (c *Config) BrightFile() string {
	return fmt.Sprintf("/sys/class/backlight/%s/brightness", c.Backlight)
}

func (c *Config) MaxBrightFile() string {
	return fmt.Sprintf("/sys/class/backlight/%s/max_brightness", c.Backlight)
}
