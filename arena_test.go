package gomemory

import (
	"bytes"
	"encoding/binary"
	"runtime"
	"testing"
	"unsafe"
)

// @Cleanup: Make sure to name properly tests, also if`s looks kinda ugly a the moment.
func TestMallocArenaNewObject(t *testing.T) {
	t.Parallel()

	count := 1000
	ints := make([]*uint32, count)

	arena := NewMallocArena(SizeOfAligned[uint32](count))
	defer arena.Free()

	for i := range count {
		x := New[uint32](arena)
		*x = uint32(i)
		ints[i] = x
	}

	for i := range count {
		v := ints[i]

		if v == nil {
			t.Errorf("expected non-nil value at index %d", i)
		} else if *v != uint32(i) {
			t.Errorf("expected %d got %d", i, *v)
		}
	}
}

func TestAlignedSizeOfType(t *testing.T) {
	t.Parallel()

	testAlignedSizeTimes[bool](t, 1000)
	testAlignedSizeTimes[int8](t, 1000)
	testAlignedSizeTimes[uint8](t, 1000)
	testAlignedSizeTimes[int16](t, 1000)
	testAlignedSizeTimes[uint16](t, 1000)
	testAlignedSizeTimes[int32](t, 1000)
	testAlignedSizeTimes[uint32](t, 1000)
	testAlignedSizeTimes[int64](t, 1000)
	testAlignedSizeTimes[uint64](t, 1000)
	testAlignedSizeTimes[int](t, 1000)
	testAlignedSizeTimes[uint](t, 1000)
	testAlignedSizeTimes[uintptr](t, 1000)
	testAlignedSizeTimes[float32](t, 1000)
	testAlignedSizeTimes[float64](t, 1000)
	testAlignedSizeTimes[complex64](t, 1000)
	testAlignedSizeTimes[complex128](t, 1000)
	testAlignedSizeTimes[string](t, 1000)

	testAlignedSizeTimes[struct {
		A int
		B uintptr
		C struct {
			CA string
			CB func(string) int
			CC float32
		}
		D complex64
		E unsafe.Pointer
	}](t, 1000)
}

func testAlignedSizeTimes[T any](t *testing.T, count int) {
	t.Helper()

	for n := 1; n <= count; n++ {
		alignedSize := SizeOfAligned[T](n)
		arena := NewMallocArena(SizeOfAligned[T](n))

		for range n {
			_ = New[T](arena)
		}

		if alignedSize != int(arena.size()) {
			t.Errorf("calculated aligned size is %d, but arena size %d for %d %T`s", alignedSize, int(arena.size()), n, *new(T))
		}

		buf := new(bytes.Buffer)
		if bufLen, err := arena.WriteRawMemory(buf); err != nil {
			t.Error(err)
		} else if bufLen != alignedSize {
			t.Errorf("calculated aligned size is %d, but WriteRawMemory size %d for %d %T`s", alignedSize, bufLen, n, *new(T))
		}

		if arena.Free(); arena.size() != 0 {
			t.Errorf("arena size is not 0 after Free")
		}

		buf = new(bytes.Buffer)
		if bufLen, err := arena.WriteRawMemory(buf); err != nil {
			t.Error(err)
		} else if bufLen != 0 {
			t.Errorf("Dump buffer expected to be 0, but got %d", bufLen)
		}
	}
}

func TestMallocArenaMemoryLayout(t *testing.T) {
	t.Parallel()

	arena := NewMallocArena(SizeOfAligned[uint32](2))
	defer arena.Free()

	x := New[uint32](arena)
	*x = 1

	y := New[uint32](arena)
	*y = 2

	buf := new(bytes.Buffer)

	bufLen, err := arena.WriteRawMemory(buf)
	if err != nil {
		t.Error(err)
	} else if bufLen != SizeOfAligned[uint32](2) {
		t.Errorf("expected %d dumped bytes, but arena reported size is %d", SizeOfAligned[uint32](2), bufLen)
	}

	// @Incomplete: Endians.
	var num [4]uint32
	if err := binary.Read(buf, binary.LittleEndian, &num); err != nil {
		t.Error(err)
	}

	if num[2] != *x {
		t.Errorf("expected %d in buffer, but got %d", *x, num[2])
	}

	if num[0] != *y {
		t.Errorf("expected %d in buffer, but got %d", *y, num[0])
	}
}

func TestMallocArenaFree(t *testing.T) {
	t.Parallel()

	arena := NewMallocArena(1024)
	defer arena.Free()

	type someStruct struct {
		A uint32
		S string
	}

	x := New[someStruct](arena)
	x.A = 1
	x.S = "foo"

	arena.Free()

	y := New[someStruct](arena)
	y.A = 2
	y.S = "bar"

	if x.A != y.A {
		t.Errorf("x.A expected %d got %d", y.A, x.A)
	}

	if x.S != y.S {
		t.Errorf("x.S expected %s got %s", y.S, x.S)
	}
}

// @Incomplete.
func TestNewStructPointer(t *testing.T) {
	t.Parallel()

	type A struct {
		aa bool
		ab int32
		ac []*A
	}

	type B struct {
		ba *A
	}

	arena := NewMallocArena(SizeOfAligned[B](2) + SizeOfAligned[A](2))
	defer arena.Free()

	b := New[B](arena)
	if b.ba != nil {
		t.Fail()
	}

	b.ba = New[A](arena)
	b.ba.aa = true
	b.ba.ab = 1
	b.ba.ac = make([]*A, 2)
	b.ba.ac[0] = &A{ab: 12}
	b.ba.ac[1] = &A{ab: 13}

	if b.ba.aa != true || b.ba.ab != 1 {
		t.Fail()
	}

	runtime.GC()

	if len(b.ba.ac) != 2 {
		t.Fail()
	}

	if b.ba.ac[0].ab != 12 || b.ba.ac[1].ab != 13 {
		t.Fail()
	}
}
