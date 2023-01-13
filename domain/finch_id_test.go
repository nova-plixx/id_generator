package domain

import (
	"reflect"
	"testing"
	"time"
)

func TestFinchIdGenerator_Generate(t *testing.T) {
	gen, _ := NewFinchIdGenerator(4521, 124, 1577836800000, &SystemClock{})
	ids, err := gen.Generate(1)
	if err != nil || ids[0] == -1 {
		t.Errorf("error while generating finch id: %v", err)
	}
}

func TestNewFinchIdGenerator_InvalidStartEpochError(t *testing.T) {
	_, err := NewFinchIdGenerator(0, 127, 4102444800000, &SystemClock{})
	checkType(t, err, reflect.TypeOf(&InvalidStartEpochError{}))
}

func TestNewFinchIdGenerator_InvalidGeneratorIdError(t *testing.T) {
	_, err := NewFinchIdGenerator(0, 524, 1577836800000, &SystemClock{})
	checkType(t, err, reflect.TypeOf(&InvalidGeneratorIdError{}))
}

func TestFinchIdGenerator_Generate_ClockShiftError(t *testing.T) {
	epoch := time.Now().UnixMilli()
	mockClock := MockClock{UnixMilli: epoch}
	gen, _ := NewFinchIdGenerator(4521, 124, 1577836800000, &mockClock)
	_, err := gen.Generate(1)
	mockClock.UnixMilli = epoch - 5
	_, err = gen.Generate(1)
	checkType(t, err, reflect.TypeOf(&ClockShiftError{}))
}

type MockClock struct {
	UnixMilli int64
}

func (m *MockClock) EpochMilli() int64 {
	return m.UnixMilli
}

func checkType(t *testing.T, obj interface{}, expectedType reflect.Type) {
	if obj != nil {
		switch actualType := reflect.TypeOf(obj); actualType {
		case expectedType:
			break
		default:
			t.Errorf("expected %v, but received %v", expectedType, actualType)
		}
	} else {
		t.Errorf("expected %v", expectedType)
	}
}
