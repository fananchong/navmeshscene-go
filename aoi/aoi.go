package aoi

type Object struct {
	Point
	mNode      *QuadTreeNode
	mQueryNext IItem
}

func (this *Object) Next() IItem                { return this.mQueryNext }
func (this *Object) setQueryNext(item IItem)    { this.mQueryNext = item }
func (this *Object) getNode() *QuadTreeNode     { return this.mNode }
func (this *Object) setNode(node *QuadTreeNode) { this.mNode = node }

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
