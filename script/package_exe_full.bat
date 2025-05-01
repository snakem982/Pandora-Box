@echo off
setlocal

set APP_NAME=PandoraBox
set ICON_FILE=icon.ico
set OUTPUT_DIR=build
set MAIN_GO_FILE=main.go

:: 检查 windres
where windres >nul 2>nul
if errorlevel 1 (
    echo ❌ 未找到 windres，请安装 MinGW 或准备一个 windres.exe
    pause
    exit /b
)

:: 生成 resource.rc
echo 1 ICON "%ICON_FILE%" > resource.rc

:: 编译 resource.syso
echo [1/3] 生成资源文件...
windres resource.rc -O coff -o resource.syso

:: 创建输出目录
if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

:: 编译 Go 程序
echo [2/3] 编译 Go 程序...
go build -ldflags="-H windowsgui" -o %OUTPUT_DIR%\%APP_NAME%.exe %MAIN_GO_FILE%

:: 清理
del resource.rc
del resource.syso

echo.
echo ✅ 完成！EXE 在 %OUTPUT_DIR%\%APP_NAME%.exe
pause
