/*
frame is the package that contains the definition of the tree window structure and decoration handling that drives rowm.

place.go - entrypoint for when a new window is created either standalone or as part of a split
frame.go - defines the tree structure and traversal of windows
container.go - defines the resizing, minimizing, and moving of a window tree as wrapped by decorations
context.go - all non trivial state is stored in the context, and is available to most operations
config.go - store of all user defined settings
decoration.go - utilities for decorations (non user created windows)
pieces.go - definitions of individual decorations and their callbacks which make up a container
taskbar.go - a taskbar decoration for displaying basic system information and showing open windows
anchor.go - utilities for defining screen anchors (preset shapes on a screen you can hotkey to)
background.go - utilities for generating backgrounds
rect.go - a basic rectangle definition for describing window shapes and locations
*/
package frame
