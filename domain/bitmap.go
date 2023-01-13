package domain

import (
	"fmt"
	"math"
)

type Bitmap interface {
	Set(pos int)
	UnSet(pos int)
	IsSet(pos int) bool
	Reset()
}

type byteBitmap struct {
	size int
	data []byte
}

type InvalidPositionError error

func NewByteBitmap(size int) Bitmap {
	data := make([]byte, int(math.Ceil(float64(size)/8)))
	return &byteBitmap{size: size, data: data}
}

func (b *byteBitmap) Set(pos int) {
	b.validate(pos)
	bytePos := pos / 8
	b.data[bytePos] = b.data[bytePos] | (1 << (pos % 8))
}

func (b *byteBitmap) UnSet(pos int) {
	b.validate(pos)
	bytePos := pos / 8
	b.data[bytePos] = b.data[bytePos] &^ (1 << (pos % 8))
}

func (b *byteBitmap) IsSet(pos int) bool {
	b.validate(pos)
	bytePos := pos / 8
	return (b.data[bytePos] & (1 << (pos % 8))) != 0
}

func (b *byteBitmap) Reset() {
	for i := range b.data {
		b.data[i] = b.data[i] & 0
	}
}

func (b *byteBitmap) validate(pos int) {
	if pos < 0 {
		panic(InvalidPositionError(fmt.Errorf("pos cannot be negative")))
	}
	if pos > (b.size - 1) {
		panic(InvalidPositionError(fmt.Errorf("pos cannot be greater than size of the bitmap")))
	}
}
