package snowflake

import (
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		id     ID
		verify string
	}{
		{ID(123), "21"},
		{ID(123123), "6vF"},
		{ID(123123123), "nHW1a"},
		{ID(1820096636282474496), "efUzLtM5yvu"},
		{ID(9223372036854775807), "EZNmktHEz5H"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test_%d", int64(test.id)), func(t *testing.T) {
			slow, err1 := test.id.baseEncode(54, alphabet)
			fast, err2 := test.id.base54()

			if err1 != nil || err2 != nil {
				t.Errorf("encoding failed: %v, %v", err1, err2)
			} else if slow != fast || fast != test.verify {
				t.Errorf("got '%s' and '%s', want '%s'", fast, slow, test.verify)
			}
		})
	}
}

// 49.10 ns/op
func BenchmarkBaseEncode(b *testing.B) {
	id := ID(1820096636282474496)
	for i := 0; i < b.N; i++ {
		_, _ = id.baseEncode(54, debugAlphabet)
	}
}

// 10.84 ns/op
func BenchmarkBase54(b *testing.B) {
	id := ID(1820096636282474496)
	for i := 0; i < b.N; i++ {
		_, _ = id.base54()
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		verify  ID
		encoded string
	}{
		{ID(123), "21"},
		{ID(123123), "6vF"},
		{ID(123123123), "nHW1a"},
		{ID(1820096636282474496), "efUzLtM5yvu"},
		{ID(9223372036854775807), "EZNmktHEz5H"},
		{Invalid, "xZNmktHEz5H"}, // overflow
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test_Decode_%s", test.encoded), func(t *testing.T) {
			id, err := Parse(test.encoded)

			if err != nil && test.verify != Invalid {
				t.Errorf("decoding failed: %v", err)
			} else if id != test.verify {
				t.Errorf("got '%v', want '%v'", int64(id), int64(test.verify))
			}
		})
	}
}

// 4.833 ns/op
func BenchmarkBaseDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Parse("8FaPRNs8Uks")
	}
}
