#!/bin/bash
# Cross-compiles codesync for macOS (Intel/Apple Silicon), Windows, and Linux
# into ./bin. Run from the repo root: ./build.sh
set -euo pipefail

cd "$(dirname "$0")"
mkdir -p bin

targets=(
	"darwin amd64 codesync-darwin-amd64"
	"darwin arm64 codesync-darwin-arm64"
	"windows amd64 codesync-windows-amd64.exe"
	"linux amd64 codesync-linux-amd64"
)

for target in "${targets[@]}"; do
	read -r goos goarch out <<<"$target"
	echo "Building $out ($goos/$goarch)..."
	GOOS="$goos" GOARCH="$goarch" go build -o "bin/$out" .
done

echo "Done. Binaries are in ./bin"
