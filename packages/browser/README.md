# @caitun/lc3

Browser LC3 encoder and decoder powered by Google `liblc3` compiled to WebAssembly.

## Install

```bash
npm install @caitun/lc3
```

## Usage

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
  bitrate: 32000,
  frameBytes: frame.byteLength
})

const decoded = decoder.decode(frame)
decoder.close()
```

If your bundler serves WASM assets from a custom location, pass `locateFile`:

```js
const encoder = await LC3Encoder.create({
  locateFile: (path) => `/assets/${path}`
})
```

## Build

```bash
npm run build
```

`build:wasm` requires Emscripten `emcc`.
