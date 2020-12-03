package frame

import (
	"fmt"
	"github.com/BurntSushi/xgbutil/xcursor"
	"github.com/spf13/viper"
	"log"
	"os/user"
	"path"
	"time"
)

type StringWithHelp struct {
	Data	string	`json:"data,omitempty"`
	Help	string	`json:"help,omitempty"`
}


type Config struct {
	Shell                     string `json:"shell,omitempty"`
	Lock                      string `json:"lock,omitempty"`
	TabByFrame                bool `json:"tab_frame,omitempty"`
	TabForward                StringWithHelp `json:"tab_forward,omitempty"`
	TabBackward               StringWithHelp `json:"tab_backward,omitempty"`
	ButtonDrag                string `json:"button_drag,omitempty"`
	ButtonClick               string `json:"button_click,omitempty"`
	SplitVertical             StringWithHelp `json:"split_vertical,omitempty"`
	SplitHorizontal           StringWithHelp `json:"split_horizontal,omitempty"`
	RunCmd                    StringWithHelp `json:"run_cmd,omitempty"`
	Shutdown                  string `json:"shutdown,omitempty"`
	CloseFrame                StringWithHelp `json:"close_frame,omitempty"`
	ToggleExpandFrame         StringWithHelp `json:"toggle_expand_frame,omitempty"`
	ToggleExternalDecorator   string `json:"toggle_external_decorator,omitempty"`
	ToggleTaskbar             string `json:"toggle_taskbar,omitempty"`
	PopFrame                  StringWithHelp `json:"pop_frame,omitempty"`
	ResetSize                 string `json:"reset_size,omitempty"`
	Minimize                  string `json:"minimize,omitempty"`
	WindowUp		  StringWithHelp `json:"window_up,omitempty"`
	WindowDown                StringWithHelp `json:"window_down,omitempty"`
	WindowLeft                StringWithHelp `json:"window_left,omitempty"`
	WindowRight               StringWithHelp `json:"window_right,omitempty"`
	VolumeUp                  string `json:"volume_up,omitempty"`
	VolumeDown                string `json:"volume_down,omitempty"`
	BrightnessUp              string `json:"brightness_up,omitempty"`
	BrightnessDown            string `json:"brightness_down,omitempty"`
	Backlight                 string `json:"backlight,omitempty"`
	VolumeMute                string `json:"volume_mute,omitempty"`
	FocusNext                 StringWithHelp `json:"focus_next,omitempty"`
	FocusPrev                 StringWithHelp `json:"focus_prev,omitempty"`
	ElemSize                  int `json:"elem_size,omitempty"`
	CloseCursor               int `json:"close_cursor,omitempty"`
	DefaultShapeRatio         Rectf
	SeparatorColor            uint32 `json:"seperator_color,omitempty"`
	GrabColor                 uint32 `json:"grab_color,omitempty"`
	FocusColor                uint32 `json:"focus_color,omitempty"`
	CloseColor                uint32 `json:"close_clor,omitempty"`
	MaximizeColor             uint32 `json:"maximize_color,omitempty"`
	MinimizeColor             uint32 `json:"minimize_color,omitempty"`
	ResizeColor               uint32 `json:"resize_color,omitempty"`
	TaskbarBaseColor          uint32 `json:"taskbar_base_color,omitempty"`
	TaskbarTextColor          uint32 `json:"taskbar_text_color,omitempty"`
	InternalPadding           int `json:"internal_padding,omitempty"`
	BackgroundImagePath       string `json:"background_image_path,omitempty"`
	BuiltinCommands           map[StringWithHelp]string `json:"builtin_commands,omitempty"`
	FocusMarkerTime           time.Duration
	DoubleClickTime           time.Duration
	TaskbarHeight             int `json:"taskbar_height,omitempty"`
	TaskbarSlideWidth         int `json:"taskbar_slide_width,omitempty"`
	TaskbarSlideActiveColor   uint32 `json:"taskbar_slide_active_color,omitempty"`
	TaskbarSlideInactiveColor uint32 `json:"taskbar_slide_inactive_color,omitempty"`
	TaskbarFontSize           float64 `json:"taskbar_font_size,omitempty"`
	TaskbarTimeBaseColor      uint32 `json:"taskbar_time_base_color,omitempty"`
	TaskbarXPad               int `json:"taskbar_x_pad,omitempty"`
	TaskbarYPad               int `json:"taskbar_y_pad,omitempty"`
	TaskbarTimeFormat         string `json:"taskbar_time_format,omitempty"`
	TaskbarBatFormat          string `json:"taskbar_bat_format,omitempty"`
	TaskbarElementShape       Rect
	TaskbarMinMaxHeight       int `json:"taskbar_min_max_height,omitempty"`
	TaskbarMinMaxColor        uint32 `json:"taskbar_min_max_color,omitempty"`
	TaskbarSlideLeft          string `json:"taskbar_slide_left,omitempty"`
	TaskbarSlideRight         string `json:"taskbar_slide_right,omitempty"`
	CutSelectFrame            string `json:"cut_select_frame,omitempty"`
	CutSelectContainer        string `json:"cut_select_container,omitempty"`
	CopySelectHorizontal      StringWithHelp `json:"copy_select_horizontal,omitempty"`
	CopySelectVertical        StringWithHelp `json:"copy_select_vertical,omitempty"`
	SuspendCommand            string `json:"suspend_command,omitempty"`
	BatteryWarningLevels      []int
	BatteryWarningDuration    time.Duration
	LaunchHelp                string `json:"launch_help,omitempty"`
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

func LoadConfig() Config {
	conf := DefaultConfig()
	viper.SetConfigName("rowm") // name of config file (without extension)
	viper.SetConfigType("json") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/rowm/")   // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		log.Println("Didn't load config file:", err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Println("Unable to decode into config struct:", err)
	}
	return conf
}


// DefaultConfig reference for key strings
// https://github.com/BurntSushi/xgbutil/blob/master/keybind/keysymdef.go
func DefaultConfig() Config {
	return Config{
		Shell:                   "/bin/bash",
		Lock:                    "Mod4-l",
		TabByFrame:              true,
		TabForward:              StringWithHelp{Data: "Mod1-tab", Help:"Tab Forward"},
		TabBackward:             StringWithHelp{Data: "Mod1-Shift-tab", Help:"Tab Backward"},
		ButtonDrag:              "1",
		ButtonClick:             "1",
		SplitVertical:           StringWithHelp{Data: "Mod4-r", Help:"Split Vertically"},
		SplitHorizontal:         StringWithHelp{Data: "Mod4-e", Help:"Split Horizontally"},
		RunCmd:                  StringWithHelp{Data: "Mod4-f", Help:"Run Command"},
		Shutdown:                "Mod4-BackSpace",
		CloseFrame:              StringWithHelp{Data: "Mod4-d", Help:"Close Frame"},
		ToggleExpandFrame:       StringWithHelp{Data: "Mod4-x", Help:"Toggle Expanded Frame"},
		ToggleExternalDecorator: "Mod4-h",
		ToggleTaskbar:           "Mod4-s",
		WindowUp:                 StringWithHelp{Data: "Mod4-up", Help:"Window up"},
		WindowDown:               StringWithHelp{Data: "Mod4-down", Help:"Window Down"},
		WindowLeft:               StringWithHelp{Data: "Mod4-left", Help:"Window Left"},
		WindowRight:              StringWithHelp{Data: "Mod4-right", Help:"Window Right"},
		PopFrame:                 StringWithHelp{Data: "Mod4-q", Help:"Pop Frame"},
		ResetSize:               "Mod4-Shift-up",
		Minimize:                "Mod4-Shift-down",
		VolumeUp:                "XF86AudioRaiseVolume",
		VolumeDown:              "XF86AudioLowerVolume",
		VolumeMute:              "XF86AudioMute",
		BrightnessUp:            "XF86MonBrightnessUp",
		BrightnessDown:          "XF86MonBrightnessDown",
		FocusNext:               StringWithHelp{Data: "Mod4-Tab", Help:"Focus Next"},
		FocusPrev:               StringWithHelp{Data: "Mod4-asciitilde", Help:"Focus Previous"},
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
		BuiltinCommands: map[StringWithHelp]string{
			StringWithHelp{Data: "Mod4-t", Help:"Terminal"}:  "x-terminal-emulator",
			StringWithHelp{Data: "Mod4-w", Help:"Chrome"}:  "google-chrome",
			StringWithHelp{Data: "Mod4-p", Help: "Gnome"}:  "XDG_CURRENT_DESKTOP=GNOME gnome-control-center",
			StringWithHelp{Data: "Mod4-o", Help:" XDG"}:"xdg-open .",
			StringWithHelp{Data: "XF86AudioPlay", Help:"Pause"}: "dbus-send --print-reply --dest=org.mpris.MediaPlayer2.spotify /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.PlayPause",
			StringWithHelp{Data: "XF86AudioPrev", Help: "Previous"}: "dbus-send --print-reply --dest=org.mpris.MediaPlayer2.spotify /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.Previous",
			StringWithHelp{Data: "XF86AudioNext", Help: "Next"}: "dbus-send --print-reply --dest=org.mpris.MediaPlayer2.spotify /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.Next",
			StringWithHelp{Data: "Print", Help: "Screenshot"}:"gnome-screenshot -i",
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
		CopySelectHorizontal:      StringWithHelp{Data: "Mod4-v", Help: "Paste Horizontally"},
		CopySelectVertical:        StringWithHelp{Data: "Mod4-b", Help: "Paste Vertically"},
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
