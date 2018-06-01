package nmscene

import (
	detour "github.com/fananchong/recastnavigation-go/Detour"
	dtcache "github.com/fananchong/recastnavigation-go/DetourTileCache"
)

type Detour struct {
	mbStaticMesh bool
	mMaxNode     int
	mMesh        *detour.DtNavMesh
	mTileCache   *dtcache.DtTileCache
	mQuery       *detour.DtNavMeshQuery
}
