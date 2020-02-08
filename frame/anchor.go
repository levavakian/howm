package frame

type AnchorType int
const (
	NONE AnchorType = iota
	FULL
	TOP
	LEFT
	RIGHT
	BOTTOM
)

func AnchorShape(ctx *Context, screen Rect, anchor AnchorType) Rect {
	if !ctx.Taskbar.Hidden && screen == ctx.Screens[0] {
		screen.H = screen.H - ctx.Config.TaskbarHeight
	}

	if anchor == TOP {
		screen.H = screen.H / 2
		return screen
	}

	if anchor == BOTTOM {
		origYEnd := screen.Y + screen.H
		screen.Y = screen.Y + screen.H/2
		screen.H = origYEnd - screen.Y
		return screen
	}

	if anchor == LEFT {
		screen.W = screen.W / 2
		return screen
	}

	if anchor == RIGHT {
		origXEnd := screen.X + screen.W
		screen.X = screen.X + screen.W/2
		screen.W = origXEnd - screen.X
		return screen
	}

	return screen
}

func AnchorMatch(ctx *Context, screen Rect, shape Rect) AnchorType {
	options := []AnchorType{
		FULL,
		TOP,
		LEFT,
		RIGHT,
		BOTTOM,
	}
	for _, opt := range options {
		if shape == AnchorShape(ctx, screen, opt) {
			return opt
		}
	}
	return NONE
}