#!/bin/bash
set -eou pipefail

rm -rf build && mkdir build
OUT_DIR="$(realpath build)"
SRC_DIR="$(realpath .)"
FRONTEND_DIR="$(realpath frontend)"

echo "Building backend..."
cd "$OUT_DIR"
go build "$SRC_DIR"


echo "Building frontend..."
cd "$FRONTEND_DIR"
npm ci
npm run build
mv dist "$OUT_DIR/pb_public"

echo "Done!"
