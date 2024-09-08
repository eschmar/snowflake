// Package snowflake provides a Twitter Snowflake generator
// with a base 54 encoding and helper methods to encode
// common cloud regions and enumerate machines.
//
// Learn more about Twitter Snowflake at [Wikipedia] and the
// original blog post [Twitter Engineering] from 2010.
//
// The following format is used:
// +----------------------------------------------------------------------+
// | 1 bit Unused | 42 bit Timestamp | 9 bit Machine ID | 12 bit Sequence |
// +----------------------------------------------------------------------+
//
// Example:
//   - ID in64: 305023354946072576
//   - Base 54 Encoded: rAuFzp2Kghj
//
// Extracted:
//   - Timestamp: 2024-08-10 09:47:50.758 +0000 UTC
//   - Machine ID: 35
//   - Sequence: 0
//
// It differs from the usual Twitter encoding, such that the machine id
// encoding takes 9 bits instead of 10, in order to have one more bit for
// the timestamp. This enables timestamps of ~140 years from epoch.
//
// Notes:
//   - Intended usage is int64 internally, base encoding user facing.
//   - IDs are roughly sortable, but timestamps are not deducible without effort
//     since a scrambled alphabet is used for encoding.
//   - Base 54 only uses characters easily readable for humans.
//   - Encoded length is at most 11 chars (log(54,9223372036854775807)<11).
//   - Machine ID uses 3 bits for the continent, 6 bits for machine enumeration.
//     That means at most 64 machines per continent.
//   - Machines can generate a new ID without coordination, however
//     a unique machine ID is required on startup.
//   - Uses monotonic clock when available to avoid duplicate ids.
//
// [Wikipedia]: https://en.wikipedia.org/wiki/Snowflake_ID
// [Twitter Engineering]: https://blog.x.com/engineering/en_us/a/2010/announcing-snowflake
package snowflake

import (
	"encoding/json"
	"math"
	"sync"
	"time"
)

// Snowflake ID
type ID int64

const Invalid ID = ID(-1)

// Wed Jan 01 2020 00:00:01.000 UTC
// time.Date(2020, time.January, 1, 0, 0, 1, 0, time.UTC).UnixMilli()
const Epoch int64 = 1577836801000

// Number of bits to encode timestamp, defined as the difference in milliseconds between current timestamp and Epoch. Max date is therefore 2159-05-15 07:35:12.103 +0000 UTC.
const bitsTimestamp int64 = 42

// Number of bits to encode the machine ID.
// 3 bits continent + 6 bits machine id, meaning max 64 machines per continent.
const bitsMachineID int64 = 9

// Number of bits to encode sequence number, if more than one ID was generated within the same millisecond.
const bitsMachineSequence int64 = 12

// Internal variables for snowflake ID generation.
var epoch time.Time
var machineId int64 = 0
var machineSequence int64 = 0

var bitMapMachineId, bitMapMachineSequence int64
var mutex sync.Mutex
var previous int64

func init() {
	// Sanity check if encoding fits in signed int64
	if bitsTimestamp+bitsMachineID+bitsMachineSequence != 63 {
		panic("invalid snowflake bit length")
	}

	// `epoch` is `Epoch` + monotonic information. A monotonic clock
	// exclusively moves forward, unlike a wall clock that can be adjusted
	// backwards. In such case, there is a chance of duplicate IDs.
	now := time.Now()
	epoch = now.Add(time.UnixMilli(Epoch).Sub(now))

	// Prepare bitmaps for bitwise operation
	bitMapMachineId = int64(math.Pow(2, float64(bitsMachineID))) - 1
	bitMapMachineSequence = int64(math.Pow(2, float64(bitsMachineSequence))) - 1

	// Pre-populates `decodeMap` to speed up parsing.
	initDecodeMap()
}

// Sets the unique machine id for snowflake generation.
// ATTENTION: If more than one server is using the same
// machine id in parallel, then the uniqueness of any
// snowflake ID can _NOT_ be guaranteed.
func SetMachineId(region string, index int64) {
	continent := getContinentCode(region)
	maxMachineNumber := int64(math.Pow(2, float64(bitsMachineID-3)))

	if continent < 0 || index < 0 || index >= maxMachineNumber {
		panic("unable to determine proper machine id")
	}

	machineId = ((continent & 0b111) << (bitsMachineID - 3)) | (index & (maxMachineNumber - 1))
}

// Generates a unique snowflake id.
func Generate() ID {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Since(epoch).Milliseconds()

	if now == previous && machineSequence == bitMapMachineSequence {
		// Reached max squence number 2^{BitsMachineSequence}.
		// Wait for the next millisecond.
		for now <= previous {
			now = time.Since(epoch).Milliseconds()
		}
	} else if now > previous {
		// Reset machine sequence for new millisecond
		machineSequence = -1
	} else if now < previous {
		// Avoid potential duplicates
		panic("attempted to generate snowflake id of the past")
	}

	// Increment machine sequence
	machineSequence = (machineSequence + 1) & bitMapMachineSequence

	// Update latest ID timestamp
	previous = now

	// Return snowflake
	return ID(now<<(bitsMachineID+bitsMachineSequence) |
		(machineId << bitsMachineSequence) |
		machineSequence)
}

// Returns the base encoded representation of a snowflake ID.
func (id ID) String() string {
	encoded, err := id.base54()
	if err != nil {
		return ""
	}

	return encoded
}

// Converts a base encoded string into a snowflake ID.
func Parse(input string) (ID, error) {
	return decode54([]byte(input))
}

// Extracts timestamp from a snowflake.
func (id ID) Time() int64 {
	return (int64(id) >> (bitsMachineID + bitsMachineSequence)) + Epoch
}

// Extracts machine id from a snowflake.
func (id ID) MachineId() int64 {
	return (int64(id) >> bitsMachineSequence) & bitMapMachineId
}

// Extracts sequence number from a snowflake.
func (id ID) MachineSequence() int64 {
	return int64(id) & bitMapMachineSequence
}

//
// Marshaler interface implementation
//

// ID to JSON marshalling.
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

// JSON to ID unmarshalling.
func (id *ID) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
		*id = Invalid
		return &ErrorInvalidJson
	}

	parsed, err := Parse(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}

	*id = parsed
	return nil
}
