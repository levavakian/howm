# howm

# Monitor hotplugging
If you are planning on hotplugging monitors, it is recommended you install arandr and autorandr.

`arandr` is a graphical utility to help you configure your screens. Once you have them configured how you like, use `autorandr --save <profile>` to save the configuration and `autorandr --default <profile>` to make it default. Once this is done the monitor setup should stay sticky through disconnects and reconnects.

# Logging
To tail the logs run `journalctl -f /usr/lib/gdm3/gdm-x-session`