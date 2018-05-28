package aoi

type Object struct {
	Point
	mNode      *QuadTreeNode
	mQueryNext IItem
	mItemNext  IItem
}

func (this *Object) Next() IItem                { return this.mQueryNext }
func (this *Object) getItemNext() IItem         { return this.mItemNext }
func (this *Object) setItemNext(item IItem)     { this.mItemNext = item }
func (this *Object) setQueryNext(item IItem)    { this.mQueryNext = item }
func (this *Object) getNode() *QuadTreeNode     { return this.mNode }
func (this *Object) setNode(node *QuadTreeNode) { this.mNode = node }

func NewScene(bounds *Rect, nodeCapacity, levelLimit, blockSize int) *QuadTree {
	this := &QuadTree{}
	this.NodeCapacity = nodeCapacity
	this.LevelLimit = levelLimit
	this.mAlloc.Init(blockSize)
	this.mRoot.Init(this, NodeTypeLeaf, 0, bounds, nil)
	return this
}

func NewDefaultSecene(bounds *Rect) *QuadTree {
	return NewScene(bounds, 16, 10, 1024)
}
