package service_test

import (
	"errors"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/protocol"
)



func TestBug25_CodecKeepsBase64Cause(t *testing.T) {
    _,err:=protocol.NewCodec("q").DecodeEnvelope("q.not-base64!");if errors.Unwrap(err)==nil{t.Fatalf("unwrap=%v",err)}
}
