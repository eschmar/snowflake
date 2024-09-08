## Snowflake
Package snowflake provides a Twitter Snowflake unique id generator with a base 54 encoding and helper methods to encode common cloud regions and enumerate machines.

See the [documentation in `snowflake.go`](https://github.com/eschmar/snowflake/blob/master/snowflake.go) for details.

Notes:
- Encoded length is at most 11 chars (`log(54,9223372036854775807)<11`).
- Supports ~140 years runtime from Epoch.
- Zero allocations for generation, encoding and decoding.

Benchmarks:
```
goos: darwin
goarch: arm64
pkg: github.com/eschmar/snowflake
cpu: Apple M2
=== RUN
BenchmarkGenerate-8      4923445      244.1 ns/op    0 B/op    0 allocs/op
BenchmarkBase54-8        110005624    10.81 ns/op    0 B/op    0 allocs/op
BenchmarkBaseDecode-8    238005690    4.854 ns/op    0 B/op    0 allocs/op
```
