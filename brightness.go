package main

import (
	"github.com/BurntSushi/wingo/prompt"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"howm/ext"
	"howm/frame"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetBrightnessAndMaxBrightness(ctx *frame.Context) (int, int, error) {
	read := func(filename string) (int, error) {
		in, err := ioutil.ReadFile(filename)
		if err != nil {
			return 0, err
		}
		return strconv.Atoi(strings.TrimSpace(string(in)))
	}

	bright, err := read(ctx.Config.BrightFile())
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}

	maxBright, err := read(ctx.Config.MaxBrightFile())
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}

	return bright, maxBright, nil
}

func RegisterBrightnessHooks(ctx *frame.Context) error {
	brightnessMod := func(increment float64) {
		if ctx.Locked {
			return
		}

		current, max, err := GetBrightnessAndMaxBrightness(ctx)
		if err != nil || max <= 0 {
			return
		}

		toPercent := func(top, bottom int) float64 {
			return float64(top) / float64(bottom)
		}

		newBright := ext.IClamp(int(ext.Clamp(toPercent(current, max)+increment, 0.0, 1.0)*float64(max)), 1, max)

		err = exec.Command("howmbright", ctx.Config.Backlight, strconv.Itoa(newBright)).Run()

		current, max, err = GetBrightnessAndMaxBrightness(ctx)
		if err != nil {
			log.Println(err)
			return
		}

		msgPrompt := prompt.NewMessage(ctx.X,
			prompt.DefaultMessageTheme, prompt.DefaultMessageConfig)
		timeout := 1 * time.Second
		msgPrompt.Show(ctx.Screens[0].ToXRect(), strconv.Itoa(int(toPercent(current, max)*100)), timeout, func(msg *prompt.Message) {})
	}

	var err error
	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		brightnessMod(.02)
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.BrightnessUp, true)
	if err != nil {
		log.Println(err)
		return err
	}

	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		brightnessMod(-.02)
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.BrightnessDown, true)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
