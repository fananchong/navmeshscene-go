package NavMeshScene

import (
	detour "github.com/fananchong/recastnavigation-go/Detour"
)

const (
	POLYAREA_GROUND int = 0
	POLYAREA_WATER  int = 1
	POLYAREA_ROAD   int = 2
	POLYAREA_DOOR   int = 3
	POLYAREA_GRASS  int = 4
	POLYAREA_JUMP   int = 5
)

const (
	POLYFLAGS_WALK     uint16 = 0x01   // Ability to walk (ground, grass, road)
	POLYFLAGS_SWIM     uint16 = 0x02   // Ability to swim (water).
	POLYFLAGS_DOOR     uint16 = 0x04   // Ability to move through doors.
	POLYFLAGS_JUMP     uint16 = 0x08   // Ability to jump.
	POLYFLAGS_DISABLED uint16 = 0x10   // Disabled polygon
	POLYFLAGS_ALL      uint16 = 0xffff // All abilities.
)

const DEFAULT_AREA_COST_GROUND float32 = 1.0
const DEFAULT_AREA_COST_WATER float32 = 10.0
const DEFAULT_AREA_COST_ROAD float32 = 1.0
const DEFAULT_AREA_COST_DOOR float32 = 1.0
const DEFAULT_AREA_COST_GRASS float32 = 2.0
const DEFAULT_AREA_COST_JUMP float32 = 1.5

const DEFAULT_INCLUDE_FLAGS uint16 = POLYFLAGS_ALL ^ POLYFLAGS_DISABLED
const DEFAULT_EXCLUDE_FLAGS uint16 = 0

type Filter struct {
	mFilter *detour.DtQueryFilter
}

func NewFilter() *Filter {
	this := &Filter{
		mFilter: detour.DtAllocDtQueryFilter(),
	}
	this.constructor()
	return this
}

func (this *Filter) constructor() {
	this.mFilter.SetAreaCost(POLYAREA_GROUND, DEFAULT_AREA_COST_GROUND)
	this.mFilter.SetAreaCost(POLYAREA_WATER, DEFAULT_AREA_COST_WATER)
	this.mFilter.SetAreaCost(POLYAREA_ROAD, DEFAULT_AREA_COST_ROAD)
	this.mFilter.SetAreaCost(POLYAREA_DOOR, DEFAULT_AREA_COST_DOOR)
	this.mFilter.SetAreaCost(POLYAREA_GRASS, DEFAULT_AREA_COST_GRASS)
	this.mFilter.SetAreaCost(POLYAREA_JUMP, DEFAULT_AREA_COST_JUMP)
	this.mFilter.SetIncludeFlags(DEFAULT_INCLUDE_FLAGS)
	this.mFilter.SetExcludeFlags(DEFAULT_EXCLUDE_FLAGS)
}

func (this *Filter) destructor() {
}

func (this *Filter) SetAreaCost(i int, cost float32) {
	this.mFilter.SetAreaCost(i, cost)
}

func (this *Filter) SetIncludeFlags(flags uint16) {
	this.mFilter.SetIncludeFlags(flags)
}

func (this *Filter) SetExcludeFlags(flags uint16) {
	this.mFilter.SetExcludeFlags(flags)
}
