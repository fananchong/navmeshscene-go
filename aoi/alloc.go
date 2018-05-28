package aoi

type Blocks struct {
	blockSize int
	mBlocks   []*[]QuadTreeNode
	mFreeList []*QuadTreeNode
}

func (this *Blocks) Init(blockSize int) {
	this.blockSize = blockSize
}

func (this *Blocks) New() *QuadTreeNode {
LABLE_DO:
	l := len(this.mFreeList)
	if l != 0 {
		item := this.mFreeList[l-1]
		this.mFreeList = this.mFreeList[:l-1]
		return item
	} else {
		this.newBlock()
		goto LABLE_DO
	}
}

func (this *Blocks) Delete(node *QuadTreeNode) {
	this.mFreeList = append(this.mFreeList, node)
}

func (this *Blocks) newBlock() {
	items := make([]QuadTreeNode, this.blockSize)
	this.mBlocks = append(this.mBlocks, &items)
	for _, item := range items {
		this.mFreeList = append(this.mFreeList, &item)
	}
}
