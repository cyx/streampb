package streampb

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
)

// NewEncoder creates a streaming protobuf encoder.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encoder wraps an underlying io.Writer and allows you to stream
// proto encodings on it.
type Encoder struct {
	w io.Writer
}

// Encode takes any proto.Message and streams it to the underlying writer.
// Messages are framed with a length prefix.
func (e *Encoder) Encode(msg proto.Message) error {
	pbdata, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	size := uint64(len(pbdata))
	maxSize := size + binary.MaxVarintLen64
	buf := make([]byte, maxSize)
	n := binary.PutUvarint(buf, size)
	buf = append(buf[:n], pbdata...)

	_, err = e.w.Write(buf)
	return err
}

// NewDecoder creates a streaming protobuf decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decoder wraps an underlying io.Reader and allows you to stream
// proto decodings on it.
type Decoder struct {
	r      io.Reader
	buf    []byte
	bufcap uint64
}

// Decode takes a proto.Message and unmarshals the next payload in the
// underlying io.Reader. It returns an EOF when it's done.
func (d *Decoder) Decode(v proto.Message) error {
	size, err := binary.ReadUvarint(d.r.(io.ByteReader))
	if err != nil {
		return err
	}

	if size > d.bufcap {
		d.buf = make([]byte, size)
		d.bufcap = size
	}

	_, err = io.ReadFull(d.r, d.buf[:size])
	if err != nil {
		return err
	}

	return proto.Unmarshal(d.buf[:size], v)
}

func translateError(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}
