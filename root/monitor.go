package root

import (
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgbutil"
	"github.com/levavakian/rowm/frame"
	"github.com/levavakian/rowm/sideloop"
	"log"
	"os/exec"
)

// MonitorScreens runs an X event loop on the side just listening to whether the root geometry has changed.
// Once we notice a change, we pause the main event loop to update the screens, and then start monitoring again.
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
