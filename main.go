package main

import (
  "log"
  "os/exec"
  "howm/frame"
  "howm/ext"
  "howm/sideloop"
  "github.com/BurntSushi/xgbutil"
  "github.com/BurntSushi/xgbutil/xevent"
  "github.com/BurntSushi/xgbutil/keybind"
  "github.com/BurntSushi/xgbutil/mousebind"
  "github.com/BurntSushi/xgbutil/xwindow"
  "github.com/BurntSushi/xgb/xproto"
  "github.com/BurntSushi/xgb/randr"
  "github.com/BurntSushi/wingo/prompt"
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

  // Make sure we can leave
  err = keybind.KeyReleaseFun(
		func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
      xevent.Quit(X)
    }).Connect(X, X.RootWin(), ctx.Config.Shutdown, true)
  if err != nil {
    log.Fatal(err)
  }

  for k, v := range(ctx.Config.BuiltinCommands) {
    ncmd := v  // force to not be a reference
    err = keybind.KeyReleaseFun(
      func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
        if ctx.SplitPrompt != nil {
          ctx.SplitPrompt.Destroy()
        }
        cmd := exec.Command(ncmd)
        err := cmd.Start()
        if err != nil {
          log.Println(err)
        }
        go func() {
          cmd.Wait()
        }()
      }).Connect(X, X.RootWin(), k, true)
    if err != nil {
      log.Println(err)
    }
  }

  err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent){
    inPrompt := prompt.NewInput(X,
      prompt.DefaultInputTheme, prompt.DefaultInputConfig)

    canc := func (inp *prompt.Input) {
      // Chill
    }

    resp := func (inp *prompt.Input, text string) {
      cmd := exec.Command(text)
      err = cmd.Start()
      if err != nil {
        log.Println(err)
      }
      go func() {
        cmd.Wait()
      }()
      inPrompt.Destroy()
    }
    
    inPrompt.Show(ctx.Screens[0].ToXRect(),
      "Command:", resp, canc)
  }).Connect(X, X.RootWin(), ctx.Config.RunCmd, true)
  if err != nil {
    log.Println(err)
  }

  err = mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
      ext.Focus(xwindow.New(X, X.RootWin()))
			xproto.AllowEvents(X.Conn(), xproto.AllowReplayPointer, 0)
    }).Connect(X, X.RootWin(), ctx.Config.ButtonClick, true, true)
    if err != nil {
      log.Println(err)
    }

  xevent.MapRequestFun(
		func(X *xgbutil.XUtil, ev xevent.MapRequestEvent) {
      frame.NewWindow(ctx, ev)
    }).Connect(X, X.RootWin())

  return err
}