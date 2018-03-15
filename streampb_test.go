package streampb

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
)

func TestRoundtripEncodeDecode(t *testing.T) {
	want := ptypes.DurationProto(42 * time.Minute)

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	err := enc.Encode(want)
	if err != nil {
		t.Fatal(err)
	}

	var got duration.Duration

	dec := NewDecoder(buf)
	if err := dec.Decode(&got); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(&got, want) {
		t.Fatalf("got: %#v want %#v", &got, want)
	}

	// Reading at the end should be an EOF
	err = dec.Decode(&got)
	if err != io.EOF {
		t.Fatalf("got err: %v want EOF", err)
	}
}

func BenchmarkDecoder(b *testing.B) {
	d := ptypes.DurationProto(42 * time.Minute)

	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	for i := 0; i < b.N; i++ {
		if err := enc.Encode(d); err != nil {
			b.Fatal(err)
		}
	}

	var dur duration.Duration
	dec := NewDecoder(buf)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if err := dec.Decode(&dur); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}
