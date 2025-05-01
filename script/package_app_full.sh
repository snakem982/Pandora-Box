#!/bin/bash

set -e

# --- 基本配置 ---
APP_NAME="MyGoApp"
APP_EXECUTABLE="my_go_program"   # Go输出的二进制文件名
MAIN_GO_FILE="main.go"            # 主Go文件
ICON_FILE="AppIcon.icns"          # 应用图标文件 (.icns)
BUNDLE_ID="com.example.$APP_NAME"
VERSION="1.0"
BUILD_DIR="./build"
APP_DIR="$BUILD_DIR/$APP_NAME.app"
DMG_NAME="$APP_NAME.dmg"

# --- 清理旧文件 ---
echo "Cleaning previous builds..."
rm -rf "$BUILD_DIR" "$DMG_NAME"
mkdir -p "$BUILD_DIR"

# --- 编译Go程序 ---
echo "Building Go executable..."
GOOS=darwin GOARCH=amd64 go build -o "$APP_EXECUTABLE" "$MAIN_GO_FILE"
# 如果是Apple Silicon Mac，改成：GOARCH=arm64

# --- 创建.app包结构 ---
echo "Creating .app bundle..."
mkdir -p "$APP_DIR/Contents/MacOS"
mkdir -p "$APP_DIR/Contents/Resources"

# --- 创建Info.plist ---
cat > "$APP_DIR/Contents/Info.plist" <<EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundleDisplayName</key>
    <string>$APP_NAME</string>
    <key>CFBundleIdentifier</key>
    <string>$BUNDLE_ID</string>
    <key>CFBundleVersion</key>
    <string>$VERSION</string>
    <key>CFBundleShortVersionString</key>
    <string>$VERSION</string>
    <key>CFBundleExecutable</key>
    <string>$APP_EXECUTABLE</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleIconFile</key>
    <string>$(basename "$ICON_FILE")</string>
</dict>
</plist>
EOL

# --- 拷贝可执行文件和图标 ---
mv "$APP_EXECUTABLE" "$APP_DIR/Contents/MacOS/"
if [ -f "$ICON_FILE" ]; then
    echo "Copying icon..."
    cp "$ICON_FILE" "$APP_DIR/Contents/Resources/"
else
    echo "Warning: Icon file '$ICON_FILE' not found. Skipping icon setup."
fi

# --- 创建 DMG 安装包 ---
echo "Creating DMG..."
hdiutil create -volname "$APP_NAME" -srcfolder "$APP_DIR" -ov -format UDZO "$DMG_NAME"

echo "✅ Build complete! DMG generated: $DMG_NAME"
