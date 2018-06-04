package tests

import (
	"fmt"
	"math"
	"math/rand"

	NavMeshScene "github.com/fananchong/navmeshscene-go"
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
	myassert(scn.Insert(&temp.Object))
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
		myassert(scn.Remove(&temp.Object))
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
		if queryArea.ContainsItem(&items[i].Object) {
			testCount++
		}
	}

	findCount := 0
	item := scn.Query1(&queryArea)
	for item != nil {
		findCount++
		item = item.QueryNext
	}
	//fmt.Printf("find obj count:%d, test count:%d, total count:%d\n", findCount, testCount, scn.GetItemCount())
	myassert(testCount == findCount)
}

func Query_by_radius(scn *aoi.QuadTree, items []*A, radius float32) {
	if len(items) == 0 {
		return
	}

	testCount := 0

	index := rand.Int() % len(items)
	var queryArea aoi.Rect
	queryArea.Init(
		items[index].X-radius,
		items[index].X+radius,
		items[index].Y-radius,
		items[index].Y+radius)

	for i := 0; i < len(items); i++ {
		if queryArea.ContainsItem(&items[i].Object) {
			testCount++
		}
	}

	findCount := 0
	item := scn.Query1(&queryArea)
	for item != nil {
		findCount++
		item = item.QueryNext
	}
	//fmt.Printf("find obj count:%d, test count:%d, total count:%d\n", findCount, testCount, scn.GetItemCount())
	myassert(testCount == findCount)
}

/// =========================================

type Player struct {
	*NavMeshScene.Agent
}

func NewPlayer() *Player {
	this := &Player{}
	this.Agent = NavMeshScene.NewAgent()
	this.Agent.OnHit = this.OnHit
	return this
}

func (this *Player) OnHit(agentId uint64) {
	this.ChangeDir()
}

func (this *Player) ChangeDir() {
	angle := float64(rand.Int() % 360)
	vx := math.Cos(math.Pi * angle / 180)
	vy := -math.Sin(math.Pi * angle / 180)
	s := math.Sqrt(vx*vx + vy*vy)
	vx = vx / s
	vy = vy / s
	v := [3]float32{float32(vx * 5), 0, float32(vy * 5)}
	this.Velocity = v
}

func InitScene(scene *NavMeshScene.Scene, path string) {
	const PLAYER_COUNT = 5000
	if ec := scene.Load(path); ec != 0 {
		panic(ec)
	}

	fmt.Println("load scene success!")

	min := scene.GetBoundsMin()
	max := scene.GetBoundsMax()
	fmt.Printf("width:%f, heigth:%f\n", max[0]-min[0], max[2]-min[2])

	for i := 0; i < PLAYER_COUNT; i++ {
		player := NewPlayer()
		scene.AddAgent(uint64(i+1), player.Agent)
		player.RandomPosition()
		player.ChangeDir()
	}
	fmt.Printf("player count: %d\n", PLAYER_COUNT)
}
