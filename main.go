// Example show-window-icons shows how to get a list of all top-level client
// windows managed by the currently running window manager, and show the icon
// for each window. (Each icon is shown by opening its own window.)
package main

import (
  "log"
  "errors"
  "os/exec"
  "time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/xgb/xproto"
  // "github.com/BurntSushi/xgbutil/xgraphics"
  "github.com/BurntSushi/xgb/xinerama"
  "howm/frame"
  "howm/ext"
  "github.com/BurntSushi/wingo/prompt"
)

var (
	// The icon width and height to use.
	// _NET_WM_ICON will be searched for an icon closest to these values.
	// The icon closest in size to what's specified here will be used.
	// The resulting icon will be scaled to this size.
	// (Set both to 0 to avoid scaling and use the biggest possible icon.)
	iconWidth, iconHeight = 300, 300
)

func main() {
  log.SetFlags(log.LstdFlags | log.Lshortfile)

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
  }
  defer X.Conn().Close()

  // Init Xinerama
  var Xin []xinerama.ScreenInfo
  if err := xinerama.Init(X.Conn()); err != nil {
    log.Fatal(err)
  }
  if xin, err := xinerama.QueryScreens(X.Conn()).Reply(); err != nil {
    log.Fatal(err)
  } else {
    Xin = xin.ScreenInfo
  }
  log.Println("found", len(Xin), "screen(s)", Xin)

  // Take ownership
  if err := ConfigRoot(X); err != nil {
    log.Fatal(err)
  }

	// All we really need to do is block, so a 'select{}' would be sufficient.
	// But running the event loop will emit errors if anything went wrong.
  xevent.Main(X)
  log.Println("Exited")
}

func ConfigRoot(X *xgbutil.XUtil) error {
  var err error

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

  c, err := frame.NewContext(X)
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

    attachF, err := func()(*frame.Frame, error){
      focus, err := xproto.GetInputFocus(X.Conn()).Reply()
      if err != nil {
        log.Println(err)
        return nil, err
      }

      found, ok := c.Tracked[focus.Focus]
      if ok {
        return found, nil
      }

      parent, err := xwindow.New(X, focus.Focus).Parent()
      if err == nil {
        found, ok = c.Tracked[parent.Id]
        if ok {
          return found, nil
        }
      }

      return nil, errors.New("not found")
    }()

    if err != nil {
      log.Println(err)
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
    c.AttachPoint = &frame.AttachTarget{
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
    c.AttachPoint = &frame.AttachTarget{
      Target: fr,
      Type: frame.VERTICAL,
    }
  }).Connect(X, X.RootWin(), "Mod4-r", true)
  if err != nil {
    log.Println(err)
  }

  err = mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
      ext.Focus(xwindow.New(X, X.RootWin()))
			xproto.AllowEvents(X.Conn(), xproto.AllowReplayPointer, 0)
		}).Connect(X, X.RootWin(), c.Config.ButtonClick, true, true)

  xevent.MapRequestFun(
		func(X *xgbutil.XUtil, ev xevent.MapRequestEvent) {
      frame.NewContainer(c, ev)
    }).Connect(X, X.RootWin())

  return err
}