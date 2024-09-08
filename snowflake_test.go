package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	SetMachineId("arn", 35)

	tests := []ID{
		Generate(),
		Generate(),
		Generate(),
		ID(9223372036854775807),
	}

	for _, id := range tests {
		fmt.Println("ID: ", id)
		fmt.Println("ID: ", int64(id))

		encoded, err := id.base54()
		if err != nil {
			t.Errorf("encoding failed: %v", err)
		}

		fmt.Println("Encoded: ", encoded)
		fmt.Println("Time: ", time.UnixMilli(id.Time()).UTC().String())
		fmt.Println("MachineId: ", id.MachineId())
		fmt.Println("MachineSequence: ", id.MachineSequence())

		fmt.Println("---")
	}
}

func TestMachineId(t *testing.T) {
	tests := []struct {
		region string
		num    int64
	}{
		{"fra", 35},
		{"lax", 4},
		// {"unk", 0},
		// {"phx", -1},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test_MachineID_%s_%d", test.region, test.num), func(t *testing.T) {
			SetMachineId(test.region, test.num)
			fmt.Println("Machine ID: ", machineId)
			fmt.Printf("Binary:      %09b\n", machineId)
		})
	}
}

func TestGenerateExceedSequence(t *testing.T) {
	var wg sync.WaitGroup

	for j := 0; j < 4; j++ {
		wg.Add(1)

		go func() {
			for i := 0; i < 5000; i++ {
				id := Generate()
				seq := id.MachineSequence()

				// TODO: Define a more useful test.
				if seq > 4092 || seq < 3 {
					fmt.Printf("[%d]: Sequence = %d\n", j, seq)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

// 244.0 ns/op
func BenchmarkGenerate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Generate()
	}
}

//
// Marshaler interface implementation
//

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		id     ID
		verify string
	}{
		{ID(123123), `"6vF"`},
		{ID(123123123), `"nHW1a"`},
		{ID(1820096636282474496), `"efUzLtM5yvu"`},
		{ID(9223372036854775807), `"EZNmktHEz5H"`},
		{ID(305023354946072576), `"8uyZY2sj3re"`},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test_Marshal_%d", int64(test.id)), func(t *testing.T) {
			bytes, err := test.id.MarshalJSON()

			if err != nil {
				t.Errorf("marshal failed: %v", err)
			} else if string(bytes) != test.verify {
				t.Errorf("got '%v', want '%v'", string(bytes), test.verify)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		json   string
		verify ID
		err    error
	}{
		{`"6vF"`, ID(123123), nil},
		{`"nHW1a"`, ID(123123123), nil},
		{`"efUzLtM5yvu"`, ID(1820096636282474496), nil},
		{`"EZNmktHEz5H"`, ID(9223372036854775807), nil},
		{`"8HH7MXkTRtr"`, ID(310311215280041986), nil},
		{`6vF`, Invalid, &ErrorInvalidJson},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test_Unmarshal_%s", test.json), func(t *testing.T) {
			var id ID
			err := id.UnmarshalJSON([]byte(test.json))

			if err != nil && !errors.Is(err, test.err) {
				t.Errorf("unexpected unmarshal error: %v", err)
			} else if id != test.verify {
				t.Errorf("got '%v', want '%v'", id, test.verify)
			}
		})
	}
}
