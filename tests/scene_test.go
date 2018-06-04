package tests

import (
	"math/rand"
	"testing"
	"time"

	NavMeshScene "github.com/fananchong/navmeshscene-go"
)

const TEST_COUNT = 10000

func Test_scene(t *testing.T) {
	const path1 = "Meshes/scene1.obj.tile.bin"
	const path2 = "Meshes/scene1.obj.tilecache.bin"

	rand.Seed(time.Now().UTC().UnixNano())
	scn1 := NavMeshScene.NewStaticScene()
	InitScene(scn1.Scene, path1)
	for i := 0; i < TEST_COUNT; i++ {
		scn1.Simulation(0.025)
	}

	scn2 := NavMeshScene.NewDynamicScene(NavMeshScene.HEIGHT_MODE_1)
	InitScene(scn2.Scene, path1)
	for i := 0; i < TEST_COUNT; i++ {
		scn2.Simulation(0.025)
	}

	scn3 := NavMeshScene.NewDynamicScene(NavMeshScene.HEIGHT_MODE_2)
	InitScene(scn3.Scene, path2)
	for i := 0; i < TEST_COUNT; i++ {
		scn3.Simulation(0.025)
	}
}
