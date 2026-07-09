#!/usr/bin/env sh
set -eu

PACKAGE_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
ROOT_DIR="$(CDPATH= cd -- "$PACKAGE_DIR/../.." && pwd)"
OUT_DIR="$PACKAGE_DIR/dist/wasm"
export EM_CACHE="$ROOT_DIR/.cache/emscripten"
mkdir -p "$EM_CACHE"
mkdir -p "$OUT_DIR"

emcc \
  "$ROOT_DIR/internal/csrc/liblc3/src/attdet.c" \
  "$ROOT_DIR/internal/csrc/liblc3/src/bits.c" \
  "$ROOT_DIR/internal/csrc/liblc3/src/bwdet.c" \
  "$ROOT_DIR/internal/csrc/liblc3/src/energy.c" \
  "$ROOT_DIR/internal/csrc/liblc3/src/lc3.c" \
  "$ROOT_DIR/internal/csrc/liblc3/src/ltpf.c" \
  "$ROOT_DIR/internal/csrc/liblc3/src/mdct.c" \
  "$ROOT_DIR/internal/csrc/liblc3/src/plc.c" \
  "$ROOT_DIR/internal/csrc/liblc3/src/sns.c" \
  "$ROOT_DIR/internal/csrc/liblc3/src/spec.c" \
  "$ROOT_DIR/internal/csrc/liblc3/src/tables.c" \
  "$ROOT_DIR/internal/csrc/liblc3/src/tns.c" \
  "$PACKAGE_DIR/src/lc3_wasm.c" \
  -I"$ROOT_DIR/internal/csrc/liblc3/include" \
  -I"$ROOT_DIR/internal/csrc/liblc3/src" \
  -O3 \
  -ffast-math \
  -s WASM=1 \
  -s MODULARIZE=1 \
  -s EXPORT_ES6=1 \
  -s ENVIRONMENT=web,worker \
  -s ALLOW_MEMORY_GROWTH=1 \
  -s EXPORTED_RUNTIME_METHODS='["HEAPF32","HEAPU8"]' \
  -s EXPORTED_FUNCTIONS='["_malloc","_free","_lc3_js_encoder_create","_lc3_js_encoder_frame_size","_lc3_js_encoder_frame_bytes","_lc3_js_encode","_lc3_js_encoder_destroy","_lc3_js_decoder_create","_lc3_js_decoder_frame_size","_lc3_js_decoder_frame_bytes","_lc3_js_decode","_lc3_js_decoder_destroy"]' \
  -o "$OUT_DIR/lc3-wasm.js"
