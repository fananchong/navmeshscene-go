package benchmarks

import "testing"

type ITestObj interface {
	GetPropX() int
	SetPropX(x int)
}

type Obj struct {
	X int
}

func (this *Obj) GetPropX() int {
	return this.X
}

func (this *Obj) SetPropX(x int) {
	this.X = x
}

func testFunc1(obj *Obj) {
	obj.X = obj.X + 1
}

func testFunc2(obj *Obj) {
	obj.SetPropX(obj.GetPropX() + 1)
}

func testFunc3(obj ITestObj) {
	obj.SetPropX(obj.GetPropX() + 1)
}

var myobj1 Obj

func Benchmark_Ptr(t *testing.B) {
	for i := 0; i < t.N; i++ {
		testFunc1(&myobj1)
	}
}

func Benchmark_PtrAndCallF(t *testing.B) {
	for i := 0; i < t.N; i++ {
		testFunc2(&myobj1)
	}
}

func Benchmark_Interface(t *testing.B) {
	for i := 0; i < t.N; i++ {
		testFunc3(&myobj1)
	}
}
