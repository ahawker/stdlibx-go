package stdlib

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"math"
	"math/rand"
	"os"
	"time"
)

// Rand a default global random number generator.
var (
	Rand   *rand.Rand
	Source rand.Source
)

const (
	alphabet     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits       = "0123456789"
	alphaNumeric = alphabet + digits
)

func init() {
	if seed, ok := os.LookupEnv("STDLIB_RANDOM_SEED"); ok {
		Source = rand.NewSource(int64(MustInt(seed)))
		Rand = rand.New(Source)
	} else {
		Source = rand.NewSource(time.Now().UnixNano())
		Rand = rand.New(Source)
	}
}

// RandomSelection returns a random item from the provided list of items.
func RandomSelection[T any](rng *rand.Rand, items []T) T {
	switch len(items) {
	case 0:
		return *new(T)
	case 1:
		return items[0]
	default:
		return items[rng.Intn(len(items))]
	}
}

// RandomExcluding generates a random value that is not in the excluded set.
//
// Note: Depending on the content of the excluded set, this function may be extremely inefficient.
func RandomExcluding[T comparable](fn func() T, exclude map[T]struct{}) T {
	for {
		value := fn()
		if _, ok := exclude[value]; !ok {
			return value
		}
	}
}

// RandomNumber returns a random number between primitive value min & max.
func RandomNumber[T constraints.Integer | constraints.Float](rng *rand.Rand) T {
	switch t := any(*new(T)).(type) {
	case uint8:
		return T(RandomNumberRange(rng, uint8(0), math.MaxUint8))
	case uint16:
		return T(RandomNumberRange(rng, uint16(0), math.MaxUint16))
	case uint32:
		return T(RandomNumberRange(rng, uint32(0), math.MaxUint32))
	case uint:
		return T(RandomNumberRange(rng, uint(0), math.MaxUint))
	case uint64:
		return T(RandomNumberRange(rng, uint64(0), math.MaxUint64))
	case int8:
		return T(RandomNumberRange(rng, int8(0), math.MaxInt8))
	case int16:
		return T(RandomNumberRange(rng, int16(0), math.MaxInt16))
	case int32:
		return T(RandomNumberRange(rng, int32(0), math.MaxInt32))
	case int:
		return T(RandomNumberRange(rng, int(0), math.MaxInt))
	case int64:
		return T(RandomNumberRange(rng, int64(0), math.MaxInt64))
	case float32:
		return T(RandomNumberRange(rng, float32(0), math.MaxFloat32))
	case float64:
		return T(RandomNumberRange(rng, float64(0), math.MaxFloat64))
	default:
		panic(fmt.Sprintf("RandomNumber[%T] not supported", t))
	}
}

// RandomNumberRange returns a random number between min (inclusive) and max (exclusive).
func RandomNumberRange[T constraints.Integer | constraints.Float](rng *rand.Rand, min, max T) T {
	switch any(min).(type) {
	case uint8, uint16, uint32, uint:
		return T(rng.Int31n(int32(max-min+1)) - int32(min))
	case uint64, int64:
		return T(rng.Int63n(int64(max-min+1)) - int64(min))
	case int8, int16, int32, int:
		return T(rng.Int31n(int32(max-min+1)) - int32(min))
	case float32:
		return T(float32(min) + rand.Float32()*float32(max-min))
	case float64:
		return T(float64(min) + rand.Float64()*float64(max-min))
	default:
		panic(fmt.Sprintf("RandomNumberRange[%T] not supported", min))
	}
}

// RandomString returns a random string between min length (inclusive) and max length (exclusive).
func RandomString[T ~string](rng *rand.Rand, min, max uint64) T {
	length := RandomNumberRange[uint64](rng, min, max)
	result := make([]rune, length)
	for i := uint64(0); i < length; i++ {
		result[i] = rune(alphaNumeric[rng.Intn(len(alphabet))])
	}
	return T(result)
}
