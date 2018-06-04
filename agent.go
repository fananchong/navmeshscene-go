package NavMeshScene

import (
	"math"
	"math/rand"

	"github.com/fananchong/navmeshscene-go/aoi"
	detour "github.com/fananchong/recastnavigation-go/Detour"
)

var DEFAULT_HALF_EXTENTS = [3]float32{0.6, 2.0, 0.6}
var ZERO = [3]float32{0, 0, 0}
var EPSILON = 0.0001

type Agent struct {
	aoi.Object
	Id          uint64
	HalfExtents [3]float32
	Position    [3]float32
	Velocity    [3]float32
	CurPolyRef  detour.DtPolyRef
	Filter      *detour.DtQueryFilter
	Scene       *Scene
	OnHit       func(uint64)
}

func NewAgent() *Agent {
	this := &Agent{}
	this.HalfExtents = DEFAULT_HALF_EXTENTS
	this.Filter = detour.DtAllocDtQueryFilter()
	return this
}

func (this *Agent) Update(delta float32) {
	if this.Velocity[0] == 0 && this.Velocity[1] == 0 && this.Velocity[2] == 0 {
		return
	}
	endPos := [3]float32{
		this.Position[0] + this.Velocity[0]*delta,
		this.Position[1] + this.Velocity[1]*delta,
		this.Position[2] + this.Velocity[2]*delta}
	if agent := this.checkPosByAOI(this.Position[0], this.Position[2], &endPos[0], &endPos[2], true); agent != 0 {
		detour.DtVcopy(this.Velocity[:], ZERO[:])
		if this.OnHit != nil {
			this.OnHit(agent)
		}
		return
	}
	var bHit bool
	var realEndPolyRef detour.DtPolyRef
	var realEndPos [3]float32
	if !this.TryMove(endPos[:], &realEndPolyRef, realEndPos[:], &bHit) {
		return
	}
	if bHit {
		detour.DtVcopy(this.Velocity[:], ZERO[:])
		if this.OnHit != nil {
			this.OnHit(0)
		}
		return
	}
	this.CurPolyRef = realEndPolyRef
	detour.DtVcopy(this.Position[:], realEndPos[:])
	if math.Abs(float64(this.X-this.Position[0])) >= EPSILON || math.Abs(float64(this.Y-this.Position[2])) >= EPSILON {
		this.X = this.Position[0]
		this.Y = this.Position[2]
		this.Scene.Update(&this.Object)
	}
}

func (this *Agent) TryMove(endPos []float32, realEndPolyRef *detour.DtPolyRef, realEndPos []float32, bHit *bool) bool {
	if this.Scene != nil {
		return this.Scene.Detour.TryMove(
			this.CurPolyRef,
			this.Position[:],
			endPos,
			this.HalfExtents[:],
			this.Filter,
			realEndPolyRef,
			realEndPos,
			bHit)
	}
	return false
}

func (this *Agent) SetPosition(v []float32) bool {
	if this.Scene != nil {
		if this.checkPosByAOI(this.Position[0], this.Position[2], &this.Position[0], &this.Position[2], false) == 0 {
			this.Scene.Detour.GetPoly(v, this.HalfExtents[:], this.Filter, &this.CurPolyRef, this.Position[:])
			this.X = this.Position[0]
			this.Y = this.Position[2]
			this.Scene.Update(&this.Object)
			return true
		}
	}
	return false
}

func (this *Agent) RandomPosition() {
	if this.Scene != nil {
	LABLE_RANDOM:
		this.Scene.Detour.RandomPosition(this.HalfExtents[:], this.Filter, rand.Float32, &this.CurPolyRef, this.Position[:])
		if this.checkPosByAOI(this.Position[0], this.Position[2], &this.Position[0], &this.Position[2], false) != 0 {
			goto LABLE_RANDOM
		}
		this.X = this.Position[0]
		this.Y = this.Position[2]
		this.Scene.Update(&this.Object)
	}
}

func (this *Agent) Raycast(endPos []float32, bHit *bool, hitPos []float32) bool {
	if this.Scene != nil {
		return this.Scene.Detour.Raycast(
			this.CurPolyRef,
			this.Position[:],
			endPos,
			this.Filter,
			bHit,
			hitPos)
	}
	return false
}

func (this *Agent) checkPosByAOI(srcX, srcY float32, dstX, dstY *float32, bMove bool) uint64 {
	rect := aoi.NewRect(*dstX-this.HalfExtents[0], *dstX+this.HalfExtents[0], *dstY-this.HalfExtents[2], *dstY+this.HalfExtents[2])
	agents := this.Scene.QuadTree.Query4(&this.Object, rect)
	for a := agents; a != nil; a = a.QueryNext {
		if a == &this.Object {
			continue
		}
		tempRect := aoi.NewRect(this.Position[0]-this.HalfExtents[0],
			this.Position[0]+this.HalfExtents[0],
			this.Position[2]-this.HalfExtents[2],
			this.Position[2]+this.HalfExtents[2])
		if rect.Intersects(tempRect) {
			return a.UserData
		}
	}
	return 0
}
