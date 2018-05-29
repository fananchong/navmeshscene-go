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
	mItems     []IItem                    // 叶子节点上的Items
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
	this.mItems = this.mItems[:0]
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
		if len(this.mItems) < this.mTree.NodeCapacity {
			if this.mBounds.ContainsItem(item) {
				item.setNode(this)
				this.mItems = append(this.mItems, item)
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

	this.mChildrens[0] = &QuadTreeNode{}
	this.mChildrens[0].Init(this.mTree, NodeTypeLeaf, this.mLevel+1, &rect0, this)
	this.mChildrens[1] = &QuadTreeNode{}
	this.mChildrens[1].Init(this.mTree, NodeTypeLeaf, this.mLevel+1, &rect1, this)
	this.mChildrens[2] = &QuadTreeNode{}
	this.mChildrens[2].Init(this.mTree, NodeTypeLeaf, this.mLevel+1, &rect2, this)
	this.mChildrens[3] = &QuadTreeNode{}
	this.mChildrens[3].Init(this.mTree, NodeTypeLeaf, this.mLevel+1, &rect3, this)

	for _, item := range this.mItems {
		index := this.mBounds.GetQuadrantWithoutBounds(item.getPostion()) - 1
		this.mChildrens[index].Insert(item)
	}
	this.mItems = this.mItems[:0]
}

func (this *QuadTreeNode) Remove(item IItem) bool {
	for index, it := range this.mItems {
		if it == item {
			this.mItems = append(this.mItems[:index], this.mItems[index+1:]...)
			this.tryMerge()
			return true
		}
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
			count += len(childrens[i].mItems)
		}

		if count <= this.mTree.NodeCapacity {
			node.mNodeType = NodeTypeLeaf
			node.mItems = node.mItems[:0]
			for i := 0; i < ChildrenNum; i++ {
				for _, item := range childrens[i].mItems {
					item.setNode(node)
					node.mItems = append(node.mItems, item)
				}
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
		for _, item := range this.mItems {
			if area.ContainsItem(item) {
				if (*head) != nil {
					(*tail).setQueryNext(item)
					*tail = item
				} else {
					*head = item
					*tail = item
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
		count += len(this.mItems)
	}
	return count
}
