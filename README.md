# howm

# Installation
To install, run `sudo ./install.sh`. It's generally not wise to run random shell scripts with sudo but this one should be pretty easy to inspect. Dokcer is a prerequisite for the script, but if you already have go installed it should be fairly easy to do it without docker.

# Default Configuration
All key combinations described here are easily editable in `config.go`. `Mod4` is used a lot, which is commonly assigned to the Windows key.

#### Logging out
Press `Mod4-Backspace`

#### Builtin Commands
Some commands have builtin keyboard shortcuts, namely;
`Mod4-t: x-terminal-emulator`
`Mod4-w: google-chrome`
`Mod4-p: XDG_CURRENT_DESKTOP=GNOME gnome-control-center`

More commands can be added to the builtin commands in `config.go`.

`Mod4-c` brings up a dialog to run an arbitrary command.

#### Window controls
`Mod4-up/left/right/down` will move the window to anchor points around the screen, as well as keep moving them across other screens if they are available.

`Mod4-d` will close a frame.

`Mod4-x` will toogle a frame as expanded, making it the same size as its container.

`Mod4-h` will hide the container decorations.

`Alt-Tab`/`Alt-Shift-Tab` works like you'd expect.

If an internal video is being fullscreened, sometimes you may need to resize or move the window a little to have the internal video fill the screen.

#### Taskbar
The taskbar can be toggled with `Mod4-s`. `Mod4-Shift-down` will minimize the focused container to the taskbar. `Mod4-Shift-left/right` will scroll the taskbar if lots of windows are open.

The display time format can be changed in `config.go`, but things may be a bit funky if the time format does not have constant size.

#### Background image
Any photo named `$USER/.config/howm/bg.png` will be used as a background image. You can change the path in `config.go`.

#### Splitting
To subdivide a window, press `Mod4-e` for a horizontal split or `Mod4-r` for a vertical split. A command window will pop up to take in a command to launch, but you can use the keyboard shortcuts to launch a builtin command, bypassing the command prompt.

#### Volume
Can be controlled with `Mod4-F1` to mute, and `Mod4-F2/F3` to raise/lower volume.

#### Brightness
Can be raised/lowered with `Mod4-F11/F12`. If brightness controls are not working, check what folder exists in `/sys/class/backlight/` and change the Backlight value in `config.go` to match.

# Monitor hotplugging
If you are planning on hotplugging monitors, it is recommended you install arandr and autorandr.

`arandr` is a graphical utility to help you configure your screens. Once you have them configured how you like, use `autorandr --save <profile>` to save the configuration and `autorandr --default <profile>` to make it default. Once this is done the monitor setup should stay sticky through disconnects and reconnects.

# Logging
To tail the logs run `journalctl -f /usr/lib/gdm3/gdm-x-session`.