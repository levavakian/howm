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
	FocusNext                 string
	FocusPrev                 string
	ElemSize                  int
	CloseCursor               int
	DefaultShapeRatio         Rectf
	SeparatorColor            uint32
	GrabColor                 uint32
	FocusColor                uint32
	CloseColor                uint32
	MaximizeColor             uint32
	MinimizeColor             uint32
	ResizeColor               uint32
	TaskbarBaseColor          uint32
	TaskbarTextColor          uint32
	InternalPadding           int
	BackgroundImagePath       string
	BuiltinCommands           map[string]string
	FocusMarkerTime           time.Duration
	DoubleClickTime           time.Duration
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
	SuspendCommand            string
	BatteryWarningLevels      []int
	BatteryWarningDuration    time.Duration
	LaunchHelp                string
	GotoKeys                  map[string]string
}

func HomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
		return ""
	}
	return usr.HomeDir
}

func (c *Config) MinShape() Rect {
	return Rect{
		X: 0,
		Y: 0,
		W: 5 * c.ElemSize,
		H: 5 * c.ElemSize,
	}
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
		FocusNext:               "Mod4-Tab",
		FocusPrev:               "Mod4-asciitilde",
		Backlight:               "intel_backlight",
		ElemSize:                10,
		CloseCursor:             xcursor.Dot,
		DefaultShapeRatio: Rectf{
			X: .1,
			Y: .1,
			W: .8,
			H: .8,
		},
		SeparatorColor:      0x777777,
		GrabColor:           0x339999,
		FocusColor:          0x9932cc,
		ResizeColor:         0x777777,
		TaskbarBaseColor:    0x222222,
		TaskbarTextColor:    0xbbbbbb,
		CloseColor:          0xff0000,
		MaximizeColor:       0x00ff00,
		MinimizeColor:       0xfdfd96,
		InternalPadding:     0,
		BackgroundImagePath: path.Join(HomeDir(), ".config/rowm/bg.png"),
		FocusMarkerTime:     time.Millisecond * 350,
		DoubleClickTime:     time.Millisecond * 500,
		BuiltinCommands: map[string]string{
			"Mod4-t":  "x-terminal-emulator",
			"Mod4-w":  "google-chrome",
			"Mod4-p":  "XDG_CURRENT_DESKTOP=GNOME gnome-control-center",
			"Mod4-o":  "xdg-open .",
			"Mod4-F5": "dbus-send --print-reply --dest=org.mpris.MediaPlayer2.spotify /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.PlayPause",
			"Mod4-F4": "dbus-send --print-reply --dest=org.mpris.MediaPlayer2.spotify /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.Previous",
			"Mod4-F6": "dbus-send --print-reply --dest=org.mpris.MediaPlayer2.spotify /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.Next",
			"Print":   "gnome-screenshot -i",
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
		BatteryWarningLevels:      []int{20, 10, 5, 1},
		BatteryWarningDuration:    time.Second * 2,
		LaunchHelp:                "Mod4-Shift-h",
		GotoKeys: map[string]string{
			"Mod4-Shift-0": "Mod4-0",
			"Mod4-Shift-1": "Mod4-1",
			"Mod4-Shift-2": "Mod4-2",
			"Mod4-Shift-3": "Mod4-3",
			"Mod4-Shift-4": "Mod4-4",
			"Mod4-Shift-5": "Mod4-5",
			"Mod4-Shift-6": "Mod4-6",
			"Mod4-Shift-7": "Mod4-7",
			"Mod4-Shift-8": "Mod4-8",
			"Mod4-Shift-9": "Mod4-9",
		},
	}
}

func (c *Config) BrightFile() string {
	return fmt.Sprintf("/sys/class/backlight/%s/brightness", c.Backlight)
}

func (c *Config) MaxBrightFile() string {
	return fmt.Sprintf("/sys/class/backlight/%s/max_brightness", c.Backlight)
}
