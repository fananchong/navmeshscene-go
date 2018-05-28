#!/bin/bash

export CURDIR=$PWD
export GOPATH=$CURDIR/../../../../

cd $CURDIR/benchmarks
go test -v -test.bench=".*" -count=1


