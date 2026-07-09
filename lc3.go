//go:build cgo

package lc3

/*
#cgo CFLAGS: -I${SRCDIR}/internal/csrc/liblc3/include -I${SRCDIR}/internal/csrc/liblc3/src -ffast-math
#cgo !darwin LDFLAGS: -lm
#include <stdlib.h>
#include "lc3.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

const (
	DefaultSampleRate      = 16000
	DefaultFrameDurationUS = 10000
	DefaultBitrate         = 32000
	MinFrameBytes          = 20
	MaxFrameBytes          = 400
)

var (
	ErrConfig = errors.New("lc3 config invalid")
	ErrInit   = errors.New("lc3 init failed")
	ErrEncode = errors.New("lc3 encode failed")
	ErrDecode = errors.New("lc3 decode failed")
)

type Config struct {
	SampleRate      int
	FrameDurationUS int
	Bitrate         int
	FrameBytes      int
}

type Encoder struct {
	handle          C.lc3_encoder_t
	mem             unsafe.Pointer
	frameSize       int
	frameBytes      int
	sampleRate      int
	frameDurationUS int
	bitrate         int
	closed          bool
}

type Decoder struct {
	handle          C.lc3_decoder_t
	mem             unsafe.Pointer
	frameSize       int
	frameBytes      int
	sampleRate      int
	frameDurationUS int
	closed          bool
}

func NewEncoder(cfg Config) (*Encoder, error) {
	params, err := normalizeConfig(cfg)
	if err != nil {
		return nil, err
	}
	size := C.lc3_encoder_size(C.int(params.FrameDurationUS), C.int(params.SampleRate))
	if size == 0 {
		return nil, ErrConfig
	}
	mem := C.malloc(C.size_t(size))
	if mem == nil {
		return nil, ErrInit
	}
	handle := C.lc3_setup_encoder(C.int(params.FrameDurationUS), C.int(params.SampleRate), 0, mem)
	if handle == nil {
		C.free(mem)
		return nil, ErrInit
	}
	return &Encoder{
		handle:          handle,
		mem:             mem,
		frameSize:       params.frameSize,
		frameBytes:      params.FrameBytes,
		sampleRate:      params.SampleRate,
		frameDurationUS: params.FrameDurationUS,
		bitrate:         params.Bitrate,
	}, nil
}

func NewDecoder(cfg Config) (*Decoder, error) {
	params, err := normalizeConfig(cfg)
	if err != nil {
		return nil, err
	}
	size := C.lc3_decoder_size(C.int(params.FrameDurationUS), C.int(params.SampleRate))
	if size == 0 {
		return nil, ErrConfig
	}
	mem := C.malloc(C.size_t(size))
	if mem == nil {
		return nil, ErrInit
	}
	handle := C.lc3_setup_decoder(C.int(params.FrameDurationUS), C.int(params.SampleRate), 0, mem)
	if handle == nil {
		C.free(mem)
		return nil, ErrInit
	}
	return &Decoder{
		handle:          handle,
		mem:             mem,
		frameSize:       params.frameSize,
		frameBytes:      params.FrameBytes,
		sampleRate:      params.SampleRate,
		frameDurationUS: params.FrameDurationUS,
	}, nil
}

func (e *Encoder) Encode(samples []float32) ([]byte, error) {
	if e == nil || e.closed || e.handle == nil {
		return nil, ErrEncode
	}
	if len(samples) != e.frameSize {
		return nil, ErrConfig
	}
	out := make([]byte, e.frameBytes)
	if C.lc3_encode(
		e.handle,
		C.LC3_PCM_FORMAT_FLOAT,
		unsafe.Pointer(&samples[0]),
		1,
		C.int(e.frameBytes),
		unsafe.Pointer(&out[0]),
	) != 0 {
		return nil, ErrEncode
	}
	return out, nil
}

func (d *Decoder) Decode(frame []byte) ([]float32, error) {
	if d == nil || d.closed || d.handle == nil {
		return nil, ErrDecode
	}
	if len(frame) == 0 {
		return nil, nil
	}
	if len(frame) > MaxFrameBytes {
		return nil, ErrConfig
	}
	pcm := make([]float32, d.frameSize)
	ret := C.lc3_decode(
		d.handle,
		unsafe.Pointer(&frame[0]),
		C.int(len(frame)),
		C.LC3_PCM_FORMAT_FLOAT,
		unsafe.Pointer(&pcm[0]),
		1,
	)
	if ret < 0 {
		return nil, ErrDecode
	}
	return pcm, nil
}

func (e *Encoder) Close() error {
	if e == nil || e.closed {
		return nil
	}
	e.closed = true
	if e.mem != nil {
		C.free(e.mem)
		e.mem = nil
		e.handle = nil
	}
	return nil
}

func (d *Decoder) Close() error {
	if d == nil || d.closed {
		return nil
	}
	d.closed = true
	if d.mem != nil {
		C.free(d.mem)
		d.mem = nil
		d.handle = nil
	}
	return nil
}

func (e *Encoder) FrameSize() int {
	if e == nil {
		return 0
	}
	return e.frameSize
}

func (e *Encoder) FrameBytes() int {
	if e == nil {
		return 0
	}
	return e.frameBytes
}

func (d *Decoder) FrameSize() int {
	if d == nil {
		return 0
	}
	return d.frameSize
}

func (d *Decoder) FrameBytes() int {
	if d == nil {
		return 0
	}
	return d.frameBytes
}

type normalizedConfig struct {
	Config
	frameSize int
}

func normalizeConfig(cfg Config) (normalizedConfig, error) {
	if cfg.SampleRate == 0 {
		cfg.SampleRate = DefaultSampleRate
	}
	if cfg.FrameDurationUS == 0 {
		cfg.FrameDurationUS = DefaultFrameDurationUS
	}
	if cfg.Bitrate == 0 {
		cfg.Bitrate = DefaultBitrate
	}

	frameSize := int(C.lc3_frame_samples(C.int(cfg.FrameDurationUS), C.int(cfg.SampleRate)))
	if frameSize <= 0 {
		return normalizedConfig{}, ErrConfig
	}
	if cfg.FrameBytes == 0 {
		cfg.FrameBytes = int(C.lc3_frame_bytes(C.int(cfg.FrameDurationUS), C.int(cfg.Bitrate)))
	}
	if cfg.FrameBytes < MinFrameBytes || cfg.FrameBytes > MaxFrameBytes {
		return normalizedConfig{}, ErrConfig
	}
	return normalizedConfig{
		Config:    cfg,
		frameSize: frameSize,
	}, nil
}
