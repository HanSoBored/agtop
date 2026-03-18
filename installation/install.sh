#!/bin/bash
set -e

# --- CONFIGURATION ---
REPO_OWNER="HanSoBored"
REPO_NAME="agtop"
BINARY_BASE_NAME="agtop"
FINAL_NAME="agtop"
INSTALL_DIR="/usr/local/bin"

# --- DETECT SYSTEM ---
OS="$(uname -s)"
ARCH="$(uname -m)"

echo "🔍 Detecting system..."
echo "   OS: $OS"
echo "   Arch: $ARCH"

SUFFIX=""
PLATFORM=""

# 1. DETECT OS & MAP ARCHITECTURE
if [ "$OS" = "Linux" ]; then
    # Check if running on Android (has Android-specific properties)
    if [ -f "/system/bin/app_process" ] || grep -q "Android" /proc/version 2>/dev/null; then
        PLATFORM="android"
        if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
            SUFFIX="android-aarch64"
        elif [[ "$ARCH" == armv7* ]] || [ "$ARCH" = "arm" ]; then
            SUFFIX="android-armv7"
        else
            echo "❌ Unsupported Architecture: $ARCH on Android"
            exit 1
        fi
    else
        PLATFORM="linux"
        if [ "$ARCH" = "x86_64" ]; then
            SUFFIX="linux-x86_64"
        elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
            SUFFIX="linux-aarch64"
        elif [[ "$ARCH" == armv7* ]] || [ "$ARCH" = "arm" ]; then
            SUFFIX="linux-armv7"
        else
            echo "❌ Unsupported Architecture: $ARCH on Linux"
            exit 1
        fi
    fi
elif [ "$OS" = "Darwin" ]; then
    PLATFORM="darwin"
    if [ "$ARCH" = "x86_64" ]; then
        SUFFIX="darwin-x86_64"
    elif [ "$ARCH" = "arm64" ]; then
        SUFFIX="darwin-aarch64"
    else
        echo "❌ Unsupported Architecture: $ARCH on macOS"
        exit 1
    fi
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

TARGET_FILE="${BINARY_BASE_NAME}-${SUFFIX}"
echo "🎯 Target Release Asset: $TARGET_FILE"

# --- DOWNLOADING ---
echo "⬇️  Downloading latest release..."
DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/latest/download/$TARGET_FILE"

# Use curl to download to temp folder
# -L follows redirects
# -f fails silently on server error (404) so we can catch it
if ! curl -f -L -o "/tmp/$BINARY_BASE_NAME" "$DOWNLOAD_URL"; then
    echo "❌ Error: Failed to download. The release asset '$TARGET_FILE' might not exist yet."
    exit 1
fi

# --- INSTALLING ---
echo "📦 Installing to $INSTALL_DIR..."
chmod +x "/tmp/$BINARY_BASE_NAME"

# Check write permissions
if [ -w "$INSTALL_DIR" ]; then
    mv "/tmp/$BINARY_BASE_NAME" "$INSTALL_DIR/$FINAL_NAME"
else
    echo "🔑 Sudo permission required to move binary to $INSTALL_DIR"
    sudo mv "/tmp/$BINARY_BASE_NAME" "$INSTALL_DIR/$FINAL_NAME"
fi

echo "✅ Installed successfully!"
echo "   Binary location: $INSTALL_DIR/$FINAL_NAME"
echo "   You can now run it using: $FINAL_NAME"
