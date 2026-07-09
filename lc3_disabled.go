//go:build !cgo

package lc3

import "errors"

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

type Encoder struct{}

type Decoder struct{}

func NewEncoder(Config) (*Encoder, error) {
	return nil, ErrInit
}

func NewDecoder(Config) (*Decoder, error) {
	return nil, ErrInit
}

func (*Encoder) Encode([]float32) ([]byte, error) {
	return nil, ErrEncode
}

func (*Decoder) Decode([]byte) ([]float32, error) {
	return nil, ErrDecode
}

func (*Encoder) Close() error {
	return nil
}

func (*Decoder) Close() error {
	return nil
}

func (*Encoder) FrameSize() int {
	return 0
}

func (*Encoder) FrameBytes() int {
	return 0
}

func (*Decoder) FrameSize() int {
	return 0
}

func (*Decoder) FrameBytes() int {
	return 0
}
