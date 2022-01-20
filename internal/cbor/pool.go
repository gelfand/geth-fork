package cbor

import (
	"fmt"
	"io"

	"github.com/ugorji/go/codec"
)

// Constant bufSize for both encoders and decoders.
const bufSize = 1 << 16

// Pool of the decoders.
var decoderPool = make(chan *codec.Decoder, 128)

// Decoder returns a new decoder that reads from the io.Reader.
func Decoder(r io.Reader) (d *codec.Decoder) {
	select {
	case d = <-decoderPool:
		d.Reset(r)
	default:
		{
			var handle codec.CborHandle
			handle.ReaderBufferSize = bufSize
			handle.ZeroCopy = true
			d = codec.NewDecoder(r, &handle)
		}
	}
	return
}

// DecoderBytes returns a new decoder that reads from the byte slice.
func DecoderBytes(r []byte) (d *codec.Decoder) {
	select {
	case d = <-decoderPool:
		d.ResetBytes(r)
	default:
		{
			var handle codec.CborHandle
			handle.ReaderBufferSize = bufSize
			handle.ZeroCopy = true
			d = codec.NewDecoderBytes(r, &handle)
		}
	}
	return
}

func returnDecoderToPool(d *codec.Decoder) {
	select {
	case decoderPool <- d:
	default:
	}
}

var encoderPool = make(chan *codec.Encoder, 128)

// Encoder returns a new encoder that writes into io.Writer.
func Encoder(w io.Writer) (e *codec.Encoder) {
	select {
	case e = <-encoderPool:
		e.Reset(w)
	default:
		{
			var handle codec.CborHandle
			handle.WriterBufferSize = bufSize
			handle.StructToArray = true
			handle.OptimumSize = true
			handle.StringToRaw = true
			e = codec.NewEncoder(w, &handle)
		}
	}
	return
}

// EncoderBytes returns a new encoder that writes into
// byte slice that w points to.
func EncoderBytes(w *[]byte) (e *codec.Encoder) {
	select {
	case e = <-encoderPool:
		e.ResetBytes(w)
	default:
		{
			var handle codec.CborHandle
			handle.WriterBufferSize = bufSize
			handle.StructToArray = true
			handle.OptimumSize = true
			handle.StringToRaw = true
			e = codec.NewEncoderBytes(w, &handle)
		}
	}
	return
}

func returnEncoderToPool(e *codec.Encoder) {
	select {
	case encoderPool <- e:
	default:
	}
}

// Return puts given Encoder or Decoder into the appropriate pool.
func Return(d interface{}) {
	switch toReturn := d.(type) {
	case *codec.Decoder: // nolint
		returnDecoderToPool(toReturn)
	case *codec.Encoder: // nolint
		returnEncoderToPool(toReturn)
	default:
		panic(fmt.Sprintf("unexpected type: %T", d))
	}
}
