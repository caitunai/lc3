import test from 'node:test'
import assert from 'node:assert/strict'
import { LC3Encoder, LC3Decoder } from '../dist/index.js'

test('LC3 browser package exports encoder and decoder classes', () => {
  assert.equal(typeof LC3Encoder.create, 'function')
  assert.equal(typeof LC3Decoder.create, 'function')
})

test('LC3 browser wrapper encodes, decodes, and closes via module factory', async () => {
  const calls = []
  const memory = {
    next: 8,
    HEAPF32: new Float32Array(4096),
    HEAPU8: new Uint8Array(4096)
  }
  const moduleFactory = async () => ({
    HEAPF32: memory.HEAPF32,
    HEAPU8: memory.HEAPU8,
    _malloc(size) {
      const ptr = memory.next
      memory.next += size + 8
      return ptr
    },
    _free(ptr) {
      calls.push(['free', ptr])
    },
    _lc3_js_encoder_create(sampleRate, frameDurationUs, bitrate, frameBytes) {
      calls.push(['encoder_create', sampleRate, frameDurationUs, bitrate, frameBytes])
      return 100
    },
    _lc3_js_encoder_frame_size() {
      return 4
    },
    _lc3_js_encoder_frame_bytes() {
      return 3
    },
    _lc3_js_encode(_handle, _pcmPtr, outPtr) {
      memory.HEAPU8.set([1, 2, 3], outPtr)
      return 3
    },
    _lc3_js_encoder_destroy(handle) {
      calls.push(['encoder_destroy', handle])
    },
    _lc3_js_decoder_create(sampleRate, frameDurationUs, bitrate, frameBytes) {
      calls.push(['decoder_create', sampleRate, frameDurationUs, bitrate, frameBytes])
      return 200
    },
    _lc3_js_decoder_frame_size() {
      return 4
    },
    _lc3_js_decoder_frame_bytes() {
      return 3
    },
    _lc3_js_decode(_handle, _framePtr, _frameLen, pcmPtr) {
      memory.HEAPF32.set([0.1, 0.2, 0.3, 0.4], pcmPtr >> 2)
      return 4
    },
    _lc3_js_decoder_destroy(handle) {
      calls.push(['decoder_destroy', handle])
    }
  })

  const encoder = await LC3Encoder.create({ moduleFactory, sampleRate: 16000 })
  const frame = encoder.encode(new Float32Array([0, 0.1, 0.2, 0.3]))
  encoder.close()

  assert.deepEqual([...frame], [1, 2, 3])

  const decoder = await LC3Decoder.create({ moduleFactory, sampleRate: 16000, frameBytes: frame.byteLength })
  const decoded = decoder.decode(frame)
  decoder.close()

  assert.deepEqual([...decoded].map((value) => Number(value.toFixed(1))), [0.1, 0.2, 0.3, 0.4])
  assert.deepEqual(calls.filter(([name]) => name.endsWith('destroy')), [
    ['encoder_destroy', 100],
    ['decoder_destroy', 200]
  ])
})
