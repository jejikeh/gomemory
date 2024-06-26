package bitset

import (
	"errors"
	"unsafe"
)

type Int interface {
	~int | ~uint | ~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64
}

const bitsInByte = 4

var ErrBitSetOverflow = errors.New("bitset overflow")

type BitSet[T Int] struct {
	bits T
}

func NewBitSet[T Int](values ...T) *BitSet[T] {
	bitset := &BitSet[T]{}

	for _, value := range values {
		bitset.Set(value)
	}

	return bitset
}

func (b *BitSet[T]) Set(v T) {
	checkBitOverflow(v)
	b.bits |= 1 << v
}

func (b *BitSet[T]) Has(v T) bool {
	checkBitOverflow(v)

	return b.bits&(1<<v) != 0
}

func (b *BitSet[T]) Check(mask *BitSet[T]) bool {
	return b.bits&mask.bits == mask.bits
}

func (b *BitSet[T]) Unset(v T) {
	checkBitOverflow(v)

	b.bits &= ^(1 << v)
}

func Sizeof[T Int](v T) int {
	return int(unsafe.Sizeof(v)) * bitsInByte
}

func checkBitOverflow[T Int](v T) {
	if unsafe.Sizeof(v)*4 < uintptr(v) {
		panic(ErrBitSetOverflow)
	}
}
