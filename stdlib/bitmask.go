package stdlib

import "strconv"

// ParseBitmask creates a new Bitmask from a binary string.
func ParseBitmask(binary string) (Bitmask, error) {
	if binary == "" {
		return Bitmask(0), nil
	}
	v, err := strconv.ParseInt(binary, 2, 32)
	if err != nil {
		return Bitmask(0), err
	}
	return Bitmask(v), nil
}

// Bitmask is a `uint8` with helper methods for bitwise operations.
type Bitmask uint8

// MarshalText implements the text marshaller method.
func (b Bitmask) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

// String returns the Bitmask in binary string (001101010) form.
func (b Bitmask) String() string {
	return strconv.FormatUint(uint64(b), 2)
}

// Clear given bits from the current mask and return a new copy.
func (b Bitmask) Clear(bits Bitmask) Bitmask {
	return b &^ bits
}

// Has checks if bits are set.
func (b Bitmask) Has(bits Bitmask) bool {
	return b&bits != 0
}

// Set bits in the current mask and return a new copy.
func (b Bitmask) Set(bits Bitmask) Bitmask {
	return b | bits
}

// Toggle bits on/off and return a new copy.
func (b Bitmask) Toggle(bits Bitmask) Bitmask {
	return b ^ bits
}
