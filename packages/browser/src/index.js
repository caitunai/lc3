const LC3_MAX_FRAME_BYTES = 400

export const DEFAULT_SAMPLE_RATE = 16000
export const DEFAULT_FRAME_DURATION_US = 10000
export const DEFAULT_BITRATE = 32000

export class LC3Encoder {
  constructor(module, handle, frameSize, frameBytes) {
    this.module = module
    this.handle = handle
    this.frameSize = frameSize
    this.frameBytes = frameBytes
    this.pcmPtr = module._malloc(frameSize * 4)
    this.outPtr = module._malloc(frameBytes)
  }

  static async create(options = {}) {
    const module = await createLC3Module(options)
    const sampleRate = options.sampleRate ?? DEFAULT_SAMPLE_RATE
    const frameDurationUs = options.frameDurationUs ?? DEFAULT_FRAME_DURATION_US
    const bitrate = options.bitrate ?? DEFAULT_BITRATE
    const frameBytes = options.frameBytes ?? 0
    const handle = module._lc3_js_encoder_create(sampleRate, frameDurationUs, bitrate, frameBytes)
    if (!handle) {
      throw new Error('LC3 encoder initialization failed.')
    }
    const resolvedFrameSize = module._lc3_js_encoder_frame_size(handle)
    const resolvedFrameBytes = module._lc3_js_encoder_frame_bytes(handle)
    if (resolvedFrameSize <= 0 || resolvedFrameBytes <= 0) {
      module._lc3_js_encoder_destroy(handle)
      throw new Error('LC3 frame parameters are invalid.')
    }
    return new LC3Encoder(module, handle, resolvedFrameSize, resolvedFrameBytes)
  }

  encode(samples) {
    if (!this.handle) {
      throw new Error('LC3 encoder is closed.')
    }
    if (samples.length !== this.frameSize) {
      throw new Error(`LC3 encoder expects ${this.frameSize} samples.`)
    }
    this.module.HEAPF32.set(samples, this.pcmPtr >> 2)
    const length = this.module._lc3_js_encode(this.handle, this.pcmPtr, this.outPtr)
    if (length < 0) {
      throw new Error('LC3 encode failed.')
    }
    return this.module.HEAPU8.slice(this.outPtr, this.outPtr + length)
  }

  close() {
    if (!this.module) {
      return
    }
    if (this.handle) {
      this.module._lc3_js_encoder_destroy(this.handle)
      this.handle = 0
    }
    if (this.pcmPtr) {
      this.module._free(this.pcmPtr)
      this.pcmPtr = 0
    }
    if (this.outPtr) {
      this.module._free(this.outPtr)
      this.outPtr = 0
    }
    this.module = null
  }
}

export class LC3Decoder {
  constructor(module, handle, frameSize, frameBytes) {
    this.module = module
    this.handle = handle
    this.frameSize = frameSize
    this.frameBytes = frameBytes
    this.framePtr = module._malloc(LC3_MAX_FRAME_BYTES)
    this.pcmPtr = module._malloc(frameSize * 4)
  }

  static async create(options = {}) {
    const module = await createLC3Module(options)
    const sampleRate = options.sampleRate ?? DEFAULT_SAMPLE_RATE
    const frameDurationUs = options.frameDurationUs ?? DEFAULT_FRAME_DURATION_US
    const bitrate = options.bitrate ?? DEFAULT_BITRATE
    const frameBytes = options.frameBytes ?? 0
    const handle = module._lc3_js_decoder_create(sampleRate, frameDurationUs, bitrate, frameBytes)
    if (!handle) {
      throw new Error('LC3 decoder initialization failed.')
    }
    const resolvedFrameSize = module._lc3_js_decoder_frame_size(handle)
    const resolvedFrameBytes = module._lc3_js_decoder_frame_bytes(handle)
    if (resolvedFrameSize <= 0 || resolvedFrameBytes <= 0) {
      module._lc3_js_decoder_destroy(handle)
      throw new Error('LC3 frame parameters are invalid.')
    }
    return new LC3Decoder(module, handle, resolvedFrameSize, resolvedFrameBytes)
  }

  decode(frame) {
    if (!this.handle) {
      throw new Error('LC3 decoder is closed.')
    }
    if (frame.byteLength > LC3_MAX_FRAME_BYTES) {
      throw new Error('LC3 frame is too large.')
    }
    this.module.HEAPU8.set(frame, this.framePtr)
    const sampleCount = this.module._lc3_js_decode(this.handle, this.framePtr, frame.byteLength, this.pcmPtr)
    if (sampleCount < 0) {
      throw new Error('LC3 decode failed.')
    }
    return this.module.HEAPF32.slice(this.pcmPtr >> 2, (this.pcmPtr >> 2) + sampleCount)
  }

  close() {
    if (!this.module) {
      return
    }
    if (this.handle) {
      this.module._lc3_js_decoder_destroy(this.handle)
      this.handle = 0
    }
    if (this.framePtr) {
      this.module._free(this.framePtr)
      this.framePtr = 0
    }
    if (this.pcmPtr) {
      this.module._free(this.pcmPtr)
      this.pcmPtr = 0
    }
    this.module = null
  }
}

export async function createLC3Module(options = {}) {
  const factory = options.moduleFactory ?? await loadLC3ModuleFactory()
  return factory({
    locateFile: options.locateFile
  })
}

async function loadLC3ModuleFactory() {
  const module = await import('./wasm/lc3-wasm.js')
  return module.default
}
