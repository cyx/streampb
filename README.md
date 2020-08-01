# streampb

`streampb` is a library heavily inspired from recordio. It frames messages
using [a varint](https://golang.org/pkg/encoding/binary/#Varint).

The underlying techniques used are [as described in the official docs][docs].

[docs]: https://developers.google.com/protocol-buffers/docs/techniques#streaming

## Getting

```
go get -u github.com/cyx/streampb
```

## Usage

```golang
// encoding
enc := streampb.NewEncoder(w)
enc.Encode(ptypes.DurationProto(42 * time.Minute))

// decoding
var d duration.Duration
dec := streampub.NewDecoder(r)
dec.Decode(&d)
```
