package main

import (
	"github.com/BurntSushi/wingo/prompt"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"howm/frame"
	"log"
	"os/exec"
	"time"
)

func Split(ctx *frame.Context) *frame.Frame {
	if ctx.SplitPrompt != nil {
		ctx.SplitPrompt.Destroy()
		ctx.SplitPrompt = nil
		ctx.AttachPoint = nil
	}

	attachF := ctx.GetFocusedFrame()
	if attachF == nil {
		msgPrompt := prompt.NewMessage(ctx.X, prompt.DefaultMessageTheme, prompt.DefaultMessageConfig)
		timeout := 1 * time.Second
		msgPrompt.Show(ctx.Screens[0].ToXRect(), "Cannot split when not focused on a window", timeout, func(msg *prompt.Message) {})
		return nil
	}

	nprompt := prompt.NewInput(ctx.X, prompt.DefaultInputTheme, prompt.DefaultInputConfig)
	ctx.SplitPrompt = nprompt

	canc := func(inp *prompt.Input) {
		if ctx.SplitPrompt == nprompt {
			ctx.SplitPrompt.Destroy()
			ctx.SplitPrompt = nil
		}
	}

	resp := func(inp *prompt.Input, text string) {
		cmd := exec.Command("bash", "-c", text)
		err := cmd.Start()
		if err != nil {
			log.Println(err)
		}
		go func() {
			cmd.Wait()
		}()
		ctx.SplitPrompt.Destroy()
		ctx.SplitPrompt = nil
	}

	ctx.SplitPrompt.Show(ctx.Screens[0].ToXRect(), "Command:", resp, canc)
	return attachF
}

func RegisterSplitHooks(ctx *frame.Context) error {
	var err error
	// Builting shortcuts
	for k, v := range ctx.Config.BuiltinCommands {
		ncmd := v // force to not be a reference
		err = keybind.KeyReleaseFun(
			func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
				if ctx.Locked {
					return
				}

				cmd := exec.Command("bash", "-c", ncmd)
				err := cmd.Start()
				if err != nil {
					log.Println(err)
				}
				go func() {
					cmd.Wait()
				}()
			}).Connect(ctx.X, ctx.X.RootWin(), k, true)
		if err != nil {
			return err
		}
	}

	// Shitty standalone launchers
	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		if ctx.Locked {
			return
		}

		inPrompt := prompt.NewInput(X, prompt.DefaultInputTheme, prompt.DefaultInputConfig)

		canc := func(inp *prompt.Input) {
			// Chill
		}

		resp := func(inp *prompt.Input, text string) {
			cmd := exec.Command("bash", "-c", text)
			err = cmd.Start()
			if err != nil {
				log.Println(err)
			}
			go func() {
				cmd.Wait()
			}()
			inPrompt.Destroy()
		}

		inPrompt.Show(ctx.Screens[0].ToXRect(), "Command:", resp, canc)

	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.RunCmd, true)
	if err != nil {
		return err
	}

	// Split launchers
	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		if ctx.Locked {
			return
		}

		fr := Split(ctx)
		if fr == nil {
			return
		}
		ctx.AttachPoint = &frame.AttachTarget{
			Target: fr,
			Type:   frame.HORIZONTAL,
		}
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.SplitHorizontal, true)
	if err != nil {
		return err
	}

	err = keybind.KeyReleaseFun(func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
		if ctx.Locked {
			return
		}

		fr := Split(ctx)
		if fr == nil {
			return
		}
		ctx.AttachPoint = &frame.AttachTarget{
			Target: fr,
			Type:   frame.VERTICAL,
		}
	}).Connect(ctx.X, ctx.X.RootWin(), ctx.Config.SplitVertical, true)
	if err != nil {
		return err
	}
	return err
}
