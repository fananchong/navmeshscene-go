package NavMeshScene

import (
	"fmt"

	"github.com/fananchong/navmeshscene-go/aoi"
)

type Scene struct {
	*aoi.QuadTree
	Detour *Detour
	Agents map[uint64]*Agent
}

func NewScene(bStatic bool) *Scene {
	this := &Scene{}
	this.QuadTree = aoi.NewDefaultSecene(aoi.NewRect(0, 1, 0, 1))
	this.Agents = make(map[uint64]*Agent)
	this.Detour = NewDetour(bStatic, 16)
	return this
}

func (this *Scene) Load(path string) int {
	retCode := this.Detour.Load(path)
	if retCode == 0 {
		bmin := this.Detour.mBoundsMin
		bmax := this.Detour.mBoundsMax
		fmt.Printf("bounds min:(%f, %f, %f)\n", bmin[0], bmin[1], bmin[2])
		fmt.Printf("bounds max:(%f, %f, %f)\n", bmax[0], bmax[1], bmax[2])
		this.QuadTree.SetBounds(aoi.NewRect(bmin[0], bmax[0], bmin[2], bmax[2]))
	}
	return retCode
}

func (this *Scene) Simulation(delta float32) {
	if this.Detour.mTileCache != nil {
		this.Detour.mTileCache.Update(delta, this.Detour.mMesh, nil)
	}
	for _, agent := range this.Agents {
		agent.Update(delta)
	}
}

func (this *Scene) AddAgent(id uint64, agent *Agent) {
	if id != 0 && agent != nil {
		agent.Id = id
		agent.Object.UserData = id
		this.Agents[id] = agent
		agent.Scene = this
	}
}

func (this *Scene) RemoveAgent(id uint64) {
	if agent, ok := this.Agents[id]; ok {
		this.Remove(&agent.Object)
		delete(this.Agents, id)
	}
}

func (this *Scene) GetBoundsMin() []float32 {
	return this.Detour.mBoundsMin[:]
}

func (this *Scene) GetBoundsMax() []float32 {
	return this.Detour.mBoundsMax[:]
}

type StaticScene struct {
	*Scene
}

func NewStaticScene() *StaticScene {
	this := &StaticScene{}
	this.Scene = NewScene(true)
	return this
}

type DynamicScene struct {
	*Scene
}

const HEIGHT_MODE_1 int = 1 // 原始的，精度不是很高，但是没多余消耗。可以通过使Tile Size变小来提高精度
const HEIGHT_MODE_2 int = 2 // 通过公共的StaticScene上，获取精确高度值。

func NewDynamicScene(heightMode int) *DynamicScene {
	this := &DynamicScene{}
	this.Scene = NewScene(false)
	this.Scene.Detour.mHeightMode = heightMode
	return this
}

func (this *DynamicScene) AddCapsuleObstacle(pos []float32, radius, height float32) uint {
	return this.Detour.AddCapsuleObstacle(pos, radius, height)
}

func (this *DynamicScene) AddBoxObstacle(bmin, bmax []float32) uint {
	return this.Detour.AddBoxObstacle(bmin, bmax)
}

func (this *DynamicScene) AddBoxObstacle2(center, halfExtents []float32, yRadians float32) uint {
	return this.Detour.AddBoxObstacle2(center, halfExtents, yRadians)
}

func (this *DynamicScene) RemoveObstacle(obstacleId uint) {
	this.Detour.RemoveObstacle(obstacleId)
}
