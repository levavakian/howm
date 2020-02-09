/*
rowm is a hybrid tiling manager aiming to replicate the ergonomics of terminator and tmux.
For more information on the usage, please read the README instead.

main.go - entrypoint, initializes connection to X and connects root callbacks
root/ - callbacks for the root window are defined here
frame/ - the majority of the logic is in this folder, defines the tree window structure and wrapping decorations
sideloop/ - utilities for running code in line with the X event loop from xgbutil
resources/ - non compiled data resources
cmd/rowmbright/ - a utility for changing backlight brightness without root privileges
ext/ - misc bits and bobs missing in either the standard library or xgbutil

dev.sh - starts up a container for developing in (required for running compile.sh or test.sh)
compile.sh - go gets/builds rowm in container
clean.sh - gets rid of container created by dev.sh
exec.sh - execs into container created by dev.sh
install.sh - (requires root) builds rowm and installs files to the system (only tested on ubuntu18.04)
test.sh - (requires xephyr to be installed and dev.sh to have been run) opens a nested x server with rowm loaded for trying out
*/
package main
