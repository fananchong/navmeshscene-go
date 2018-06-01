package NavMeshScene

import (
	detour "github.com/fananchong/recastnavigation-go/Detour"
	dtcache "github.com/fananchong/recastnavigation-go/DetourTileCache"
	"github.com/fananchong/recastnavigation-go/fastlz"
)

type FastLZCompressor struct{}

func (this *FastLZCompressor) MaxCompressedSize(bufferSize int32) int32 {
	return int32(float64(bufferSize) * 1.05)
}
func (this *FastLZCompressor) Compress(buffer []byte, bufferSize int32, compressed []byte, maxCompressedSize int32, compressedSize *int32) detour.DtStatus {
	*compressedSize = int32(fastlz.Fastlz_compress(buffer, int(bufferSize), compressed))
	return detour.DT_SUCCESS
}
func (this *FastLZCompressor) Decompress(compressed []byte, compressedSize int32, buffer []byte, maxBufferSize int32, bufferSize *int32) detour.DtStatus {
	*bufferSize = int32(fastlz.Fastlz_decompress(compressed, int(compressedSize), buffer, int(maxBufferSize)))
	if *bufferSize < 0 {
		return detour.DT_FAILURE
	} else {
		return detour.DT_SUCCESS
	}
}

type MeshProcess struct{}

func (this *MeshProcess) Process(params *detour.DtNavMeshCreateParams, polyAreas []uint8, polyFlags []uint16) {
	// Update poly flags from areas.
	for i := 0; i < int(params.PolyCount); i++ {
		if polyAreas[i] == dtcache.DT_TILECACHE_WALKABLE_AREA {
			polyAreas[i] = uint8(POLYAREA_GROUND)
		}
		if polyAreas[i] == uint8(POLYAREA_GROUND) ||
			polyAreas[i] == uint8(POLYAREA_GRASS) ||
			polyAreas[i] == uint8(POLYAREA_ROAD) {
			polyFlags[i] = POLYFLAGS_WALK
		} else if polyAreas[i] == uint8(POLYAREA_WATER) {
			polyFlags[i] = POLYFLAGS_SWIM
		} else if polyAreas[i] == uint8(POLYAREA_DOOR) {
			polyFlags[i] = POLYFLAGS_WALK | POLYFLAGS_DOOR
		}
	}

	// TODO: Pass in off-mesh connections.
}
