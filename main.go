// Example show-window-icons shows how to get a list of all top-level client
// windows managed by the currently running window manager, and show the icon
// for each window. (Each icon is shown by opening its own window.)
package main

import (
  "log"
  "os/exec"
  // "time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/xgb/xproto"
  // "github.com/BurntSushi/xgbutil/xgraphics"
  "github.com/BurntSushi/xgb/xinerama"
  "howm/frame"
  // "github.com/BurntSushi/wingo/prompt"
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

  // err = keybind.KeyPressFun(
	// 	func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
  //     focus, err := xproto.GetInputFocus(X.Conn()).Reply()
  //     if err != nil {
  //       log.Println(err)
  //       return
  //     }

  //     attachFrame, ok := frame.TrackedFrames[focus.Focus]
  //     if !ok {
  //       parent, err := xwindow.New(X, focus.Focus).Parent()
  //       if err == nil {
  //         attachFrame, ok = frame.TrackedFrames[parent.Id]
  //       }
  //       if !ok || err != nil {
  //         msgPrompt := prompt.NewMessage(X,
  //           prompt.DefaultMessageTheme, prompt.DefaultMessageConfig)
  //         timeout := 2 * time.Second
  //         msgPrompt.Show(xwindow.RootGeometry(X), "Cannot split when not focused on a window", timeout, func(msg *prompt.Message){})
  //         return
  //       }
  //     }

  //     inpPrompt := prompt.NewInput(X,
  //       prompt.DefaultInputTheme, prompt.DefaultInputConfig)

  //     canc := func (inp *prompt.Input) {
  //       log.Println("canceled")
  //     }

  //     resp := func (inp *prompt.Input, text string) {
  //       cmd := exec.Command(text)
  //       err = cmd.Start()
  //       if err != nil {
  //         log.Println(err)
  //       }
  //       go func() {
  //         cmd.Wait()
  //       }()
  //       inpPrompt.Destroy()
  //     }
      
  //     inpPrompt.Show(xwindow.RootGeometry(X),
  //       "Command:", resp, canc)
      
  //     frame.NextFrameAttachHorizontal = attachFrame
  //   }).Connect(X, X.RootWin(), "Mod4-e", true)
  // if err != nil {
  //   log.Println(err)
  // }

  // xevent.MapRequestFun(
	// 	func(X *xgbutil.XUtil, ev xevent.MapRequestEvent) {
  //     if frame.NextFrameAttachHorizontal != nil {
  //       frame.NewHorizontal(X, frame.NextFrameAttachHorizontal, ev)
  //       frame.NextFrameAttachHorizontal = nil
  //       frame.NextFrameAttachVertical = nil
  //     } else if frame.NextFrameAttachVertical != nil {
  //       frame.NewVertical(X, frame.NextFrameAttachVertical, ev)
  //       frame.NextFrameAttachHorizontal = nil
  //       frame.NextFrameAttachVertical = nil
  //     } else {
  //       frame.New(X, ev)
  //     }
  //   }).Connect(X, X.RootWin())

  xevent.MapRequestFun(
		func(X *xgbutil.XUtil, ev xevent.MapRequestEvent) {
      log.Println(ev)
      frame.NewContainer(c, ev)
    }).Connect(X, X.RootWin())

  return err
}