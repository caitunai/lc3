#include <stdlib.h>
#include "lc3.h"

typedef struct {
  lc3_encoder_t encoder;
  lc3_decoder_t decoder;
  void *mem;
  int frame_size;
  int frame_bytes;
} Lc3JsCodec;

static int lc3_js_frame_bytes(int frame_duration_us, int bitrate, int frame_bytes) {
  if (frame_bytes > 0) {
    return frame_bytes;
  }
  return lc3_frame_bytes(frame_duration_us, bitrate);
}

Lc3JsCodec *lc3_js_encoder_create(int sample_rate, int frame_duration_us, int bitrate, int frame_bytes) {
  int resolved_frame_bytes = lc3_js_frame_bytes(frame_duration_us, bitrate, frame_bytes);
  int frame_size = lc3_frame_samples(frame_duration_us, sample_rate);
  unsigned mem_size = lc3_encoder_size(frame_duration_us, sample_rate);
  if (resolved_frame_bytes <= 0 || frame_size <= 0 || mem_size == 0) {
    return NULL;
  }

  Lc3JsCodec *codec = (Lc3JsCodec *)calloc(1, sizeof(Lc3JsCodec));
  if (codec == NULL) {
    return NULL;
  }
  codec->mem = malloc(mem_size);
  if (codec->mem == NULL) {
    free(codec);
    return NULL;
  }
  codec->encoder = lc3_setup_encoder(frame_duration_us, sample_rate, 0, codec->mem);
  if (codec->encoder == NULL) {
    free(codec->mem);
    free(codec);
    return NULL;
  }
  codec->frame_size = frame_size;
  codec->frame_bytes = resolved_frame_bytes;
  return codec;
}

int lc3_js_encoder_frame_size(Lc3JsCodec *codec) {
  if (codec == NULL) {
    return 0;
  }
  return codec->frame_size;
}

int lc3_js_encoder_frame_bytes(Lc3JsCodec *codec) {
  if (codec == NULL) {
    return 0;
  }
  return codec->frame_bytes;
}

int lc3_js_encode(Lc3JsCodec *codec, float *pcm, unsigned char *out) {
  if (codec == NULL || pcm == NULL || out == NULL || codec->encoder == NULL) {
    return -1;
  }
  if (lc3_encode(codec->encoder, LC3_PCM_FORMAT_FLOAT, pcm, 1, codec->frame_bytes, out) != 0) {
    return -1;
  }
  return codec->frame_bytes;
}

void lc3_js_encoder_destroy(Lc3JsCodec *codec) {
  if (codec == NULL) {
    return;
  }
  free(codec->mem);
  free(codec);
}

Lc3JsCodec *lc3_js_decoder_create(int sample_rate, int frame_duration_us, int bitrate, int frame_bytes) {
  int resolved_frame_bytes = lc3_js_frame_bytes(frame_duration_us, bitrate, frame_bytes);
  int frame_size = lc3_frame_samples(frame_duration_us, sample_rate);
  unsigned mem_size = lc3_decoder_size(frame_duration_us, sample_rate);
  if (resolved_frame_bytes <= 0 || frame_size <= 0 || mem_size == 0) {
    return NULL;
  }

  Lc3JsCodec *codec = (Lc3JsCodec *)calloc(1, sizeof(Lc3JsCodec));
  if (codec == NULL) {
    return NULL;
  }
  codec->mem = malloc(mem_size);
  if (codec->mem == NULL) {
    free(codec);
    return NULL;
  }
  codec->decoder = lc3_setup_decoder(frame_duration_us, sample_rate, 0, codec->mem);
  if (codec->decoder == NULL) {
    free(codec->mem);
    free(codec);
    return NULL;
  }
  codec->frame_size = frame_size;
  codec->frame_bytes = resolved_frame_bytes;
  return codec;
}

int lc3_js_decoder_frame_size(Lc3JsCodec *codec) {
  if (codec == NULL) {
    return 0;
  }
  return codec->frame_size;
}

int lc3_js_decoder_frame_bytes(Lc3JsCodec *codec) {
  if (codec == NULL) {
    return 0;
  }
  return codec->frame_bytes;
}

int lc3_js_decode(Lc3JsCodec *codec, unsigned char *frame, int frame_len, float *pcm) {
  if (codec == NULL || frame == NULL || frame_len <= 0 || pcm == NULL || codec->decoder == NULL) {
    return -1;
  }
  if (lc3_decode(codec->decoder, frame, frame_len, LC3_PCM_FORMAT_FLOAT, pcm, 1) < 0) {
    return -1;
  }
  return codec->frame_size;
}

void lc3_js_decoder_destroy(Lc3JsCodec *codec) {
  if (codec == NULL) {
    return;
  }
  free(codec->mem);
  free(codec);
}
