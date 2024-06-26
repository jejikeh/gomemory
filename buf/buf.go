package buf

import (
	"errors"
	"reflect"
	"unsafe"
)

var ErrBufOverflow = errors.New("buffer overflow")

//go:linkname runtime_mallocgc runtime.mallocgc
func runtime_mallocgc(size uintptr, typ uintptr, needzero bool) unsafe.Pointer

type Buf[T any] struct {
	mem   []T
	tType uintptr
	index int
}

func New[T any](count int, ts ...T) *Buf[T] {
	var t T
	if len(ts) > 0 {
		t = ts[0]
	} else {
		t = *new(T)
	}

	ttype := runtimeTypeOf(t)
	tSize := indirectSizeof(t)
	size := tSize * uintptr(count)

	return &Buf[T]{
		mem: ptrToSlice[T](
			runtime_mallocgc(
				size,
				ttype,
				true,
			),
			int(size/tSize),
		),
		tType: ttype,
		index: 0,
	}
}

func (b *Buf[T]) Store(construct ...func(*T)) (*T, error) {
	if b.index >= len(b.mem) {
		return nil, ErrBufOverflow
	}

	t := &b.mem[b.index]
	b.index++

	for _, c := range construct {
		c(t)
	}

	return t, nil
}

func (b *Buf[T]) Length() int {
	return b.index
}

func (b *Buf[T]) Load(idx int) *T {
	if idx >= b.index {
		return nil
	}

	return &b.mem[idx]
}

func runtimeTypeOf(i any) uintptr {
	return (*[2]uintptr)(unsafe.Pointer(&i))[0]
}

func indirectSizeof(t any) uintptr {
	return reflect.Indirect(reflect.ValueOf(t)).Type().Size()
}

func ptrToSlice[T any](ptr unsafe.Pointer, count int) []T {
	var ret []T
	s := (*struct {
		ptr unsafe.Pointer
		len int
		cap int
	})(unsafe.Pointer(&ret))

	s.ptr = ptr
	s.cap = count
	s.len = s.cap

	return ret
}
