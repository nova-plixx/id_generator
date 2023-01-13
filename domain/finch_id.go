package domain

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type IdGenerator interface {
	Generate(n int) ([]int64, error)
}

type finchIdGenerator struct {
	sharedSeed         int64
	generatorId        int
	startEpoch         int64
	randomPartBits     int
	uniquePartBits     int
	usedRandomBitmap   Bitmap
	randomPartEpoch    int64
	lastProcessedEpoch int64
	clock              Clock
	mutex              sync.Mutex
	maxGenerateLimit   int
}

type InvalidStartEpochError struct{ err error }
type InvalidGeneratorIdError struct{ err error }
type UnsupportedGenerateLimit struct{ err error }
type ClockShiftError struct{ err error }

func (e *InvalidStartEpochError) Error() string   { return e.err.Error() }
func (e *InvalidGeneratorIdError) Error() string  { return e.err.Error() }
func (e *UnsupportedGenerateLimit) Error() string { return e.err.Error() }
func (e *ClockShiftError) Error() string          { return e.err.Error() }

func NewFinchIdGenerator(sharedSeed int64, generatorId int, startEpoch int64, clock Clock) (IdGenerator, error) {
	currentEpoch := time.Now().UnixMilli()
	if currentEpoch < startEpoch {
		return nil, &InvalidStartEpochError{fmt.Errorf("startEpoch cannot be a future date")}
	}
	const randomPartBits = 12
	const uniquePartBits = 10
	var maxPossibleGeneratorId = int(math.Pow(2, float64(uniquePartBits-1))) - 1
	if generatorId > maxPossibleGeneratorId {
		return nil, &InvalidGeneratorIdError{fmt.Errorf("generatorId cannot be greater than %v", maxPossibleGeneratorId)}
	}
	usedRandomBitmap := NewByteBitmap(int(math.Ceil(math.Pow(2, float64(randomPartBits)))))
	return &finchIdGenerator{
		sharedSeed:         sharedSeed,
		generatorId:        generatorId,
		startEpoch:         startEpoch,
		usedRandomBitmap:   usedRandomBitmap,
		randomPartBits:     randomPartBits,
		uniquePartBits:     uniquePartBits,
		randomPartEpoch:    0,
		lastProcessedEpoch: 0,
		clock:              clock,
		mutex:              sync.Mutex{},
		maxGenerateLimit:   5,
	}, nil
}

func (f *finchIdGenerator) Generate(n int) ([]int64, error) {
	if n < 1 || n > f.maxGenerateLimit {
		return nil, &UnsupportedGenerateLimit{
			fmt.Errorf("the limit should be between 1 and %v (both inclusive)", f.maxGenerateLimit),
		}
	}
	f.mutex.Lock()
	defer f.mutex.Unlock()
	var res = make([]int64, n)
	for i := 0; i < n; i++ {
		id, err := f.generate()
		if err != nil {
			return nil, err
		}
		res[i] = id
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	return res, nil
}

func (f *finchIdGenerator) generate() (int64, error) {
	randomPartShift := f.uniquePartBits
	epochShift := randomPartShift + f.randomPartBits
	currentEpoch := f.clock.EpochMilli()
	if currentEpoch < f.lastProcessedEpoch {
		return -1, &ClockShiftError{fmt.Errorf("the clock shifted backwards")}
	}
	f.lastProcessedEpoch = currentEpoch
	epoch := currentEpoch - f.startEpoch
	prefix := (epoch << epochShift) | (f.perEpochRandomPart(epoch) << randomPartShift)
	uniquePart := f.generateUniquePart(prefix)
	return prefix | uniquePart, nil
}

func (f *finchIdGenerator) generateUniquePart(prefix int64) int64 {
	seed := f.sharedSeed | prefix
	arr := make([]int, int(math.Pow(2, float64(f.uniquePartBits-1))))
	for i := range arr {
		arr[i] = i
	}
	f.shuffle(arr, seed)
	return int64(arr[f.generatorId])
}

func (f *finchIdGenerator) shuffle(slice []int, seed int64) []int {
	n := len(slice)
	rng := rand.New(rand.NewSource(seed))
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = rng.Intn(n)
	}
	permutation := make([]int, n)
	for i, j := range indices {
		permutation[i] = slice[j]
	}
	return permutation
}

func (f *finchIdGenerator) perEpochRandomPart(epoch int64) int64 {
	if epoch > f.randomPartEpoch {
		f.randomPartEpoch = epoch
		f.usedRandomBitmap.Reset()
	}
	for i := 0; i < 64; i++ {
		randomPart := int64(rand.Intn(int(math.Pow(2, float64(f.randomPartBits-1)))))
		if f.usedRandomBitmap.IsSet(int(randomPart)) {
			continue
		}
		f.usedRandomBitmap.Set(int(randomPart))
		return randomPart
	}
	return f.perEpochRandomPart(epoch + 1)
}
