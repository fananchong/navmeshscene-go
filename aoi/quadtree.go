package aoi

type QuadTree struct {
	NodeCapacity int          // 节点容量
	LevelLimit   int          // 层数限制
	mRoot        QuadTreeNode // 根节点
}

func (this *QuadTree) Insert(item *Object) bool { return this.mRoot.Insert(item) }

func (this *QuadTree) Remove(item *Object) bool {
	node := item.mNode
	if node != nil {
		return node.Remove(item)
	} else {
		return false
	}
}

func (this *QuadTree) Query1(area *Rect) *Object {
	var head, tail *Object
	this.mRoot.Query(area, &head, &tail)
	if tail != nil {
		tail.QueryNext = nil
	}
	return head
}

func (this *QuadTree) Query2(source *Object, radius float32) *Object {
	var area Rect
	area.Init(source.X-radius, source.X+radius, source.Y-radius, source.Y+radius)
	return this.Query4(source, &area)
}

func (this *QuadTree) Query3(source *Object, halfExtentsX, halfExtentsY float32) *Object {
	var area Rect
	area.Init(source.X-halfExtentsX, source.X+halfExtentsX, source.Y-halfExtentsY, source.Y+halfExtentsY)
	return this.Query4(source, &area)
}

func (this *QuadTree) Query4(source *Object, area *Rect) *Object {
	var head, tail *Object
	node := source.mNode
	if node != nil && node.mBounds.ContainsRect(area) {
		node.Query(area, &head, &tail)
	} else {
		this.mRoot.Query(area, &head, &tail)
	}
	if tail != nil {
		tail.QueryNext = nil
	}
	return head
}

func (this *QuadTree) Update(item *Object) bool {
	node := item.mNode
	if node != nil {
		if node.mBounds.ContainsItem(item) {
			return true
		}
		node.Remove(item)
	}
	return this.mRoot.Insert(item)
}

func (this *QuadTree) GetItemCount() int    { return this.mRoot.GetItemCount() }
func (this *QuadTree) GetBounds() *Rect     { return &this.mRoot.mBounds }
func (this *QuadTree) SetBounds(rect *Rect) { this.mRoot.mBounds = *rect }
