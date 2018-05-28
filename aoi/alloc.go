package aoi

type Blocks struct {
	blockSize int
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
	for i := 0; i < this.blockSize; i++ {
		this.mFreeList = append(this.mFreeList, &QuadTreeNode{})
	}
}
