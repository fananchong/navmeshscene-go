package aoi

type EQuadrant int

const (
	UnknowQuadrant EQuadrant = 0
	RightTop       EQuadrant = 1 // 右上：象限一
	LeftTop        EQuadrant = 2 // 左上：象限二
	LeftDown       EQuadrant = 3 // 左下：象限三
	RightDown      EQuadrant = 4 // 右下：象限四
)

type Rect struct {
	mLeft   float32
	mRight  float32
	mBottom float32
	mTop    float32
	mMidX   float32
	mMidY   float32
}

func NewRect(left, right, bottom, top float32) *Rect {
	this := &Rect{}
	this.Init(left, right, bottom, top)
	return this
}

func (this *Rect) Init(left, right, bottom, top float32) {
	this.mLeft = left
	this.mRight = right
	this.mTop = top
	this.mBottom = bottom
	this.mMidX = left + (right-left)/2
	this.mMidY = bottom + (top-bottom)/2
}

func (this *Rect) Left() float32   { return this.mLeft }
func (this *Rect) Right() float32  { return this.mRight }
func (this *Rect) Bottom() float32 { return this.mBottom }
func (this *Rect) Top() float32    { return this.mTop }
func (this *Rect) MidX() float32   { return this.mMidX }
func (this *Rect) MidY() float32   { return this.mMidY }

func (this *Rect) ContainsRect(rect *Rect) bool {
	return (this.mLeft <= rect.mLeft &&
		this.mBottom <= rect.mBottom &&
		rect.mRight <= this.mRight &&
		rect.mTop <= this.mTop)
}

func (this *Rect) Contains(point *Point) bool {
	return (point.X >= this.mLeft && point.X <= this.mRight &&
		point.Y >= this.mBottom && point.Y <= this.mTop)
}

func (this *Rect) ContainsItem(item *Object) bool {
	return (item.X >= this.mLeft && item.X <= this.mRight &&
		item.Y >= this.mBottom && item.Y <= this.mTop)
}

func (this *Rect) Intersects(rect *Rect) bool {
	return !(this.mRight < rect.mLeft ||
		rect.mRight < this.mLeft ||
		this.mTop < rect.mBottom ||
		rect.mTop < this.mBottom)
}

func (this *Rect) GetQuadrant(item *Object) EQuadrant {
	if this.ContainsItem(item) {
		if item.Y >= this.mMidY {
			if item.X >= this.mMidX {
				return RightTop
			} else {
				return LeftTop
			}
		} else {
			if item.X >= this.mMidX {
				return RightDown
			} else {
				return LeftDown
			}
		}
	} else {
		return UnknowQuadrant
	}
}

func (this *Rect) GetQuadrantWithoutBounds(item *Object) EQuadrant {
	if item.Y >= this.mMidY {
		if item.X >= this.mMidX {
			return RightTop
		} else {
			return LeftTop
		}
	} else {
		if item.X >= this.mMidX {
			return RightDown
		} else {
			return LeftDown
		}
	}
}
