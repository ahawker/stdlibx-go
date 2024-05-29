package stdlib

type Equaler[T any] interface {
	Equal(x T, y T) bool
}

type ValueObject interface{}

func VO[T ValueObject](v T) T {
	return v
}

func E[T any](v T) bool {
	switch any(v).(type) {
	case Entity:
		return true
	default:
		return false
	}
}

type ID struct {
	Type  string
	Value string
}

type Entity interface {
	EntityID() ID
}
