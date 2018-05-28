package tests

import (
	"math"
	"math/rand"

	"github.com/fananchong/navmeshscene-go/aoi"
)

func myassert(cond bool) {
	if cond == false {
		panic("")
	}
}

func randArray(src []*A) []*A {
	dest := make([]*A, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}

type A struct {
	aoi.Object
}

func Add(scn *aoi.QuadTree, items *[]*A) {
	x := float32(rand.Int()%int(scn.GetBounds().Right()-scn.GetBounds().Left())) + scn.GetBounds().Left()
	y := float32(rand.Int()%int(scn.GetBounds().Top()-scn.GetBounds().Bottom())) + scn.GetBounds().Bottom()
	temp := &A{}
	temp.X = x
	temp.Y = y
	myassert(scn.Insert(temp))
	*items = append(*items, temp)
}

func _test_delete(scn *aoi.QuadTree, items *[]*A, count int) {
	itemsNum := len(*items)
	if itemsNum == 0 {
		return
	}
	*items = randArray(*items)
	for i := 0; i < int(math.Min(float64(count), float64(itemsNum))); i++ {
		temp := (*items)[len(*items)-1]
		myassert(scn.Remove(temp))
		*items = (*items)[:len(*items)-1]
	}
}

func Query(scn *aoi.QuadTree, items []*A) {
	testCount := 0
	var queryArea aoi.Rect
	queryArea.Init(
		float32(rand.Int()%10),
		float32(rand.Int()%int(scn.GetBounds().Right()-scn.GetBounds().Left()))+scn.GetBounds().Left(),
		float32(rand.Int()%10),
		float32(rand.Int()%int(scn.GetBounds().Top()-scn.GetBounds().Bottom()))+scn.GetBounds().Bottom())

	for i := 0; i < len(items); i++ {
		if queryArea.ContainsItem(items[i]) {
			testCount++
		}
	}

	findCount := 0
	item := scn.Query1(&queryArea)
	for item != nil {
		findCount++
		item = item.Next()
	}
	//fmt.Printf("find obj count:%d, test count:%d, total count:%d\n", findCount, testCount, scn.GetItemCount())
	myassert(testCount == findCount)
}

func Query_by_radius(scn *aoi.QuadTree, items []*A, radius float32) {
	testCount := 0

	index := rand.Int() % len(items)
	var queryArea aoi.Rect
	queryArea.Init(
		items[index].X-radius,
		items[index].X+radius,
		items[index].Y-radius,
		items[index].Y+radius)

	for i := 0; i < len(items); i++ {
		if queryArea.ContainsItem(items[i]) {
			testCount++
		}
	}

	findCount := 0
	item := scn.Query1(&queryArea)
	for item != nil {
		findCount++
		item = item.Next()
	}
	//fmt.Printf("find obj count:%d, test count:%d, total count:%d\n", findCount, testCount, scn.GetItemCount())
	myassert(testCount == findCount)
}
