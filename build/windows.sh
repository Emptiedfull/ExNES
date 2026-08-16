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

echo "go build complete"

echo " bundling DLLs"
cd "$STAGE"
 
for pass in 1 2 3 4 5; do
    before=$(find . -maxdepth 1 -type f -name '*.dll' | wc -l)
 
      for f in "$APP.exe" ./*.dll; do
        [ -e "$f" ] || continue
        
        ldd "$f" 2>/dev/null | awk '{print $1}' | while read -r name; do
            case "$name" in *.dll|*.DLL) ;; *) continue ;; esac
            src="/mingw64/bin/$name"
            [ -f "$src" ] && cp -n "$src" . || true
        done || true
    done
    after=$(find . -maxdepth 1 -type f -name '*.dll' | wc -l)
    echo "    pass $pass: $before -> $after DLLs"
    [ "$before" = "$after" ] && break
done
 
echo "    final contents:"
ls -1 | sed 's/^/      /'

cd "$DIST"
zip -qr "Win-$ARCH.zip" "$(basename "$STAGE")"

ls -lh "$DIST"/*.zip
