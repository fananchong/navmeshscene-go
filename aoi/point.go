package aoi

type Point struct {
	X float32
	Y float32
}

func (this *Point) getPostion() *Point {
	return this
}

type Size Point