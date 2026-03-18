#!/bin/bash
set -e

BINARY_NAME="agtop"
BUILD_DIR="../../build"
INSTALL_DIR="/usr/local/bin"

echo "🐹 agtop Build Script"
echo "===================="
echo ""
echo "Select build target:"
echo "1. Linux"
echo "2. Android Static"
echo ""
read -p "Enter choice [1-2]: " choice

cd cmd/agtop/
mkdir -p "$BUILD_DIR"

case $choice in
    1)
        echo "🐹 Building for Linux (Release Mode)..."
        go build -ldflags="-s -w" -o "$BUILD_DIR/$BINARY_NAME"
        ;;
    2)
        echo "🐹 Building for Android Static (Release Mode)..."
        CGO_ENABLED=0 GOOS=android GOARCH=arm64 go build -ldflags="-s -w" -o "$BUILD_DIR/${BINARY_NAME}_android_static"
        ;;
    *)
        echo "❌ Invalid choice"
        exit 1
        ;;
esac

echo ""
echo "📦 Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    if [ "$choice" = "2" ]; then
        cp "$BUILD_DIR/${BINARY_NAME}_android_static" "$INSTALL_DIR/$BINARY_NAME"
    else
        cp "$BUILD_DIR/$BINARY_NAME" "$INSTALL_DIR/"
    fi
else
    if [ "$choice" = "2" ]; then
        sudo cp "$BUILD_DIR/${BINARY_NAME}_android_static" "$INSTALL_DIR/$BINARY_NAME"
    else
        sudo cp "$BUILD_DIR/$BINARY_NAME" "$INSTALL_DIR/"
    fi
fi

echo "✅ Success! Installed"
