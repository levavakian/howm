package main

import (
  "log"
  "os/exec"
  "howm/frame"
  "howm/sideloop"
  "github.com/BurntSushi/xgbutil"
  "github.com/BurntSushi/xgbutil/xevent"
  "github.com/BurntSushi/xgbutil/keybind"
  "github.com/BurntSushi/xgbutil/mousebind"
  "github.com/BurntSushi/xgbutil/xwindow"
  "github.com/BurntSushi/xgb/xproto"
  "github.com/BurntSushi/xgb/randr"
)

func main() {
  log.SetFlags(log.LstdFlags | log.Lshortfile)
  log.Println("HOme Window Manager")
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
  ctx, err := frame.NewContext(X)
  if err != nil {
    log.Fatal(err)
  }

  // Set x cursor
  err = exec.Command(ctx.Config.Shell, "-c", "xsetroot -cursor_name arrow").Run()
  if err != nil {
    log.Println(err)
  }

  evMasks := xproto.EventMaskPropertyChange |
		xproto.EventMaskFocusChange |
		xproto.EventMaskButtonPress |
		xproto.EventMaskButtonRelease |
		xproto.EventMaskStructureNotify |
		xproto.EventMaskSubstructureNotify |
    xproto.EventMaskSubstructureRedirect
  
  err = xwindow.New(X, X.RootWin()).Listen(evMasks)
  if err != nil {
    log.Println(err)
  }

  // Start monitor for screens
  MonitorScreens(ctx, inj)

  // Add base control hooks
  err = RegisterBaseHooks(ctx)
  if err != nil {
    log.Fatal(err)
  }

  // Add splitting hooks
  err = RegisterSplitHooks(ctx)
  if err != nil {
    log.Fatal(err)
  }

  // Add volume hooks
  err = RegisterVolumeHooks(ctx)
  if err != nil {
    log.Fatal(err)
  }

  // Add backlight hooks
  err = RegisterBrightnessHooks(ctx)
  if err != nil {
    log.Fatal(err)
  }

  // Add alttab-like hooks
  RegisterChooseHooks(ctx)

  return err
}