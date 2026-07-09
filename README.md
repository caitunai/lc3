# lc3

Go and browser LC3 encoder/decoder packages backed by Google `liblc3`.

This repository publishes two artifacts from one source tree:

- Go module: `github.com/caitunai/lc3`
- npm package: `@caitun/lc3`

The bundled `liblc3` source lives in `internal/csrc/liblc3`, so users installing the Go module do not need git submodules.

## Go

Install:

```bash
go get github.com/caitunai/lc3
```

Encode and decode one LC3 frame:

```go
package main

import "github.com/caitunai/lc3"

func main() {
	encoder, err := lc3.NewEncoder(lc3.Config{
		SampleRate:      16000,
		FrameDurationUS: 10000,
		Bitrate:         32000,
	})
	if err != nil {
		panic(err)
	}
	defer encoder.Close()

	pcm := make([]float32, encoder.FrameSize())
	frame, err := encoder.Encode(pcm)
	if err != nil {
		panic(err)
	}

	decoder, err := lc3.NewDecoder(lc3.Config{
		SampleRate:      16000,
		FrameDurationUS: 10000,
		FrameBytes:      len(frame),
	})
	if err != nil {
		panic(err)
	}
	defer decoder.Close()

	decoded, err := decoder.Decode(frame)
	if err != nil {
		panic(err)
	}
	_ = decoded
}
```

The Go package uses cgo. With `CGO_ENABLED=0`, constructors return `ErrInit` and the package still compiles.

## Browser

Install:

```bash
npm install @caitun/lc3
```

Use:

```js
import { LC3Encoder, LC3Decoder } from '@caitun/lc3'

const encoder = await LC3Encoder.create({
  sampleRate: 16000,
  frameDurationUs: 10000,
  bitrate: 32000
})

const pcm = new Float32Array(encoder.frameSize)
const frame = encoder.encode(pcm)
encoder.close()

const decoder = await LC3Decoder.create({
  sampleRate: 16000,
  frameDurationUs: 10000,
  frameBytes: frame.byteLength
})

const decoded = decoder.decode(frame)
decoder.close()
```

The npm package includes prebuilt Emscripten WASM files under `dist/wasm`.

## Build

Go:

```bash
go test ./...
CGO_ENABLED=0 go test ./...
```

Browser:

```bash
cd packages/browser
npm run build
npm test
```

`npm run build:wasm` requires Emscripten `emcc`.

## Release

The Go module is released by pushing a semantic version tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The GitHub Actions workflow also publishes the browser package to npm when a
`v*` tag is pushed. Before tagging, make sure the npm version matches the tag:

```bash
cd packages/browser
npm version 0.1.0 --no-git-tag-version
```

The publishing job:

```bash
npm run build
npm test
npm pack --dry-run
npm publish --access public --provenance
```

## License

This repository is MIT licensed. The bundled Google `liblc3` source is Apache-2.0 licensed; see `internal/csrc/liblc3/LICENSE`.
