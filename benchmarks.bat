set CURDIR=%~dp0
set BASEDIR=%CURDIR%\..\..\..\..\
set GOPATH=%BASEDIR%
echo %GOPATH%

cd %CURDIR%\benchmarks
call go test scene_test.go -v -test.bench=".*" -count=1

cd %CURDIR%