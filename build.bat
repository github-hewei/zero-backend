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

rem 需要生成 Swagger 文档的模块
set "swag_modules=admin api"

for %%a in (%apps%) do (
    for %%s in (%swag_modules%) do (
        if "%%a"=="%%s" (
            echo [%%a] Running swag init...
            swag init -g .\cmd\%%a\main.go -o .\docs\%%a --pd
            if !errorlevel! neq 0 (
                echo [%%a] swag init failed! Aborting.
                exit /b 1
            )
        )
    )

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
