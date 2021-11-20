// Package cbor implements encoding and decoding as defined in RFC7049.

package cbor

import (
	"io"
)

// Marshal writes CBOR encoding of v into the given io.Writer.
func Marshal(dst io.Writer, v interface{}) error {
	e := Encoder(dst)
	err := e.Encode(v)
	returnEncoderToPool(e)

	return err
}

// Unmarshal parses CBOR-encoded data and stores the result
// in the value pointed to by v.
func Unmarshal(data io.Reader, v interface{}) error {
	d := Decoder(data)
	err := d.Decode(v)
	returnDecoderToPool(d)

	return err
}

// MustMarshal marshals v and calls panic if it returns non-nil error.
func MustMarshal(dst io.Writer, v interface{}) {
	if err := Marshal(dst, v); err != nil {
		panic(err)
	}
}

// MustUnmarshal unmarshals data and calls panic if it returns non-nil error.
func MustUnmarshal(data io.Reader, v interface{}) {
	if err := Unmarshal(data, v); err != nil {
		panic(err)
	}
}
