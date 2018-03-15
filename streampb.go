package streampb

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	// prefixSize is the number of bytes we preallocate for storing
	// our big endian lenth prefix buffer.
	prefixSize = 4

	// maxSize is the maximum length of proto messages we expect to decode.
	maxSize = 65535
)

// NewEncoder creates a streaming protobuf encoder.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w, prefixBuf: make([]byte, prefixSize)}
}

// Encoder wraps an underlying io.Writer and allows you to stream
// proto encodings on it.
type Encoder struct {
	w         io.Writer
	prefixBuf []byte
}

// Encode takes any proto.Message and streams it to the underlying writer.
// Messages are framed with a length prefix.
func (e *Encoder) Encode(msg proto.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(e.prefixBuf, uint32(len(buf)))

	if _, err := e.w.Write(e.prefixBuf); err != nil {
		return errors.Wrap(err, "failed writing length prefix")
	}

	_, err = e.w.Write(buf)
	return errors.Wrap(err, "failed writing marshaled data")
}

// NewDecoder creates a streaming protobuf decoder. It currently assumes a max
// of 64KiB for all protobuf messages.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:         r,
		prefixBuf: make([]byte, prefixSize),
		buf:       make([]byte, maxSize),
	}
}

// Decoder wraps an underlying io.Reader and allows you to stream
// proto decodings on it.
type Decoder struct {
	r         io.Reader
	prefixBuf []byte
	buf       []byte
}

// Decode takes a proto.Message and unmarshals the next payload in the
// underlying io.Reader. It returns an EOF when it's done.
func (d *Decoder) Decode(v proto.Message) error {
	_, err := io.ReadFull(d.r, d.prefixBuf)
	if err != nil {
		return err
	}

	n := binary.BigEndian.Uint32(d.prefixBuf)

	idx := uint32(0)
	for idx < n {
		m, err := d.r.Read(d.buf[idx:n])
		if err != nil {
			return errors.Wrap(translateError(err), "failed reading marshaled data")
		}
		idx += uint32(m)
	}
	return proto.Unmarshal(d.buf[:n], v)
}

func translateError(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}
