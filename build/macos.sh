#!/usr/bin/env bash

set -euo pipefail

APP=exnes
BUNDLE=ExNES.app
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST="$ROOT/dist"
APPDIR="$DIST/$BUNDLE"
ARCH="$(uname -m)"

rm -rf "$APPDIR" "$DIST/ExNES-$ARCH.dmg"
mkdir -p "$APPDIR/Contents/MacOS" "$APPDIR/Contents/Resources" "$APPDIR/Contents/Frameworks"

CGO_ENABLED=1 go build -trimpath -ldflags "-s -w" -o "$APPDIR/Contents/MacOS/$APP" ./cmd/sdl

cp "$ROOT/build/Info.plist" "$APPDIR/Contents/Info.plist"

ICON="$ROOT/build/exnes-1024.png"
ICONSET="$DIST/exnes.iconset"

rm -rf "$ICONSET"
mkdir -p "$ICONSET"
 
for s in 16 32 128 256 512; do
    sips -z $s          $s          "$ICON" --out "$ICONSET/icon_${s}x${s}.png"     >/dev/null
    sips -z $((s * 2))  $((s * 2))  "$ICON" --out "$ICONSET/icon_${s}x${s}@2x.png"  >/dev/null
done
iconutil -c icns "$ICONSET" -o "$APPDIR/Contents/Resources/$APP.icns"
rm -rf "$ICONSET"

dylibbundler --overwrite-files --bundle-deps --create-dir --fix-file "$APPDIR/Contents/MacOS/$APP" --dest-dir "$APPDIR/Contents/Frameworks/" --install-path "@executable_path/../Frameworks/"

codesign --force --deep --sign - "$APPDIR"
codesign --verify --verbose "$APPDIR"

hdiutil create -volname ExNES -srcfolder "$APPDIR" -ov -format UDZO    "$DIST/ExNES-$ARCH.dmg" >/dev/null

ls -lh "$DIST"

echo "done"


