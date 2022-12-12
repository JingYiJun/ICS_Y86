# ICS_Y86_Backend

An ICS_Y86 Backend Written in Golang

Stage 2

## Features

- WASM, packed with frontend

## Usage

### prerequisite

Go 1.19

### build and run

```shell
export GOARCH=wasm
export GOOS=js
go build -o .\build\ICS_Y86_Backend.wasm
```

导入到前端并加载

教程：[Go WebAssembly (Wasm) 简明教程](https://geektutu.com/post/quick-go-wasm.html)

## License

MIT License

Copyright (c) 2022-present ck ct fkx