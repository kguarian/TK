#!/usr/bin/bash

GOOS=js GOARCH=wasm go build -o=wasm_lib.wasm wasmBuild
sudo cp wasm_lib.wasm ../pkg/wasm_lib.wasm
