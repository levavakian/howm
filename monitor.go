package main

import (
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgbutil"
	"howm/frame"
	"howm/sideloop"
	"log"
	"os/exec"
)

func MonitorScreens(ctx *frame.Context, inj *sideloop.Injector) {
	go func() {
		XR, _ := xgbutil.NewConn()
		defer XR.Conn().Close()

		err := randr.Init(XR.Conn())
		if err != nil {
			log.Fatal(err)
		}

		err = randr.SelectInputChecked(XR.Conn(), XR.RootWin(),
			randr.NotifyMaskScreenChange).Check()
		if err != nil {
			log.Fatal(err)
		}

		for {
			_, err := XR.Conn().WaitForEvent()
			if err != nil {
				log.Println(err)
			}
			exec.Command("xrandr", "--auto").Run()
			inj.Do(func() {
				ctx.UpdateScreens()
			})
		}
	}()
}
