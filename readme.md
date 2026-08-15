## Snowflake
Package snowflake provides a Twitter Snowflake unique id generator with a base 54 encoding and helper methods to encode common cloud regions and enumerate machines.

See the [documentation in `snowflake.go`](https://github.com/phasesoftware/snowflake/blob/master/snowflake.go) for details.

```sh
go get -u github.com/phasesoftware/snowflake@v1.0.1
```

**Notes**:
- Only uses easily readable characters for humans.
- Encoded id length is at most 11 chars (`log(54,9223372036854775807)<11`).
- Supports ~140 years runtime from Epoch using 42 timestamp bits.
- Zero allocations for generation and decoding; encoding allocates only the returned string.

**Benchmarks**:
```sh
$ go test -bench=. -benchmem -v
...
goos: darwin
goarch: arm64
pkg: github.com/phasesoftware/snowflake
cpu: Apple M4 Pro

BenchmarkGenerate-14       4920823   244.000 ns/op  0 B/op  0 allocs/op
BenchmarkBase54-14       176290560     6.784 ns/op  0 B/op  0 allocs/op
BenchmarkBaseDecode-14   314045384     3.827 ns/op  0 B/op  0 allocs/op
BenchmarkBaseEncode-14    40301504     29.54 ns/op 16 B/op  1 allocs/op
```
