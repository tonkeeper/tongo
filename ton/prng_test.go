package ton

import (
	"reflect"
	"testing"
)

func createValidatorPRNG() (*ValidatorPRNG, error) {
	return NewValidatorPRNG([32]byte{}, 0x8000000000000000, -1, 0)
}

func TestNextUInt64(t *testing.T) {
	tests := []struct {
		name  string
		times int
		// expected output
		outputs []uint64
	}{
		{
			name:    "test next uint64 1 time",
			times:   1,
			outputs: []uint64{6186953295200455061},
		},
		{
			name:  "test next uint64 5 times",
			times: 5,
			outputs: []uint64{
				6186953295200455061, 9716249430906648876, 893850564141714240, 16362499097668570104,
				7550721807492789767,
			},
		},
		{
			name:  "test next uint64 10 times",
			times: 10,
			outputs: []uint64{
				6186953295200455061, 9716249430906648876, 893850564141714240, 16362499097668570104,
				7550721807492789767, 8027788155046975774, 2198044665159296191, 15889925754150310949,
				2854201576873883948, 3908958851740847745,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prng, err := createValidatorPRNG()
			if err != nil {
				t.Fatalf("cannot create validator prng: %v", err)
			}
			output := make([]uint64, test.times)
			for i := 0; i < test.times; i++ {
				output[i] = prng.NextUInt64()
			}
			if !reflect.DeepEqual(output, test.outputs) {
				t.Errorf("incorrect values: got %v, want %v", output, test.outputs)
			}
		})
	}
}

func TestNextRanged(t *testing.T) {
	tests := []struct {
		name string
		rng  uint64
		// expected output
		output uint64
	}{
		{
			name:   "test next ranged with range 5",
			rng:    5,
			output: 1,
		},
		{
			name:   "test next ranged with range 18324",
			rng:    18324,
			output: 6145,
		},
		{
			name:   "test next ranged with range 10000000000",
			rng:    10000000000,
			output: 3353954101,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prng, err := createValidatorPRNG()
			if err != nil {
				t.Fatalf("cannot create validator prng: %v", err)
			}
			output := prng.NextRanged(test.rng)
			if output != test.output {
				t.Errorf("incorrect value: got %v, want %v", output, test.output)
			}
		})
	}
}

func TestIncreaseSeed(t *testing.T) {
	seed1 := [32]byte{}
	seed1[31] = 1
	seed40 := [32]byte{}
	seed40[31] = 40
	seed256 := [32]byte{}
	seed256[30] = 1
	tests := []struct {
		name  string
		times int
		// expected seed
		seed [32]byte
	}{
		{
			name:  "test increase seed 1 time",
			times: 1,
			seed:  seed1,
		},
		{
			name:  "test increase seed 40 time",
			times: 40,
			seed:  seed40,
		},
		{
			name:  "test increase seed 256 time", // todo is it okay?
			times: 256,
			seed:  seed256,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prng, err := createValidatorPRNG()
			if err != nil {
				t.Fatalf("cannot create validator prng: %v", err)
			}
			for i := 0; i < test.times; i++ {
				prng.increaseSeed()
			}
			if !reflect.DeepEqual(prng.descr.seed, test.seed) {
				t.Errorf("incorrect seed after increases: got %v, want %v", prng.descr.seed, test.seed)
			}
		})
	}
}

func TestIncreaseHash(t *testing.T) {
	seed32 := make([]byte, 32)
	seed32[31] = 32
	hash32 := []byte{
		74, 27, 97, 202, 222, 150, 35, 10, 94, 215, 240, 213, 147, 229, 252, 235, 220, 93, 61, 153, 58, 129, 85, 207, 18, 223, 177, 238, 191, 27,
		82, 201, 215, 138, 181, 138, 211, 64, 181, 135, 235, 229, 167, 89, 39, 106, 210, 242, 97, 239, 129, 126, 111, 113, 182, 53, 72, 200, 103,
		177, 156, 208, 84, 20,
	}
	tests := []struct {
		name string
		seed []byte
		// expected hash
		hash []byte
	}{
		{
			name: "test rebuild hash",
			seed: seed32,
			hash: hash32,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prng, err := createValidatorPRNG()
			if err != nil {
				t.Fatalf("cannot create validator prng: %v", err)
			}
			copy(prng.descr.seed[:], test.seed)
			prng.rebuildHash()
			if !reflect.DeepEqual(prng.descr.hash, test.hash) {
				t.Errorf("incorrect hash after rebuild: got %v, want %v", prng.descr.hash, test.hash)
			}
		})
	}
}
