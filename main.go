package main

import (
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/levavakian/rowm/frame"
	"github.com/levavakian/rowm/root"
	"github.com/levavakian/rowm/sideloop"
	"log"
	"os"
	"io"
	"os/exec"
	"time"
)

func main() {
	logFile, err := os.OpenFile("/var/log/rowm.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
	    panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("RoWM Window Manager")
	log.Println("Hybrid Floating and Tiling Window Manager")
	log.Println("Carbonated for your pleasure")

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	defer X.Conn().Close()

	inj := sideloop.NewInjector()

	// Configure root hooks
	if err := ConfigRoot(X, inj); err != nil {
		log.Fatal(err)
	}

	pingBefore, pingAfter, pingQuit := xevent.MainPing(X)
	for {
		select {
		case <-pingBefore:
			// Wait for the event to finish processing.
			<-pingAfter
		case <-inj.WorkRequest:
			<-inj.WorkNotify
		case <-pingQuit:
			log.Println("xevent loop has quit")
			return
		}
	}
}

func ConfigRoot(X *xgbutil.XUtil, inj *sideloop.Injector) error {
	var err error

	// Init
	err = randr.Init(X.Conn())
	if err != nil {
		log.Fatal(err)
	}
	keybind.Initialize(X)
	mousebind.Initialize(X)

	// Create context
	exec.Command("xrandr", "--auto").Run()
	ctx, err := frame.NewContext(X, inj)
	if err != nil {
		log.Fatal(err)
	}

	// Set x cursor
	err = exec.Command(ctx.Config.Shell, "-c", "xsetroot -cursor_name arrow").Run()
	if err != nil {
		log.Fatal(err)
	}

	// Listen for events we're interested in
	evMasks := xproto.EventMaskPropertyChange |
		xproto.EventMaskFocusChange |
		xproto.EventMaskButtonPress |
		xproto.EventMaskButtonRelease |
		xproto.EventMaskStructureNotify |
		xproto.EventMaskSubstructureNotify |
		xproto.EventMaskSubstructureRedirect

	err = xwindow.New(X, X.RootWin()).Listen(evMasks)
	if err != nil {
		log.Fatal(err)
	}

	// Start monitor for screens
	root.MonitorScreens(ctx, inj)

	// Add base control hooks
	err = root.RegisterBaseHooks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Add splitting hooks
	err = root.RegisterSplitHooks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Add volume hooks
	err = root.RegisterVolumeHooks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Add backlight hooks
	err = root.RegisterBrightnessHooks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Add alttab-like hooks
	root.RegisterChooseHooks(ctx)

	// Add taskbar hooks
	err = root.RegisterTaskbarHooks(ctx)
	if err != nil {
		log.Fatal(err)
	}
	sideloop.NewRepeater(func() { ctx.Taskbar.Update(ctx) }, 1*time.Second, inj)

	return err
}
