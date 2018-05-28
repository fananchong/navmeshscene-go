package aoi

type ENodeType int

const (
	NodeTypeNormal ENodeType = 0 // 非叶节点
	NodeTypeLeaf   ENodeType = 1 // 叶子节点
)

const ChildrenNum int = 4

type QuadTreeNode struct {
	mNodeType  ENodeType                  // 节点类型
	mLevel     int                        // 当前节点所在层级
	mBounds    Rect                       // 节点边框范围
	mParent    *QuadTreeNode              // 父节点
	mChildrens [ChildrenNum]*QuadTreeNode // 孩子节点
	mItemCount int                        // 叶子节点上的Item数量
	mItems     IItem                      // 叶子节点上的Items
	mTree      *QuadTree                  // 所在树
}

func (this *QuadTreeNode) Init(tree *QuadTree, t ENodeType, lvl int, rect *Rect, parent *QuadTreeNode) {
	this.mNodeType = t
	this.mLevel = lvl
	this.mBounds = *rect
	this.mParent = parent
	for i := 0; i < ChildrenNum; i++ {
		this.mChildrens[i] = nil
	}
	this.mItemCount = 0
	this.mItems = nil
	this.mTree = tree
}

func (this *QuadTreeNode) Insert(item IItem) bool {
LABLE_NORMAL:
	if this.mNodeType == NodeTypeNormal {
		index := this.mBounds.GetQuadrant(item.getPostion()) - 1
		if index >= 0 {
			return this.mChildrens[index].Insert(item)
		} else {
			return false
		}
	} else {
		if this.mItemCount < this.mTree.NodeCapacity {
			if this.mBounds.ContainsItem(item) {
				this.mItemCount++
				item.setItemNext(this.mItems)
				this.mItems = item
				item.setNode(this)
				return true
			} else {
				return false
			}
		} else {
			if this.mLevel+1 >= this.mTree.LevelLimit {
				return false
			}
			this.split()
			goto LABLE_NORMAL
		}
	}
}

func (this *QuadTreeNode) split() {
	this.mNodeType = NodeTypeNormal

	// 第一象限，右上
	var rect0 Rect
	rect0.Init(this.mBounds.MidX(), this.mBounds.Right(), this.mBounds.MidY(), this.mBounds.Top())

	// 第二象限，左上
	var rect1 Rect
	rect1.Init(this.mBounds.Left(), this.mBounds.MidX(), this.mBounds.MidY(), this.mBounds.Top())

	// 第三象限，左下
	var rect2 Rect
	rect2.Init(this.mBounds.Left(), this.mBounds.MidX(), this.mBounds.Bottom(), this.mBounds.MidY())

	// 第四象限，右下
	var rect3 Rect
	rect3.Init(this.mBounds.MidX(), this.mBounds.Right(), this.mBounds.Bottom(), this.mBounds.MidY())

	this.mChildrens[0] = this.mTree.mAlloc.New()
	this.mChildrens[0].Init(this.mTree, NodeTypeLeaf, this.mLevel+1, &rect0, this)
	this.mChildrens[1] = this.mTree.mAlloc.New()
	this.mChildrens[1].Init(this.mTree, NodeTypeLeaf, this.mLevel+1, &rect1, this)
	this.mChildrens[2] = this.mTree.mAlloc.New()
	this.mChildrens[2].Init(this.mTree, NodeTypeLeaf, this.mLevel+1, &rect2, this)
	this.mChildrens[3] = this.mTree.mAlloc.New()
	this.mChildrens[3].Init(this.mTree, NodeTypeLeaf, this.mLevel+1, &rect3, this)

	for it := this.mItems; it != nil; {
		head := it.getItemNext()
		index := this.mBounds.GetQuadrantWithoutBounds(it.getPostion()) - 1
		this.mChildrens[index].Insert(it)
		it = head
	}
	this.mItemCount = 0
	this.mItems = nil
}

func (this *QuadTreeNode) Remove(item IItem) bool {
	var pre IItem
	it := this.mItems
	for it != nil {
		head := it.getItemNext()
		if it == item {
			this.mItemCount--
			if pre != nil {
				pre.setItemNext(it.getItemNext())
			} else {
				this.mItems = this.mItems.getItemNext()
			}
			this.tryMerge()
			return true
		} else {
			pre = it
		}
		it = head
	}
	return false
}

func (this *QuadTreeNode) tryMerge() {
	node := this.mParent
	for node != nil {
		count := 0
		childrens := node.mChildrens
		for i := 0; i < ChildrenNum; i++ {
			if childrens[i].mNodeType != NodeTypeLeaf {
				return
			}
			count += childrens[i].mItemCount
		}

		if count <= this.mTree.NodeCapacity {
			node.mNodeType = NodeTypeLeaf
			node.mItemCount = 0
			node.mItems = nil
			for i := 0; i < ChildrenNum; i++ {
				it := childrens[i].mItems
				for it != nil {
					head := it.getItemNext()
					node.mItemCount++
					it.setItemNext(node.mItems)
					node.mItems = it
					it.setNode(node)
					it = head
				}
				this.mTree.mAlloc.Delete(childrens[i])
			}
			node = node.mParent
		} else {
			break
		}
	}
}

func (this *QuadTreeNode) Query(area *Rect, head, tail *IItem) {
	if this.mNodeType == NodeTypeNormal {
		for i := 0; i < ChildrenNum; i++ {
			if area.Intersects(&(this.mChildrens[i].mBounds)) {
				this.mChildrens[i].Query(area, head, tail)
			}
		}
	} else {
		for it := this.mItems; it != nil; it = it.getItemNext() {
			if area.ContainsItem(it) {
				if (*head) != nil {
					(*tail).setQueryNext(it)
					*tail = it
				} else {
					*head = it
					*tail = it
				}
			}
		}
	}
}

func (this *QuadTreeNode) GetItemCount() int {
	count := 0
	if this.mNodeType == NodeTypeNormal {
		for i := 0; i < ChildrenNum; i++ {
			count += this.mChildrens[i].GetItemCount()
		}
	} else {
		count += this.mItemCount
	}
	return count
}
