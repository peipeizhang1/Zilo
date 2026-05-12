#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
API_FILE="$ROOT_DIR/api/workflow.api"
OUT_DIR="$ROOT_DIR/api"

if command -v goctl >/dev/null 2>&1; then
  GOCTL_BIN="goctl"
else
  GOCTL_BIN="$HOME/go/bin/goctl"
fi

echo "[1/2] generating go-zero api code into $OUT_DIR ..."
"$GOCTL_BIN" api go -api "$API_FILE" -dir "$OUT_DIR"

echo "[2/2] done."
echo "next: cd $ROOT_DIR && go mod tidy"
