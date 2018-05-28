package tests

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fananchong/navmeshscene-go/aoi"
)

func test1() {
	var rect aoi.Rect
	rect.Init(0, 1000, 0, 1000)
	scn := aoi.NewDefaultSecene(&rect)

	// 测试插入
	var items []*A
	for i := 0; i < 8192; i++ {
		Add(scn, &items)
	}

	// 测试查询
	for i := 0; i < 1000; i++ {
		Query(scn, items)
		Query_by_radius(scn, items, float32(rand.Int()%200+50))
	}

	// 测试删除
	itemsNum := len(items)
	_test_delete(scn, &items, itemsNum)
	items = nil
	//fmt.Printf("delete obj count:%d, total count:%d\n", itemsNum, scn.GetItemCount())
}

func test2() {
	var rect aoi.Rect
	rect.Init(0, 1000, 0, 1000)
	scn := aoi.NewDefaultSecene(&rect)

	var items []*A
	for {
		op := rand.Int() % 10

		if op <= 6 {
			Add(scn, &items)
		} else if op <= 8 {
			itemsNum := len(items)
			_test_delete(scn, &items, itemsNum%3+1)
		} else {
			Query(scn, items)
			Query_by_radius(scn, items, float32(rand.Int()%200+50))
		}
	}
}

func Test_aoi(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	test1()
	test2()
}
