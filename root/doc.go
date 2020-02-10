/*
root contains the files that define the callbacks on the root window

base.go - minimal callbacks such as locking the screen, focusing, shutting down the window manager, and creating a window
brightness.go - callbacks for raising/lowering the backlight
choose.go - callbacks implementing an alt-tab like interface
launchers.go - callbacks for prompts that launch new windows (including the paritioning launch)
monitor.go - a side event loop that monitors for changes of the screen configuration
volume.go - callbacks for raising/lowering/muting volume
taskbar.go - callbacks for interacting with the taskbar
*/
package root
