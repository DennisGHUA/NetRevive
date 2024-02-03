@echo off
echo [1/2] Updating packages...
go get -u ./...
echo [1/2] Packages updated and ready to compile.

echo.
echo [2/2] Compiling...
go build

if exist Build (
    move /y *.exe Build > nul
) else (
    mkdir Build
    move /y *.exe Build > nul
)

echo [2/2] Compilation done.
echo.

echo Press any key to exit...
pause > nul
