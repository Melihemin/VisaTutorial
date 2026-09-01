#!/bin/bash
# Build script - Derleme betiği

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC_DIR="$PROJECT_ROOT/src"
BUILD_DIR="$PROJECT_ROOT/build"
CONFIG_DIR="$PROJECT_ROOT/config"

echo "🔨 Randevu Tracker Derleniyor..."
echo "================================"
echo "Kaynak: $SRC_DIR"
echo "Build: $BUILD_DIR"
echo ""

# Ensure build directory exists
mkdir -p "$BUILD_DIR"

# Go to source directory and build all Go files
cd "$SRC_DIR" || exit 1

echo "📦 Tüm Go dosyaları derleniyor..."
go build -o "$BUILD_DIR/randevu_tracker" ./*.go 2>&1

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Derleme başarılı!"
    echo "📍 Binary konumu: $BUILD_DIR/randevu_tracker"
    ls -lh "$BUILD_DIR/randevu_tracker"
    echo ""
    echo "💡 Başlatmak için:"
    echo "   $BUILD_DIR/randevu_tracker"
else
    echo ""
    echo "❌ Derleme başarısız!"
    exit 1
fi
