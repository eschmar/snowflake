package snowflake

import (
	"math"
)

// Scrambled version of "0123456789abcdefghjkmnprstuvwxyzACDEFGHJKLMNPQRTUVWXYZ".
const alphabet string = "g82FcYyTeUr0vsn1Jb9NmLMPuHGhVztRp4f3jDk5Zd6ECaw7AWQKXx"

// Supports up to base 84. Tests will fail if map is changed (intentionally).
const debugAlphabet string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-._~:?#[]@!$&'()*+,;%="

// Lookup alphabet char to its position in the alphabet.
var decodeMap [256]byte

// Pre-populates `decodeMap` to speed up parsing.
// ~20x speedup using [256]byte lookup compared to map[byte]byte.
func initDecodeMap() {
	// Invalid characters are marked with 0xFF.
	for i := 0; i < len(decodeMap); i++ {
		decodeMap[i] = 0xFF
	}

	for i := 0; i < len(alphabet); i++ {
		decodeMap[alphabet[i]] = byte(i)
	}
}

// Returns the base 54 encoded representation of a snowflake.
func (id ID) base54() (string, error) {
	if id < 0 {
		return "", &ErrorInvalid
	} else if id < 54 {
		return string(alphabet[id]), nil
	}

	// 11 is ceil(log(54, MAX_INT64))
	b := make([]byte, 11)
	i := 10

	for id >= 54 {
		b[i] = alphabet[id%54]
		id /= 54
		i--
	}

	b[i] = alphabet[id]

	return string(b[i:]), nil
}

// Converts a base 54 encoded string into a snowflake ID.
func decode54(b []byte) (ID, error) {
	var id int64

	for i := range b {
		if decodeMap[b[i]] == 0xFF {
			return Invalid, &ErrorInvalidByte
		}

		// Example: 'Wef' is [42 12 3]
		// 42*54^2 + 12*54 + 3 == 123123 == (((42*54) + 12) * 54 + 3)
		id = id*54 + int64(decodeMap[b[i]])
	}

	// Overflow sanity check, slow down ~0.2 ns/op
	if id < 0 {
		return Invalid, &ErrorInvalid
	}

	return ID(id), nil
}

// Private encode method for testing and verifying different
// bases and encoding alphabets. Not optimised.
func (id ID) baseEncode(base int64, encodeMap string) (string, error) {
	if id < 0 {
		return "", &ErrorInvalid
	} else if int(base) > len(encodeMap) {
		return "", &ErrorEncodeMapLength
	}

	val := int64(id)

	if val < base {
		return string(encodeMap[id]), nil
	}

	length := int(math.Ceil(logBase(float64(base), float64(id))))
	b := make([]byte, length)
	i := length - 1

	for val >= base {
		b[i] = encodeMap[val%base]
		val /= base
		i--
	}

	b[i] = encodeMap[val%base]
	return string(b), nil
}

// log(base, x) = log(x) / log(base)
func logBase(base float64, x float64) float64 {
	return math.Log(x) / math.Log(base)
}
