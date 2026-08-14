#!/usr/bin/env bash

set -euo pipefail

APP=exnes
ARCH="${ARCH:-x86_64}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST="$ROOT/dist"
STAGE="$DIST/ExNES-windows-$ARCH"

rm -rf "$STAGE" "$DIST/ExNES-windows-$ARCH.zip"
mkdir -p "$STAGE"

CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -H windowsgui" -o "$STAGE/$APP.exe" ./cmd/sdl

cd "$STAGE"

for i in 1 2 3 4 5; do
    pre=$(find . -maxdepth 1 -type f -name '*.dll' | wc -l)
       for f in "$APP.exe" ./*.dll; do
        [ -e "$f" ] || continue
        ldd "$f" 2>/dev/null \
            | grep -iE '=> /(mingw64|ucrt64|clang64)/' \
            | awk '{print $3}' | sort -u \
            | while read -r dll; do
                  [ -f "$dll" ] && cp -n "$dll" . || true
              done
    done
    post=$(find . -maxdepth 1 -type f -name '*.dll' | wc -l)
    [ "$pre" = "$post" ] && break 

done 

cd "$DIST"
zip -qr "Win-$ARCH.zip" "$(basename "$STAGE")"

ls -lh "$DIST"/*.zip
