package aoi

type IItem interface {
	getPostion() *Point
	getItemNext() IItem
	setItemNext(item IItem)
	setQueryNext(item IItem)
	getNode() *QuadTreeNode
	setNode(node *QuadTreeNode)
}

type QuadTree struct {
	NodeCapacity int          // 节点容量
	LevelLimit   int          // 层数限制
	mRoot        QuadTreeNode // 根节点
	mAlloc       Blocks       // 节点分配器
}

func (this *QuadTree) Insert(item IItem) bool { return this.mRoot.Insert(item) }

func (this *QuadTree) Remove(item IItem) bool {
	node := item.getNode()
	if node != nil {
		return node.Remove(item)
	} else {
		return false
	}
}

func (this *QuadTree) Query1(area *Rect) IItem {
	var head, tail IItem
	this.mRoot.Query(area, &head, &tail)
	if tail != nil {
		tail.setQueryNext(nil)
	}
	return head
}

func (this *QuadTree) Query2(source IItem, radius float32) IItem {
	var area Rect
	pos := source.getPostion()
	area.Init(pos.X-radius, pos.X+radius, pos.Y-radius, pos.Y+radius)
	return this.Query4(source, &area)
}

func (this *QuadTree) Query3(source IItem, halfExtentsX, halfExtentsY float32) IItem {
	var area Rect
	pos := source.getPostion()
	area.Init(pos.X-halfExtentsX, pos.X+halfExtentsX, pos.Y-halfExtentsY, pos.Y+halfExtentsY)
	return this.Query4(source, &area)
}

func (this *QuadTree) Query4(source IItem, area *Rect) IItem {
	var head, tail IItem
	node := source.getNode()
	if node != nil && node.mBounds.ContainsRect(area) {
		node.Query(area, &head, &tail)
	} else {
		this.mRoot.Query(area, &head, &tail)
	}
	if tail != nil {
		tail.setQueryNext(nil)
	}
	return head
}

func (this *QuadTree) Update(item IItem) bool {
	node := item.getNode()
	if node != nil {
		if node.mBounds.Contains(item.getPostion()) {
			return true
		}
		node.Remove(item)
	}
	return this.mRoot.Insert(item)
}

func (this *QuadTree) GetItemCount() int    { return this.mRoot.GetItemCount() }
func (this *QuadTree) GetBounds() *Rect     { return &this.mRoot.mBounds }
func (this *QuadTree) SetBounds(rect *Rect) { this.mRoot.mBounds = *rect }
