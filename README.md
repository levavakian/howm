# howm

# Monitor hotplugging
If you are planning on hotplugging monitors, it is recommended you install arandr and autorandr.

`arandr` is a graphical utility to help you configure your screens. Once you have them configured how you like, use `autorandr --save <profile>` to save the configuration and `autorandr --default <profile>` to make it default. Once this is done the monitor setup should stay sticky through disconnects and reconnects.

# Brightness
If brightness controls are not working, check what folder exists in `/sys/class/backlight/` and change the Backlight value in `config.go` to match.

# Logging
To tail the logs run `journalctl -f /usr/lib/gdm3/gdm-x-session`.