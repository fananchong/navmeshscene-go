package benchmarks

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fananchong/navmeshscene-go/aoi"
	"github.com/fananchong/navmeshscene-go/tests"
)

const w float32 = 1000
const h float32 = 1000
const r float32 = 0.6

const PLAYER_COUNT int = 5000

var items []*tests.A
var scn *aoi.QuadTree = aoi.NewDefaultSecene(aoi.NewRect(0, w, 0, h))

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < PLAYER_COUNT; i++ {
		tests.Add(scn, &items)
	}
}

func Benchmark_Add(t *testing.B) {
	var items1 []*tests.A
	var scn1 *aoi.QuadTree = aoi.NewDefaultSecene(aoi.NewRect(0, w, 0, h))

	t.N = 5000
	for i := 0; i < t.N; i++ {
		tests.Add(scn1, &items1)
	}
}

func Benchmark_Query1(t *testing.B) {
	for i := 0; i < t.N; i++ {
		for j := 0; j < PLAYER_COUNT; j++ {
			scn.Query2(&items[j].Object, r)
		}
	}
}

func Benchmark_Query2(t *testing.B) {
	for i := 0; i < t.N; i++ {
		for j := 0; j < PLAYER_COUNT; j++ {
			var rect aoi.Rect
			rect.Init(items[j].X-r, items[j].X+r, items[j].Y-r, items[j].Y+r)
			scn.Query1(&rect)
		}
	}
}
