package gomemory

import (
	"math/rand/v2"
	"strconv"
	"testing"
)

// @Cleanup: Refactor these tests.

func TestBitSetSetBit(t *testing.T) {
	t.Parallel()

	b := NewBitSet[int]()
	tests := make(map[int]bool, 32)

	for range 32 {
		r := rand.IntN(33)
		testSetBit(t, b, r)

		tests[r] = true
	}

	for i := range 32 {
		if b.Has(i) != tests[i] {
			t.Errorf("for %d expected %t, but got %t in bitset[%v]", i, tests[i], b.Has(i), strconv.FormatInt(int64(b.bits), 2))
		}
	}
}

func TestBitSetClear(t *testing.T) {
	t.Parallel()

	b := NewBitSet[uint8]()

	testSetBit(t, b, 2)
	testSetBit(t, b, 0)
	testSetBit(t, b, 3)

	testClearBit(t, b, 3)
	testClearBit(t, b, 2)
	testClearBit(t, b, 0)
}

func testSetBit[T Int](t *testing.T, b *BitSet[T], v T) {
	t.Helper()
	b.Set(v)

	if !b.Has(v) {
		t.Errorf("failed to set %d in bitset=[%s]", v, strconv.FormatInt(int64(b.bits), 2))
	}
}

func testClearBit[T Int](t *testing.T, b *BitSet[T], v T) {
	t.Helper()

	if !b.Has(v) {
		t.Errorf("failed to set %d in bitset=[%s]", v, strconv.FormatInt(int64(b.bits), 2))
	}

	b.Clear(v)

	if b.Has(v) {
		t.Errorf("bit %d was not clear in bitset=[%s]", v, strconv.FormatInt(int64(b.bits), 2))
	}
}
