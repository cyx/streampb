# streampb

`streampb` is a library for streaming protobuf messages. It frames messages using
a big endian length prefix [as described in the official docs](https://developers.google.com/protocol-buffers/docs/techniques#streaming)

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
