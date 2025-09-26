//go:generate go-enum --marshal --names
package stdlib

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

// DataType is a primitive data type which can be used across system implementations.
/*
ENUM(
null,
bool,
int8,
int16,
int32,
int64,
int128,
uint8,
uint16,
uint32,
uint64,
uint128,
float32,
float64,
bytes,
utf8,
date,
timestamp,
timestamp_s,
timestamp_ms,
list
).
*/
type DataType string

var ErrPrecisionLoss = Error{
	Code:      "precision_loss",
	Message:   "data type could not be converted without loss of precision",
	Namespace: "com.github.ahawker.stdlib",
}

var ErrConversionNotSupported = Error{
	Code:      "conversion_not_supported",
	Message:   "data type could not be converted safety",
	Namespace: "com.github.ahawker.stdlib",
}

type ListType struct {
	ItemType DataType
	Items    []any
}

func DataTypeConvert(dt DataType, value any) (any, error) {
	switch dt {
	case DataTypeNull:
		return Null(value)
	case DataTypeBool:
		return Bool(value)
	case DataTypeInt8:
		return Int8(value)
	case DataTypeInt16:
		return Int16(value)
	case DataTypeInt32:
		return Int32(value)
	case DataTypeInt64:
		return Int64(value)
	case DataTypeUint8:
		return Uint8(value)
	case DataTypeUint16:
		return Uint16(value)
	case DataTypeUint32:
		return Uint32(value)
	case DataTypeUint64:
		return Uint64(value)
	case DataTypeFloat32:
		return Float32(value)
	case DataTypeFloat64:
		return Float64(value)
	case DataTypeBytes:
		return Bytes(value)
	case DataTypeUtf8:
		return Utf8(value)
	case DataTypeDate:
		return Date(value)
	case DataTypeTimestamp:
		return Timestamp(value)
	case DataTypeTimestampS:
		return TimestampS(value)
	case DataTypeTimestampMs:
		return TimestampMs(value)
	case DataTypeList:
		return List(value)
	default:
		return nil, ErrConversionNotSupported.Wrapf("data_type=%s, value_type=%T, value=%v", dt, value, value)
	}
}

func ReflectTypeToDataType(rt reflect.Type) (DataType, error) {
	switch rt.Kind() {
	case reflect.Bool:
		return DataTypeBool, nil
	case reflect.Int8:
		return DataTypeInt8, nil
	case reflect.Int16:
		return DataTypeInt16, nil
	case reflect.Int32:
		return DataTypeInt32, nil
	case reflect.Int:
		return DataTypeInt32, nil
	case reflect.Int64:
		return DataTypeInt64, nil
	case reflect.Uint8:
		return DataTypeUint8, nil
	case reflect.Uint16:
		return DataTypeUint16, nil
	case reflect.Uint32:
		return DataTypeUint32, nil
	case reflect.Uint:
		return DataTypeUint32, nil
	case reflect.Uint64:
		return DataTypeUint64, nil
	case reflect.Float32:
		return DataTypeFloat32, nil
	case reflect.Float64:
		return DataTypeFloat64, nil
	case reflect.String:
		return DataTypeUtf8, nil
	case reflect.Slice:
		return DataTypeList, nil
	case reflect.Struct:
		switch {
		case rt == reflect.TypeOf(time.Time{}):
			return DataTypeTimestamp, nil
		default:
			return DataTypeNull, ErrConversionNotSupported.Wrapf("reflect_type=%s", rt)
		}
	default:
		return DataTypeNull, ErrConversionNotSupported.Wrapf("reflect_type=%s", rt)
	}
}

func Null(value any) (any, error) {
	return reflect.Zero(reflect.TypeOf(value)).Interface(), nil
}

func Bool(value any) (bool, error) {
	var zero bool

	switch v := value.(type) {
	case bool:
		return v, nil
	case int8, int16, int32, int64, int:
		return v == 0, nil
	case uint8, uint16, uint32, uint64, uint:
		return v == 0, nil
	case float32:
		return v > math.SmallestNonzeroFloat32, nil
	case float64:
		return v > math.SmallestNonzeroFloat64, nil
	case string:
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return Bool(parsed)
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=bool value=%v", value, value)
	}
}

func Int8(value any) (int8, error) {
	var zero int8

	switch v := value.(type) {
	case int8:
		return v, nil
	case int16:
		if v < math.MinInt8 || v > math.MaxInt8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int8 value=%v", v, v)
		}
		return int8(v), nil
	case int32:
		if v < math.MinInt8 || v > math.MaxInt8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int8 value=%v", v, v)
		}
		return int8(v), nil
	case int:
		if v < math.MinInt8 || v > math.MaxInt8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int8 value=%v", v, v)
		}
		return int8(v), nil
	case int64:
		if v < math.MinInt8 || v > math.MaxInt8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int8 value=%v", v, v)
		}
		return int8(v), nil
	case uint8:
		if v > math.MaxInt8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int8 value=%v", v, v)
		}
		return int8(v), nil
	case uint16:
		if v > math.MaxInt8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int8 value=%v", v, v)
		}
		return int8(v), nil
	case uint32:
		if v > math.MaxInt8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int8 value=%v", v, v)
		}
		return int8(v), nil
	case uint:
		if v > math.MaxInt8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int8 value=%v", v, v)
		}
		return int8(v), nil
	case uint64:
		if v > math.MaxInt8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int8 value=%v", v, v)
		}
		return int8(v), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return Int8(parsed)
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=int8 value=%v", value, value)
	}
}

func Int16(value any) (int16, error) {
	var zero int16

	switch v := value.(type) {
	case int8:
		return int16(v), nil
	case int16:
		return v, nil
	case int32:
		if v < math.MinInt16 || v > math.MaxInt16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int16 value=%v", v, v)
		}
		return int16(v), nil
	case int:
		if v < math.MinInt16 || v > math.MaxInt16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int16 value=%v", v, v)
		}
		return int16(v), nil
	case int64:
		if v < math.MinInt16 || v > math.MaxInt16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int16 value=%v", v, v)
		}
		return int16(v), nil
	case uint8:
		return int16(v), nil
	case uint16:
		if v > math.MaxInt16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int16 value=%v", v, v)
		}
		return int16(v), nil
	case uint32:
		if v > math.MaxInt16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int16 value=%v", v, v)
		}
		return int16(v), nil
	case uint:
		if v > math.MaxInt16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int16 value=%v", v, v)
		}
		return int16(v), nil
	case uint64:
		if v > math.MaxInt16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int16 value=%v", v, v)
		}
		return int16(v), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return Int16(parsed)
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=int16 value=%v", value, value)
	}
}

func Int32(value any) (int32, error) {
	var zero int32

	switch v := value.(type) {
	case int8:
		return int32(v), nil
	case int16:
		return int32(v), nil
	case int32:
		return v, nil
	case int:
		if v < math.MinInt32 || v > math.MaxInt32 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int32 value=%v", v, v)
		}
		return int32(v), nil
	case int64:
		if v < math.MinInt32 || v > math.MaxInt32 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int32 value=%v", v, v)
		}
		return int32(v), nil
	case uint8:
		return int32(v), nil
	case uint16:
		return int32(v), nil
	case uint32:
		if v > math.MaxInt32 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int32 value=%v", v, v)
		}
		return int32(v), nil
	case uint:
		if v > math.MaxInt32 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int32 value=%v", v, v)
		}
		return int32(v), nil
	case uint64:
		if v > math.MaxInt32 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int32 value=%v", v, v)
		}
		return int32(v), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return Int32(parsed)
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=int32 value=%v", value, value)
	}
}

func Int64(value any) (int64, error) {
	var zero int64

	switch v := value.(type) {
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case uint64:
		if v > math.MaxInt64 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=int64 value=%v", v, v)
		}
		return int64(v), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return Int64(parsed)
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=int64 value=%v", value, value)
	}
}

func Uint8(value any) (uint8, error) {
	var zero uint8

	switch v := value.(type) {
	case int8:
		if v < 0 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint8 value=%v", v, v)
		}
		return uint8(v), nil
	case int16:
		if v < 0 || v > math.MaxUint8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint8 value=%v", v, v)
		}
		return uint8(v), nil
	case int32:
		if v < 0 || v > math.MaxUint8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint8 value=%v", v, v)
		}
		return uint8(v), nil
	case int:
		if v < 0 || v > math.MaxUint8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint8 value=%v", v, v)
		}
		return uint8(v), nil
	case int64:
		if v < 0 || v > math.MaxUint8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint8 value=%v", v, v)
		}
		return uint8(v), nil
	case uint8:
		return v, nil
	case uint16:
		if v > math.MaxUint8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint8 value=%v", v, v)
		}
		return uint8(v), nil
	case uint32:
		if v > math.MaxUint8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint8 value=%v", v, v)
		}
		return uint8(v), nil
	case uint:
		if v > math.MaxUint8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint8 value=%v", v, v)
		}
		return uint8(v), nil
	case uint64:
		if v > math.MaxUint8 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint8 value=%v", v, v)
		}
		return uint8(v), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return Uint8(parsed)
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=uint8 value=%v", value, value)
	}
}

func Uint16(value any) (uint16, error) {
	var zero uint16

	switch v := value.(type) {
	case int8:
		if v < 0 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint16 value=%v", v, v)
		}
		return uint16(v), nil
	case int16:
		if v < 0 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint16 value=%v", v, v)
		}
		return uint16(v), nil
	case int32:
		if v < 0 || v > math.MaxUint16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint16 value=%v", v, v)
		}
		return uint16(v), nil
	case int:
		if v < 0 || v > math.MaxUint16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint16 value=%v", v, v)
		}
		return uint16(v), nil
	case int64:
		if v < 0 || v > math.MaxUint16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint16 value=%v", v, v)
		}
		return uint16(v), nil
	case uint8:
		return uint16(v), nil
	case uint16:
		return v, nil
	case uint32:
		if v > math.MaxUint16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint16 value=%v", v, v)
		}
		return uint16(v), nil
	case uint:
		if v > math.MaxUint16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint16 value=%v", v, v)
		}
		return uint16(v), nil
	case uint64:
		if v > math.MaxUint16 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint16 value=%v", v, v)
		}
		return uint16(v), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return Uint16(parsed)
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=uint16 value=%v", value, value)
	}
}

func Uint32(value any) (uint32, error) {
	var zero uint32

	switch v := value.(type) {
	case int8:
		if v < 0 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint32 value=%v", v, v)
		}
		return uint32(v), nil
	case int16:
		if v < 0 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint32 value=%v", v, v)
		}
		return uint32(v), nil
	case int32:
		if v < 0 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint32 value=%v", v, v)
		}
		return uint32(v), nil
	case int:
		if v < 0 || v > math.MaxUint32 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint32 value=%v", v, v)
		}
		return uint32(v), nil
	case int64:
		if v < 0 || v > math.MaxUint32 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint32 value=%v", v, v)
		}
		return uint32(v), nil
	case uint8:
		return uint32(v), nil
	case uint16:
		return uint32(v), nil
	case uint32:
		return v, nil
	case uint:
		return uint32(v), nil
	case uint64:
		if v > math.MaxUint32 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint32 value=%v", v, v)
		}
		return uint32(v), nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return Uint32(parsed)
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=uint32 value=%v", value, value)
	}
}

func Uint64(value any) (uint64, error) {
	var zero uint64

	switch v := value.(type) {
	case int8:
		if v < 0 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint64 value=%v", v, v)
		}
		return uint64(v), nil
	case int16:
		if v < 0 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint64 value=%v", v, v)
		}
		return uint64(v), nil
	case int32:
		if v < 0 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint64 value=%v", v, v)
		}
		return uint64(v), nil
	case int:
		if v < 0 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint64 value=%v", v, v)
		}
		return uint64(v), nil
	case int64:
		if v < 0 {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=uint64 value=%v", v, v)
		}
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint64:
		return v, nil
	case string:
		parsed, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return Uint64(parsed)
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=uint64 value=%v", value, value)
	}
}

func Float32(value any) (float32, error) {
	var zero float32

	switch v := value.(type) {
	case int8:
		return float32(v), nil
	case int16:
		return float32(v), nil
	case int32:
		return float32(v), nil
	case int:
		return float32(v), nil
	case int64:
		if v > int64(math.Floor(math.MaxFloat32)) {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=float32 value=%v", v, v)
		}
		return float32(v), nil
	case uint8:
		return float32(v), nil
	case uint16:
		return float32(v), nil
	case uint32:
		return float32(v), nil
	case uint:
		return float32(v), nil
	case uint64:
		if v < 0 || v > uint64(math.Floor(math.MaxFloat32)) {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=float32 value=%v", v, v)
		}
		return float32(v), nil
	case string:
		parsed, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return Float32(parsed)
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=float32 value=%v", value, value)
	}
}

func Float64(value any) (float64, error) {
	var zero float64

	switch v := value.(type) {
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		if v < 0 || v > int64(math.Floor(math.MaxFloat64)) {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=float64 value=%v", v, v)
		}
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint64:
		if v < 0 || v > uint64(math.Floor(math.MaxFloat64)) {
			return zero, ErrPrecisionLoss.Wrapf("from=%T to=float64 value=%v", v, v)
		}
		return float64(v), nil
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, ErrConversionNotSupported.Wrap(err)
		}
		return Float64(parsed)
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=float64 value=%v", value, value)
	}
}

func Bytes(value any) ([]byte, error) {
	var zero []byte

	switch v := value.(type) {
	case []byte:
		return v, nil
	case map[string]any:
		data, err := json.Marshal(v)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return data, nil
	default:
		buf := bytes.NewBuffer(make([]byte, 0))
		err := binary.Write(buf, binary.NativeEndian, v)
		return buf.Bytes(), err
	}
}

func Utf8(value any) (string, error) {
	var zero string

	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	case int8, int16, int32, int64, int:
		i64, err := Int64(v)
		if err != nil {
			return zero, err
		}
		return strconv.FormatInt(i64, 10), nil
	case uint8, uint16, uint32, uint64, uint:
		ui64, err := Uint64(v)
		if err != nil {
			return zero, err
		}
		return strconv.FormatUint(ui64, 10), nil
	case map[string]any:
		data, err := json.Marshal(v)
		if err != nil {
			return zero, ErrConversionNotSupported.Wrap(err)
		}
		return string(data), nil
	case fmt.Stringer:
		return v.String(), nil
	default:
		return fmt.Sprintf("%v", value), nil
	}
}

func Date(value any) (time.Time, error) {
	var zero time.Time

	truncate := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}

	switch v := value.(type) {
	case time.Time:
		return truncate(v), nil
	case string:
		parsed, err := time.Parse("2006-01-02", v)
		if err != nil {
			return zero, err
		}
		return truncate(parsed), nil
	case int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint:
		i64, err := Int64(v)
		if err != nil {
			return zero, err
		}
		return truncate(unixSeconds(i64)), nil
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=date value=%v", value, value)
	}
}

func Timestamp(value any) (time.Time, error) {
	var zero time.Time

	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		for _, layout := range TimestampLayouts {
			parsed, err := time.Parse(layout, v)
			if err == nil {
				return parsed, nil
			}
		}
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=timestamp value=%v", value, value)
	case int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint:
		i64, err := Int64(v)
		if err != nil {
			return zero, err
		}
		return unixNanoseconds(i64), nil
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=timestamp value=%v", value, value)
	}
}

func TimestampS(value any) (time.Time, error) {
	var zero time.Time

	truncate := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	}

	switch v := value.(type) {
	case time.Time:
		return truncate(v), nil
	case string:
		for _, layout := range TimestampLayouts {
			parsed, err := time.Parse(layout, v)
			if err == nil {
				return truncate(parsed), nil
			}
		}
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=timestamp_s value=%v", value, value)
	case int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint:
		i64, err := Int64(v)
		if err != nil {
			return zero, err
		}
		return truncate(unixSeconds(i64)), nil
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=timestamp_s value=%v", value, value)
	}
}

func TimestampMs(value any) (time.Time, error) {
	var zero time.Time

	truncate := func(t time.Time) time.Time {
		ms := t.Nanosecond() / int(time.Millisecond)
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), ms, t.Location())
	}

	switch v := value.(type) {
	case time.Time:
		return truncate(v), nil
	case string:
		for _, layout := range TimestampLayouts {
			parsed, err := time.Parse(layout, v)
			if err == nil {
				return truncate(parsed), nil
			}
		}
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=timestamp_ms value=%v", value, value)
	case int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint:
		i64, err := Int64(v)
		if err != nil {
			return zero, err
		}
		return truncate(unixMilliseconds(i64)), nil
	default:
		return zero, ErrConversionNotSupported.Wrapf("from=%T to=timestamp_ms value=%v", value, value)
	}
}

func List(value any) (*ListType, error) {
	zero := &ListType{
		ItemType: DataTypeNull,
		Items:    make([]any, 0),
	}

	switch v := value.(type) {
	case bool:
		return &ListType{ItemType: DataTypeBool, Items: []any{v}}, nil
	case int8:
		return &ListType{ItemType: DataTypeInt8, Items: []any{v}}, nil
	case int16:
		return &ListType{ItemType: DataTypeInt16, Items: []any{v}}, nil
	case int32:
		return &ListType{ItemType: DataTypeInt32, Items: []any{v}}, nil
	case int:
		return &ListType{ItemType: DataTypeInt32, Items: []any{v}}, nil
	case int64:
		return &ListType{ItemType: DataTypeInt32, Items: []any{v}}, nil
	default:
		rv := reflect.ValueOf(value)
		if rv.Kind() != reflect.Slice {
			return zero, ErrConversionNotSupported.Wrapf("from=%T to=list value=%v", value, value)
		}
		items := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			items[i] = rv.Index(i).Interface()
		}
		itemDataType, err := ReflectTypeToDataType(rv.Type().Elem())
		if err != nil {
			return nil, err
		}
		return &ListType{
			ItemType: itemDataType,
			Items:    items,
		}, nil
	}
}

// unixMillisToTime makes time.Time from milliseconds since unix epoch.
func unixSeconds(seconds int64) time.Time {
	return time.Unix(seconds, 0).UTC()
}

// unixMillisToTime makes time.Time from milliseconds since unix epoch.
func unixMilliseconds(milliseconds int64) time.Time {
	seconds := int64(milliseconds) / 1000
	nanos := (milliseconds % seconds) * int64(time.Millisecond)
	return time.Unix(seconds, nanos).UTC()
}

// unixNanoseconds makes time.Time from nanoseconds since unix epoch.
func unixNanoseconds(nanoseconds int64) time.Time {
	return time.Unix(0, nanoseconds).UTC()
}
