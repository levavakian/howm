package main

import (
  "log"
  "os/exec"
  "time"
  "howm/frame"
  "howm/ext"
  "howm/background"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/xgb/xproto"
  "github.com/BurntSushi/wingo/prompt"
)

func main() {
  log.SetFlags(log.LstdFlags | log.Lshortfile)

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
  }
  defer X.Conn().Close()

  // Take ownership
  if err := ConfigRoot(X); err != nil {
    log.Fatal(err)
  }

  xevent.Main(X)
  log.Println("Exited")
}

func ConfigRoot(X *xgbutil.XUtil) error {
  var err error

  // Set x cursor
  err = exec.Command("xsetroot", "-cursor_name", "arrow").Run()
  if err != nil {
    log.Println(err)
  }

  // Create context
  ctx, err := frame.NewContext(X)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("found", len(ctx.ScreenInfos), "screen(s)", ctx.ScreenInfos)

  // Background Images
  background.GenerateBackgrounds(ctx)

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

  keybind.Initialize(X)
	mousebind.Initialize(X)

  err = keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
      xevent.Quit(X)
    }).Connect(X, X.RootWin(), "Mod4-BackSpace", true)
  if err != nil {
    log.Println(err)
  }

  err = keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
      cmd := exec.Command("x-terminal-emulator")
      err := cmd.Start()
      if err != nil {
        log.Println(err)
      }
      go func() {
        cmd.Wait()
      }()
    }).Connect(X, X.RootWin(), "Mod4-t", true)
  if err != nil {
    log.Println(err)
  }

  splitF := func() *frame.Frame {
    if err != nil {
      log.Println(err)
      return nil
    }

    attachF := ctx.GetFocusedFrame()

    if attachF == nil {
      msgPrompt := prompt.NewMessage(X,
        prompt.DefaultMessageTheme, prompt.DefaultMessageConfig)
      timeout := 2 * time.Second
      msgPrompt.Show(xwindow.RootGeometry(X), "Cannot split when not focused on a window", timeout, func(msg *prompt.Message){})
      return nil
    }

    inpPrompt := prompt.NewInput(X,
      prompt.DefaultInputTheme, prompt.DefaultInputConfig)

    canc := func (inp *prompt.Input) {
      log.Println("canceled")
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
      inpPrompt.Destroy()
    }
    
    inpPrompt.Show(xwindow.RootGeometry(X),
      "Command:", resp, canc)

    return attachF
  }

  err = keybind.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent){
    fr := splitF()
    if fr == nil {
      return
    }
    ctx.AttachPoint = &frame.AttachTarget{
      Target: fr,
      Type: frame.HORIZONTAL,
    }
  }).Connect(X, X.RootWin(), "Mod4-e", true)
  if err != nil {
    log.Println(err)
  }

  err = keybind.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent){
    fr := splitF()
    if fr == nil {
      return
    }
    ctx.AttachPoint = &frame.AttachTarget{
      Target: fr,
      Type: frame.VERTICAL,
    }
  }).Connect(X, X.RootWin(), "Mod4-r", true)
  if err != nil {
    log.Println(err)
  }

  err = keybind.KeyPressFun(func(X *xgbutil.XUtil, e xevent.KeyPressEvent){
    inpPrompt := prompt.NewInput(X,
      prompt.DefaultInputTheme, prompt.DefaultInputConfig)

    canc := func (inp *prompt.Input) {
      log.Println("canceled")
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
      inpPrompt.Destroy()
    }
    
    inpPrompt.Show(xwindow.RootGeometry(X),
      "Command:", resp, canc)
  }).Connect(X, X.RootWin(), "Mod4-c", true)
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
      log.Println(ev)
      frame.NewWindow(ctx, ev)
    }).Connect(X, X.RootWin())

  return err
}