@echo off
setlocal enabledelayedexpansion

echo.

if not "%~1"=="" (
    rem 指定了模块，检查是否存在
    if not exist ".\cmd\%~1\" (
        echo Error: module "%~1" not found under .\cmd\
        exit /b 1
    )
    set "apps=%~1"
) else (
    rem 未指定，自动扫描 cmd 下的所有目录
    set "apps="
    for /d %%d in (.\cmd\*) do (
        set "apps=!apps! %%~nxd"
    )
)

for %%a in (%apps%) do (
    echo [%%a] Running wire...
    wire .\cmd\%%a
    if !errorlevel! neq 0 (
        echo [%%a] wire failed! Aborting.
        exit /b 1
    )

    echo [%%a] Running go build...
    go build -o %%a.exe .\cmd\%%a
    if !errorlevel! neq 0 (
        echo [%%a] go build failed! Aborting.
        exit /b 1
    )

    echo [%%a] Done.
    echo.
)

echo All builds completed successfully!
endlocal
