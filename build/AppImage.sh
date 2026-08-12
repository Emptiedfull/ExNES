#!/usr/bin/env bash

set -euo pipefail

APP=exnes
ARCH="${ARCH:-x86_64}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APPDIR="$ROOT/build/${APP}.AppDir"
DIST="$ROOT/dist"

echo "AHHH OK"

rm -rf "$APPDIR" "$DIST"
mkdir -p "$APPDIR/usr/bin" "$APPDIR/usr/share/applications" "$APPDIR/usr/share/icons/hicolor/256x256/apps"  "$DIST"

echo "setup binary"

CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o "$APPDIR/usr/bin/$APP" ./cmd/sdl

cp "$ROOT/build/${APP}.desktop" "$APPDIR/usr/share/applications/"
cp "$ROOT/build/${APP}.png" "$APPDIR/usr/share/icons/hicolor/256x256/apps"

export APPIMAGE_EXTRACT_AND_RUN=1
export OUTPUT="ExNES.AppImage"

linuxdeploy --appdir "$APPDIR" \
            --desktop-file "$APPDIR/usr/share/applications/${APP}.desktop" \
            --icon-file "$APPDIR/usr/share/icons/hicolor/256x256/apps/${APP}.png" \
            --plugin gtk \
            --output appimage

mv ./*.AppImage "$DIST/"
chmod +x "$DIST"/*.AppImage

