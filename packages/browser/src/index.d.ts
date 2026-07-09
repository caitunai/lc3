export interface LC3ModuleOptions {
  locateFile?: (path: string) => string
  moduleFactory?: (options?: { locateFile?: (path: string) => string }) => Promise<unknown>
}

export interface LC3CodecOptions extends LC3ModuleOptions {
  sampleRate?: number
  frameDurationUs?: number
  bitrate?: number
  frameBytes?: number
}

export declare const DEFAULT_SAMPLE_RATE: 16000
export declare const DEFAULT_FRAME_DURATION_US: 10000
export declare const DEFAULT_BITRATE: 32000

export declare class LC3Encoder {
  readonly frameSize: number
  readonly frameBytes: number
  static create(options?: LC3CodecOptions): Promise<LC3Encoder>
  encode(samples: Float32Array): Uint8Array
  close(): void
}

export declare class LC3Decoder {
  readonly frameSize: number
  readonly frameBytes: number
  static create(options?: LC3CodecOptions): Promise<LC3Decoder>
  decode(frame: Uint8Array): Float32Array
  close(): void
}

export declare function createLC3Module(options?: LC3ModuleOptions): Promise<unknown>
