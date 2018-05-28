#!/bin/bash

export CURDIR=$PWD
export GOPATH=$CURDIR/../../../../

go test -v -test.bench=".*" -count=1


