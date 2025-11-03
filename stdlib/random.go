package stdlib

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Global random number generator.
var (
	globalRandom *Random
	globalLock   sync.Mutex
)

const (
	alphabet     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits       = "0123456789"
	alphaNumeric = alphabet + digits
)

func init() {
	var seed int64
	if s, ok := os.LookupEnv("STDLIB_RANDOM_SEED"); ok {
		seed = int64(MustInt(s))
	} else {
		seed = time.Now().UnixNano()
	}
	globalRandom = NewRandom(seed)
}

// GetGlobal returns the global Random instance and locks it for exclusive use.
func GetGlobal() *Random {
	globalLock.Lock()
	return globalRandom
}

// ReturnGlobal returns the global Random instance and unlocks it.
func ReturnGlobal(random *Random) {
	defer globalLock.Unlock()
	if random != globalRandom {
		panic("ReturnGlobal received non-global Random instance")
	}
}

// NewRandom creates a new Random instance with the provided seed.
func NewRandom(seed int64) *Random {
	source := rand.NewSource(seed)
	return &Random{
		Rand:   rand.New(source),
		Source: source,
		Seed:   seed,
	}
}

// Random represents a random number generator with its source and seed.
//
// Callers are free to create their own and pass them into functions or
// use the 'Get' and 'Return' functions to borrow the global one.
type Random struct {
	Rand   *rand.Rand
	Source rand.Source
	Seed   int64
}

// RandomSelection returns a random item from the provided list of items.
func RandomSelection[T any](r *Random, items []T) T {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	switch len(items) {
	case 0:
		return *new(T)
	case 1:
		return items[0]
	default:
		return items[r.Rand.Intn(len(items))]
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
func RandomNumber[T constraints.Integer | constraints.Float](r *Random) T {
	switch t := any(*new(T)).(type) {
	case uint8:
		return T(RandomNumberRange(r, uint8(0), math.MaxUint8))
	case uint16:
		return T(RandomNumberRange(r, uint16(0), math.MaxUint16))
	case uint32:
		return T(RandomNumberRange(r, uint32(0), math.MaxUint32))
	case uint:
		return T(RandomNumberRange(r, uint(0), math.MaxUint))
	case uint64:
		return T(RandomNumberRange(r, uint64(0), math.MaxUint64))
	case int8:
		return T(RandomNumberRange(r, int8(0), math.MaxInt8))
	case int16:
		return T(RandomNumberRange(r, int16(0), math.MaxInt16))
	case int32:
		return T(RandomNumberRange(r, int32(0), math.MaxInt32))
	case int:
		return T(RandomNumberRange(r, int(0), math.MaxInt))
	case int64:
		return T(RandomNumberRange(r, int64(0), math.MaxInt64))
	case float32:
		return T(RandomNumberRange(r, float32(0), math.MaxFloat32))
	case float64:
		return T(RandomNumberRange(r, float64(0), math.MaxFloat64))
	default:
		panic(fmt.Sprintf("RandomNumber[%T] not supported", t))
	}
}

// RandomNumberRange returns a random number between min (inclusive) and max (exclusive).
func RandomNumberRange[T constraints.Integer | constraints.Float](r *Random, min, max T) T {
	switch any(min).(type) {
	case uint8, uint16, uint32, uint:
		return T(r.Rand.Int31n(int32(max-min+1)) - int32(min))
	case uint64, int64:
		return T(r.Rand.Int63n(int64(max-min+1)) - int64(min))
	case int8, int16, int32, int:
		return T(r.Rand.Int31n(int32(max-min+1)) - int32(min))
	case float32:
		return T(float32(min) + r.Rand.Float32()*float32(max-min))
	case float64:
		return T(float64(min) + r.Rand.Float64()*float64(max-min))
	default:
		panic(fmt.Sprintf("RandomNumberRange[%T] not supported", min))
	}
}

// RandomString returns a random string between min length (inclusive) and max length (exclusive).
func RandomString[T ~string](r *Random, min, max uint64) T {
	length := RandomNumberRange[uint64](r, min, max)
	result := make([]rune, length)
	for i := uint64(0); i < length; i++ {
		result[i] = rune(alphaNumeric[r.Rand.Intn(len(alphabet))])
	}
	return T(result)
}
