package lc3

import (
	"errors"
	"math"
	"testing"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	encoder, err := NewEncoder(Config{})
	if err != nil {
		if errors.Is(err, ErrInit) {
			t.Skipf("lc3 unavailable: %v", err)
		}
		t.Fatalf("new encoder: %v", err)
	}
	defer func() {
		if err := encoder.Close(); err != nil {
			t.Fatalf("close encoder: %v", err)
		}
	}()

	decoder, err := NewDecoder(Config{FrameBytes: encoder.FrameBytes()})
	if err != nil {
		t.Fatalf("new decoder: %v", err)
	}
	defer func() {
		if err := decoder.Close(); err != nil {
			t.Fatalf("close decoder: %v", err)
		}
	}()

	pcm := make([]float32, encoder.FrameSize())
	for i := range pcm {
		pcm[i] = float32(math.Sin(float64(i) * 0.04))
	}
	frame, err := encoder.Encode(pcm)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if len(frame) != encoder.FrameBytes() {
		t.Fatalf("frame bytes = %d, want %d", len(frame), encoder.FrameBytes())
	}
	decoded, err := decoder.Decode(frame)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(decoded) != encoder.FrameSize() {
		t.Fatalf("decoded samples = %d, want %d", len(decoded), encoder.FrameSize())
	}
}

func TestInvalidConfig(t *testing.T) {
	_, err := NewEncoder(Config{SampleRate: 11025})
	if errors.Is(err, ErrInit) {
		t.Skipf("lc3 unavailable: %v", err)
	}
	if !errors.Is(err, ErrConfig) {
		t.Fatalf("error = %v, want ErrConfig", err)
	}
}
