// Code generated by go-enum DO NOT EDIT.
// Version: 0.6.0
// Revision: 919e61c0174b91303753ee3898569a01abb32c97
// Build Date: 2023-12-18T15:54:43Z
// Built By: goreleaser

package stdlib

import (
	"fmt"
	"strings"
)

const (
	// DataTypeNull is a DataType of type null.
	DataTypeNull DataType = "null"
	// DataTypeBool is a DataType of type bool.
	DataTypeBool DataType = "bool"
	// DataTypeInt8 is a DataType of type int8.
	DataTypeInt8 DataType = "int8"
	// DataTypeInt16 is a DataType of type int16.
	DataTypeInt16 DataType = "int16"
	// DataTypeInt32 is a DataType of type int32.
	DataTypeInt32 DataType = "int32"
	// DataTypeInt64 is a DataType of type int64.
	DataTypeInt64 DataType = "int64"
	// DataTypeUint8 is a DataType of type uint8.
	DataTypeUint8 DataType = "uint8"
	// DataTypeUint16 is a DataType of type uint16.
	DataTypeUint16 DataType = "uint16"
	// DataTypeUint32 is a DataType of type uint32.
	DataTypeUint32 DataType = "uint32"
	// DataTypeUint64 is a DataType of type uint64.
	DataTypeUint64 DataType = "uint64"
	// DataTypeFloat32 is a DataType of type float32.
	DataTypeFloat32 DataType = "float32"
	// DataTypeFloat64 is a DataType of type float64.
	DataTypeFloat64 DataType = "float64"
	// DataTypeBytes is a DataType of type bytes.
	DataTypeBytes DataType = "bytes"
	// DataTypeUtf8 is a DataType of type utf8.
	DataTypeUtf8 DataType = "utf8"
	// DataTypeDate is a DataType of type date.
	DataTypeDate DataType = "date"
	// DataTypeTimestamp is a DataType of type timestamp.
	DataTypeTimestamp DataType = "timestamp"
	// DataTypeTimestampS is a DataType of type timestamp_s.
	DataTypeTimestampS DataType = "timestamp_s"
	// DataTypeTimestampMs is a DataType of type timestamp_ms.
	DataTypeTimestampMs DataType = "timestamp_ms"
	// DataTypeList is a DataType of type list.
	DataTypeList DataType = "list"
)

var ErrInvalidDataType = fmt.Errorf("not a valid DataType, try [%s]", strings.Join(_DataTypeNames, ", "))

var _DataTypeNames = []string{
	string(DataTypeNull),
	string(DataTypeBool),
	string(DataTypeInt8),
	string(DataTypeInt16),
	string(DataTypeInt32),
	string(DataTypeInt64),
	string(DataTypeUint8),
	string(DataTypeUint16),
	string(DataTypeUint32),
	string(DataTypeUint64),
	string(DataTypeFloat32),
	string(DataTypeFloat64),
	string(DataTypeBytes),
	string(DataTypeUtf8),
	string(DataTypeDate),
	string(DataTypeTimestamp),
	string(DataTypeTimestampS),
	string(DataTypeTimestampMs),
	string(DataTypeList),
}

// DataTypeNames returns a list of possible string values of DataType.
func DataTypeNames() []string {
	tmp := make([]string, len(_DataTypeNames))
	copy(tmp, _DataTypeNames)
	return tmp
}

// String implements the Stringer interface.
func (x DataType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x DataType) IsValid() bool {
	_, err := ParseDataType(string(x))
	return err == nil
}

var _DataTypeValue = map[string]DataType{
	"null":         DataTypeNull,
	"bool":         DataTypeBool,
	"int8":         DataTypeInt8,
	"int16":        DataTypeInt16,
	"int32":        DataTypeInt32,
	"int64":        DataTypeInt64,
	"uint8":        DataTypeUint8,
	"uint16":       DataTypeUint16,
	"uint32":       DataTypeUint32,
	"uint64":       DataTypeUint64,
	"float32":      DataTypeFloat32,
	"float64":      DataTypeFloat64,
	"bytes":        DataTypeBytes,
	"utf8":         DataTypeUtf8,
	"date":         DataTypeDate,
	"timestamp":    DataTypeTimestamp,
	"timestamp_s":  DataTypeTimestampS,
	"timestamp_ms": DataTypeTimestampMs,
	"list":         DataTypeList,
}

// ParseDataType attempts to convert a string to a DataType.
func ParseDataType(name string) (DataType, error) {
	if x, ok := _DataTypeValue[name]; ok {
		return x, nil
	}
	return DataType(""), fmt.Errorf("%s is %w", name, ErrInvalidDataType)
}

// MarshalText implements the text marshaller method.
func (x DataType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *DataType) UnmarshalText(text []byte) error {
	tmp, err := ParseDataType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
