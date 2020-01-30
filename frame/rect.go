package frame

import (
	"howm/ext"
)

type Rect struct {
	X int
	Y int
	W int
	H int
}

type Rectf struct {
	X float64
	Y float64
	W float64
	H float64
}

func AreaOfIntersection(shapeA, shapeB Rect) int {
	xminmax := ext.IMin(shapeA.X + shapeA.W, shapeB.X + shapeB.W)
	xmaxmin := ext.IMax(shapeA.X, shapeB.X)
	yminmax := ext.IMin(shapeA.Y + shapeA.H, shapeB.Y + shapeB.H)
	ymaxmin := ext.IMax(shapeA.Y, shapeB.Y)

	dx := xminmax - xmaxmin
	dy := yminmax - ymaxmin

	if dx >= 0 && dy >= 0 {
		return dx * dy
	} else {
		return 0
	}
}