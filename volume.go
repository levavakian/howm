package main

import (
	"log"
	"time"
	"strconv"
	"os/exec"
	"regexp"
	"howm/frame"
	"howm/ext"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/wingo/prompt"
)

func GetCurrentAudio() (int, error) {
	out, err := exec.Command("amixer", "sget", "Master").Output()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	r := regexp.MustCompile("\\[(\\d+)%\\]")
	intstr := r.FindStringSubmatch(string(out))
	if len(intstr) < 1 {
		return strconv.Atoi("")
	}
	return strconv.Atoi(intstr[1])
}

type VolumeContext struct {
	Volume int
}

func RegisterVolumeHooks(ctx *frame.Context) error {
	VolumeBeforeMute := &VolumeContext{}
	audioMod := func(increment int) {
		target := 0
		current, err := GetCurrentAudio()
		if err != nil {
			log.Println(err)
			return
		}
		if increment != 0 {
			target = ext.IClamp(current + increment, 0, 100)
		}

		if increment == 0 {
			if VolumeBeforeMute.Volume == 0 {
				err = exec.Command("amixer", "sset", "Master", "0%").Run()
				VolumeBeforeMute.Volume = current
			} else {
				err = exec.Command("amixer", "sset", "Master", strconv.Itoa(VolumeBeforeMute.Volume) + "%").Run()
				VolumeBeforeMute.Volume = 0
			}
		} else {
			err = exec.Command("amixer", "sset", "Master", strconv.Itoa(target) + "%").Run()
		}
		if err != nil {
			log.Println(err)
		}
		current, err = GetCurrentAudio()
		if err != nil {
			log.Println(err)
			return
		}
		
		msgPrompt := prompt.NewMessage(ctx.X,
			prompt.DefaultMessageTheme, prompt.DefaultMessageConfig)
		timeout := 1 * time.Second
		msgPrompt.Show(ctx.Screens[0].ToXRect(), strconv.Itoa(current), timeout, func(msg *prompt.Message){})
	}

	var err error
	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent){
		audioMod(2)	
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.VolumeUp, true)
	if err != nil {
		log.Println(err)
		return err
	}

	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent){
		audioMod(-2)	
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.VolumeDown, true)
	if err != nil {
		log.Println(err)
		return err
	}

	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent){
		audioMod(0)	
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.VolumeMute, true)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}