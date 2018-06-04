package tests

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	NavMeshScene "github.com/fananchong/navmeshscene-go"
)

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

const PLAYER_COUNT = 5000
const TEST_COUNT = 10000

func test(scene *NavMeshScene.Scene, path string) {
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

	for i := 0; i < TEST_COUNT; i++ {
		scene.Simulation(0.025)
	}
}

func Test_scene(t *testing.T) {
	const path1 = "Meshes/scene1.obj.tile.bin"
	const path2 = "Meshes/scene1.obj.tilecache.bin"

	rand.Seed(time.Now().UTC().UnixNano())

	scn1 := NavMeshScene.NewStaticScene()
	test(scn1.Scene, path1)

	scn2 := NavMeshScene.NewDynamicScene(NavMeshScene.HEIGHT_MODE_1)
	test(scn2.Scene, path1)

	scn3 := NavMeshScene.NewDynamicScene(NavMeshScene.HEIGHT_MODE_2)
	test(scn3.Scene, path2)
}
