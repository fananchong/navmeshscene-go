#!/bin/bash

export CURDIR=$PWD
export GOPATH=$CURDIR/../../../../

cd $CURDIR/benchmarks
go test scene_test.go -v -test.bench=".*" -count=1


