package domain

import (
	"testing"
)

func TestByteBitmap_Set(t *testing.T) {
	bitmap := NewByteBitmap(10)
	bitmap.Set(8)
	isSet := bitmap.IsSet(8)
	if !isSet {
		t.Errorf("unable to set the bit")
	}
}

func TestByteBitmap_UnSet(t *testing.T) {
	bitmap := NewByteBitmap(10)
	bitmap.Set(8)
	bitmap.UnSet(8)
	isSet := bitmap.IsSet(8)
	if isSet {
		t.Errorf("unable to unset the bit")
	}
}

func TestByteBitmap_Reset(t *testing.T) {
	bitmap := NewByteBitmap(10)
	bitmap.Set(8)
	bitmap.Set(5)
	bitmap.Set(1)
	bitmap.Reset()
	if bitmap.IsSet(1) || bitmap.IsSet(5) || bitmap.IsSet(8) {
		t.Errorf("unable to reset the bitmap")
	}
}
