package tlb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/boc"
)

type arrayItemWithRef struct {
	Mode Uint8
	Cell boc.Cell `tlb:"^"`
}

func TestTolkArrayRoundTrip(t *testing.T) {
	t.Run("items with refs", func(t *testing.T) {
		for _, count := range []int{0, 1, 3, 4, 5, 9, MaxTolkArrayLen} {
			items := make([]arrayItemWithRef, 0, count)
			for i := 0; i < count; i++ {
				cell := boc.NewCell()
				require.NoError(t, cell.WriteUint(uint64(i), 32))
				items = append(items, arrayItemWithRef{Mode: Uint8(i), Cell: *cell})
			}
			decoded := marshalAndBack(t, items)
			require.Len(t, decoded, count, "count=%v", count)
			for i := range decoded {
				assert.Equal(t, items[i].Mode, decoded[i].Mode, "count=%v: item %v", count, i)
				want, err := items[i].Cell.HashString()
				require.NoError(t, err)
				got, err := decoded[i].Cell.HashString()
				require.NoError(t, err)
				assert.Equal(t, want, got, "count=%v: item %v holds another ref", count, i)
			}
		}
	})
	// an item without refs makes chunks bounded by bits rather than by refs
	t.Run("inline items", func(t *testing.T) {
		for _, count := range []int{1, 31, 32, 100, MaxTolkArrayLen} {
			items := make([]Uint32, 0, count)
			for i := 0; i < count; i++ {
				items = append(items, Uint32(i))
			}
			decoded := marshalAndBack(t, items)
			require.Len(t, decoded, count, "count=%v", count)
			assert.Equal(t, items, decoded, "count=%v", count)
		}
	})
}

// marshalAndBack writes an array, checks that the chunks look the way the readers
// of a Tolk array expect them to, and reads it back.
func marshalAndBack[T any](t *testing.T, items []T) []T {
	t.Helper()
	c := boc.NewCell()
	require.NoError(t, StoreTolkArray(c, &Encoder{}, items))

	checked := *c
	length, err := checked.ReadUint(8)
	require.NoError(t, err)
	require.EqualValues(t, len(items), length, "the array declares another length")
	var head Maybe[Ref[boc.Cell]]
	require.NoError(t, head.UnmarshalTLB(&checked, NewDecoder()))
	require.Equal(t, len(items) > 0, head.Exists, "a head chunk is there only for a non-empty array")
	for head.Exists {
		chunk := head.Value.Value
		require.NoError(t, head.UnmarshalTLB(&chunk, NewDecoder()))
		require.False(t, chunk.IsEmpty(), "a chunk holds no items")
	}

	decoded, err := LoadTolkArray[T](c, NewDecoder())
	require.NoError(t, err)
	return decoded
}

func TestTolkArrayTooLong(t *testing.T) {
	err := StoreTolkArray(boc.NewCell(), &Encoder{}, make([]Uint32, MaxTolkArrayLen+1))
	require.Error(t, err, "an array longer than MaxTolkArrayLen has no length byte to fit in")
}
