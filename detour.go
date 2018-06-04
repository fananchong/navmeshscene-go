package NavMeshScene

import (
	"io/ioutil"
	"reflect"
	"strings"
	"sync"
	"unsafe"

	detour "github.com/fananchong/recastnavigation-go/Detour"
	dtcache "github.com/fananchong/recastnavigation-go/DetourTileCache"
)

type Detour struct {
	mbStaticMesh bool
	mMaxNode     int
	mMesh        *detour.DtNavMesh
	mTileCache   *dtcache.DtTileCache
	mQuery       *detour.DtNavMeshQuery

	mBoundsMin [3]float32
	mBoundsMax [3]float32

	mHeightMode          int
	mTcomp               FastLZCompressor
	mTmproc              MeshProcess
	mQueryForHeightMode2 *detour.DtNavMeshQuery
}

func NewDetour(bStaticMesh bool, maxNode uint16) *Detour {
	detour.DtAssert(maxNode != 0)
	this := &Detour{
		mbStaticMesh: bStaticMesh,
		mMaxNode:     int(maxNode),
	}
	return this
}

var mStaticMesh map[string]*detour.DtNavMesh
var mStaticMeshMutex sync.Mutex

const FILE_SUFFIX_0 string = ".tile.bin"
const FILE_SUFFIX_1 string = ".tilecache.bin"

func (this *Detour) Load(path string) int {
	detour.DtAssert(strings.HasSuffix(path, FILE_SUFFIX_0) || strings.HasSuffix(path, FILE_SUFFIX_1))
	errCode := 0
	var mesh *detour.DtNavMesh
	if this.mbStaticMesh {
		mesh = this.createStaticMesh(path, &errCode)
	} else {
		var tempMesh *detour.DtNavMesh
		if this.mHeightMode == HEIGHT_MODE_2 {
			tempPath := strings.Replace(path, FILE_SUFFIX_1, FILE_SUFFIX_0, 1)
			tempMesh = this.createStaticMesh(tempPath, &errCode)
			if errCode != 0 {
				return errCode
			}
			this.mQueryForHeightMode2 = detour.DtAllocNavMeshQuery()
			if this.mQueryForHeightMode2 == nil {
				return 1
			}

			status := this.mQueryForHeightMode2.Init(tempMesh, this.mMaxNode)
			if !detour.DtStatusSucceed(status) {
				return 2
			}
		}
		tempPath := strings.Replace(path, FILE_SUFFIX_0, FILE_SUFFIX_1, 1)
		mesh = this.loadDynamicMesh(tempPath, &errCode)
	}

	if errCode != 0 {
		return errCode
	}

	this.mQuery = detour.DtAllocNavMeshQuery()
	if this.mQuery == nil {
		return 3
	}

	status := this.mQuery.Init(mesh, this.mMaxNode)
	if !detour.DtStatusSucceed(status) {
		return 4
	}
	return 0
}

func (this *Detour) createStaticMesh(path string, errCode *int) *detour.DtNavMesh {
	mStaticMeshMutex.Lock()
	defer mStaticMeshMutex.Unlock()
	if m, ok := mStaticMesh[path]; ok {
		return m
	} else {
		mesh := this.loadStaticMesh(path, errCode)
		if *errCode == 0 {
			mStaticMesh[path] = mesh
		}
		return mesh
	}
}

type NavMeshSetHeader struct {
	magic      int32
	version    int32
	numTiles   int32
	params     detour.DtNavMeshParams
	boundsMinX float32
	boundsMinY float32
	boundsMinZ float32
	boundsMaxX float32
	boundsMaxY float32
	boundsMaxZ float32
}

type NavMeshTileHeader struct {
	tileRef  detour.DtTileRef
	dataSize int32
}

type TileCacheSetHeader struct {
	magic       int32
	version     int32
	numTiles    int32
	meshParams  detour.DtNavMeshParams
	cacheParams dtcache.DtTileCacheParams
	boundsMinX  float32
	boundsMinY  float32
	boundsMinZ  float32
	boundsMaxX  float32
	boundsMaxY  float32
	boundsMaxZ  float32
}

type TileCacheTileHeader struct {
	tileRef  dtcache.DtCompressedTileRef
	dataSize int32
}

const NAVMESHSET_MAGIC int32 = int32('M')<<24 | int32('S')<<16 | int32('E')<<8 | int32('T')
const NAVMESHSET_VERSION int32 = 1
const TILECACHESET_MAGIC int32 = int32('T')<<24 | int32('S')<<16 | int32('E')<<8 | int32('T')
const TILECACHESET_VERSION int32 = 1

func (this *Detour) loadStaticMesh(path string, errCode *int) *detour.DtNavMesh {
	*errCode = 0
	meshData, err := ioutil.ReadFile(path)
	if err != nil {
		*errCode = 101
		return nil
	}

	// Read header.
	header := (*NavMeshSetHeader)(unsafe.Pointer(&(meshData[0])))
	if header.magic != NAVMESHSET_MAGIC {
		*errCode = 103
		return nil
	}
	if header.version != NAVMESHSET_VERSION {
		*errCode = 104
		return nil
	}

	this.mBoundsMin[0] = header.boundsMinX
	this.mBoundsMin[1] = header.boundsMinY
	this.mBoundsMin[2] = header.boundsMinZ
	this.mBoundsMax[0] = header.boundsMaxX
	this.mBoundsMax[1] = header.boundsMaxY
	this.mBoundsMax[2] = header.boundsMaxZ

	mesh := detour.DtAllocNavMesh()
	if mesh == nil {
		*errCode = 105
		return nil
	}
	state := mesh.Init(&header.params)
	if detour.DtStatusFailed(state) {
		*errCode = 106
		return nil
	}

	// Read tiles.
	d := int32(unsafe.Sizeof(*header))
	for i := 0; i < int(header.numTiles); i++ {
		tileHeader := (*NavMeshTileHeader)(unsafe.Pointer(&(meshData[d])))
		if tileHeader.tileRef == 0 || tileHeader.dataSize == 0 {
			break
		}
		d += int32(unsafe.Sizeof(*tileHeader))
		data := meshData[d : d+tileHeader.dataSize]
		state = mesh.AddTile(data, int(tileHeader.dataSize), detour.DT_TILE_FREE_DATA, tileHeader.tileRef, nil)
		if detour.DtStatusFailed(state) {
			*errCode = 108
			return nil
		}
		d += tileHeader.dataSize
	}
	return mesh
}

func (this *Detour) loadDynamicMesh(path string, errCode *int) *detour.DtNavMesh {
	*errCode = 0
	meshData, err := ioutil.ReadFile(path)
	if err != nil {
		*errCode = 201
		return nil
	}

	// Read header.
	header := (*TileCacheSetHeader)(unsafe.Pointer(&(meshData[0])))
	if header.magic != TILECACHESET_MAGIC {
		*errCode = 203
		return nil
	}
	if header.version != TILECACHESET_VERSION {
		*errCode = 204
		return nil
	}

	this.mBoundsMin[0] = header.boundsMinX
	this.mBoundsMin[1] = header.boundsMinY
	this.mBoundsMin[2] = header.boundsMinZ
	this.mBoundsMax[0] = header.boundsMaxX
	this.mBoundsMax[1] = header.boundsMaxY
	this.mBoundsMax[2] = header.boundsMaxZ

	defer func() {
		if *errCode != 0 {
			this.mMesh = nil
			this.mTileCache = nil
		}
	}()

	this.mMesh = detour.DtAllocNavMesh()
	if this.mMesh == nil {
		*errCode = 205
		return nil
	}
	status := this.mMesh.Init(&header.meshParams)
	if detour.DtStatusFailed(status) {
		*errCode = 206
		return nil
	}

	this.mTileCache = dtcache.DtAllocTileCache()
	if this.mTileCache == nil {
		*errCode = 207
		return nil
	}

	status = this.mTileCache.Init(&header.cacheParams, &this.mTcomp, &this.mTmproc)
	if detour.DtStatusFailed(status) {
		*errCode = 208
		return nil
	}

	// Read tiles.
	d := int(unsafe.Sizeof(*header))
	for i := 0; i < int(header.numTiles); i++ {
		tileHeader := (*TileCacheTileHeader)(unsafe.Pointer(&(meshData[d])))
		d += int(unsafe.Sizeof(*tileHeader))
		if tileHeader.tileRef == 0 || tileHeader.dataSize == 0 {
			break
		}

		var tempData []byte
		sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&tempData)))
		sliceHeader.Cap = int(tileHeader.dataSize)
		sliceHeader.Len = int(tileHeader.dataSize)
		sliceHeader.Data = uintptr(unsafe.Pointer(&meshData[d]))
		d += int(tileHeader.dataSize)
		data := make([]byte, tileHeader.dataSize)
		copy(data, tempData)

		var tile dtcache.DtCompressedTileRef
		status = this.mTileCache.AddTile(data, tileHeader.dataSize, dtcache.DT_COMPRESSEDTILE_FREE_DATA, &tile)
		if detour.DtStatusFailed(status) {
			*errCode = 211
			return nil
		}

		if tile != 0 {
			this.mTileCache.BuildNavMeshTile(tile, this.mMesh)
		} else {
			*errCode = 212
			return nil
		}
	}
	return this.mMesh
}

func (this *Detour) TryMove(
	startPolyRef detour.DtPolyRef,
	startPos []float32,
	endPos []float32,
	halfExtents []float32,
	filter *detour.DtQueryFilter,
	realEndPolyRef *detour.DtPolyRef,
	realEndPos []float32,
	bHit *bool) bool {
	*bHit = false
	if this.mQuery == nil {
		return false
	}
	const visitedNodeCount = 16
	var visited [visitedNodeCount]detour.DtPolyRef
	nvisited := 0
	status := this.mQuery.MoveAlongSurface(
		startPolyRef,
		startPos,
		endPos,
		filter,
		realEndPos,
		visited[:],
		&nvisited,
		visitedNodeCount,
		bHit)

	if detour.DtStatusDetail(status, detour.DT_INVALID_PARAM) {
		var tempRef detour.DtPolyRef
		var tempPos [3]float32
		this.mQuery.FindNearestPoly(startPos, halfExtents, filter, &tempRef, tempPos[:])
		startPolyRef = tempRef
		detour.DtVcopy(startPos, tempPos[:])

		status = this.mQuery.MoveAlongSurface(
			startPolyRef,
			startPos,
			endPos,
			filter,
			realEndPos,
			visited[:],
			&nvisited,
			visitedNodeCount,
			bHit)
	}

	if !detour.DtStatusSucceed(status) {
		return false
	}

	*realEndPolyRef = visited[nvisited-1]

	if this.mHeightMode != HEIGHT_MODE_2 {
		var h float32 = 0
		this.mQuery.GetPolyHeight(*realEndPolyRef, realEndPos, &h)
		realEndPos[1] = h
	} else {
		var tempRef detour.DtPolyRef
		var tempPos [3]float32
		this.mQueryForHeightMode2.FindNearestPoly(realEndPos, halfExtents, filter, &tempRef, tempPos[:])
		realEndPos[1] = tempPos[1]
	}
	return true
}

func (this *Detour) GetPoly(
	pos []float32,
	halfExtents []float32,
	filter *detour.DtQueryFilter,
	nearestRef *detour.DtPolyRef,
	nearestPt []float32) bool {
	if this.mQuery == nil {
		return false
	}
	var nRef detour.DtPolyRef
	status := this.mQuery.FindNearestPoly(pos, halfExtents, filter, &nRef, nearestPt)
	if !detour.DtStatusSucceed(status) {
		return false
	}
	*nearestRef = nRef
	return true
}

func (this *Detour) Raycast(
	startPolyRef detour.DtPolyRef,
	startPos []float32,
	endPos []float32,
	filter *detour.DtQueryFilter,
	bHit *bool,
	hitPos []float32) bool {
	if this.mQuery == nil {
		return false
	}
	var t float32 = 0
	var hitNormal [3]float32
	const polysCount = 16
	var polys [polysCount]detour.DtPolyRef
	npolys := 0
	status := this.mQuery.Raycast(startPolyRef, startPos, endPos, filter,
		&t, hitNormal[:], polys[:], &npolys, polysCount)
	if !detour.DtStatusSucceed(status) {
		return false
	}
	*bHit = (t <= 1)
	if *bHit {
		detour.DtVlerp(hitPos, startPos, endPos, t)
		if npolys > 0 {
			var h float32 = 0
			this.mQuery.GetPolyHeight(polys[npolys-1], hitPos, &h)
			hitPos[1] = h
		}
	}
	return true
}

func (this *Detour) RandomPosition(halfExtents []float32,
	filter *detour.DtQueryFilter,
	frand func() float32,
	randomRef *detour.DtPolyRef,
	randomPt []float32) bool {
	if this.mQuery == nil {
		return false
	}
	var ref detour.DtPolyRef
	status := this.mQuery.FindRandomPoint(filter, frand, &ref, randomPt)
	if !detour.DtStatusSucceed(status) {
		return false
	}
	*randomRef = ref
	if this.mHeightMode == HEIGHT_MODE_2 {
		var tempRef detour.DtPolyRef
		var tempPos [3]float32
		this.mQueryForHeightMode2.FindNearestPoly(randomPt, halfExtents, filter, &tempRef, tempPos[:])
		randomPt[1] = tempPos[1]
	}
	return true
}

func (this *Detour) AddCapsuleObstacle(pos []float32, radius, height float32) uint {
	if this.mTileCache == nil {
		return 0
	}
	var obstacleId dtcache.DtObstacleRef
	status := this.mTileCache.AddObstacle(pos, radius, height, &obstacleId)
	if !detour.DtStatusSucceed(status) {
		return 0
	}
	return (uint)(obstacleId)
}

func (this *Detour) AddBoxObstacle(bmin []float32, bmax []float32) uint {
	if this.mTileCache == nil {
		return 0
	}
	var obstacleId dtcache.DtObstacleRef
	status := this.mTileCache.AddBoxObstacle(bmin, bmax, &obstacleId)
	if !detour.DtStatusSucceed(status) {
		return 0
	}
	return (uint)(obstacleId)
}

func (this *Detour) AddBoxObstacle2(center []float32, halfExtents []float32, yRadians float32) uint {
	if this.mTileCache == nil {
		return 0
	}
	var obstacleId dtcache.DtObstacleRef
	status := this.mTileCache.AddBoxObstacle2(center, halfExtents, yRadians, &obstacleId)
	if !detour.DtStatusSucceed(status) {
		return 0
	}
	return (uint)(obstacleId)
}

func (this *Detour) RemoveObstacle(obstacleId uint) {
	if this.mTileCache == nil || obstacleId == 0 {
		return
	}
	this.mTileCache.RemoveObstacle((dtcache.DtObstacleRef)(obstacleId))
}
