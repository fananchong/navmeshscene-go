package aoi

type Object struct {
	Point
	mNode     *QuadTreeNode
	QueryNext *Object
}

func NewScene(bounds *Rect, nodeCapacity, levelLimit int) *QuadTree {
	this := &QuadTree{}
	this.NodeCapacity = nodeCapacity
	this.LevelLimit = levelLimit
	this.mRoot.Init(this, NodeTypeLeaf, 0, bounds, nil)
	return this
}

func NewDefaultSecene(bounds *Rect) *QuadTree {
	return NewScene(bounds, 16, 10)
}
