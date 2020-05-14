package wlog

import (
	"unsafe"
	"time"
	"fmt"
)

const (
	nilStr = "<nil>"
)

type Field struct {
	Key string
	Val FieldVal
}

// FieldVal interface is used to serialize a value of "Field" into the given buf.
// Third-party developers can decide how to serialize a value of a any type.
type FieldVal interface {
	// Encode encodes the value of field to the given buf for the enc.
	Encode(enc ObjEncoder, buf *Buffer)
}

type BoolVal bool

func (v BoolVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendBool(buf, bool(v))
}

type ByteVal byte

func (v ByteVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendByte(buf, byte(v))
}

type IntVal int64

func (v IntVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendInt64(buf, int64(v))
}

type UintVal uint64

func (v UintVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendUint64(buf, uint64(v))
}

type FloatVal float64

func (v FloatVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendFloat64(buf, float64(v))
}

type ComplexVal complex128

func (v ComplexVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendComplex128(buf, complex128(v))
}

type StringVal string

func (v StringVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendString(buf, string(v))
}

type BytesVal []byte

func (v BytesVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendBytes(buf, v)
}

type ByteStringVal []byte

func (v ByteStringVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendByteString(buf, v)
}

type DurationVal time.Duration

func (v DurationVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendDuration(buf, time.Duration(v))
}

type TimeVal time.Time

func (v TimeVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendTime(buf, time.Time(v))
}

// Bool Returns a Field with the given key and value.
func Bool(key string, val bool) Field {
	return Field{Key: key, Val: BoolVal(val)}
}

// Uint Returns a Field with the given key and value.
func Uint(key string, val uint) Field {
	return Uint64(key, uint64(val))
}

// Uint8 Returns a Field with the given key and value.
func Uint8(key string, val uint8) Field {
	return Uint64(key, uint64(val))
}

// Uint16 Returns a Field with the given key and value.
func Uint16(key string, val uint16) Field {
	return Uint64(key, uint64(val))
}

// Uint32 Returns a Field with the given key and value.
func Uint32(key string, val uint32) Field {
	return Uint64(key, uint64(val))
}

// Uint64 Returns a Field with the given key and value.
func Uint64(key string, val uint64) Field {
	return Field{Key: key, Val: UintVal(val)}
}

// Uintptr Returns a Field with the given key and value.
func Uintptr(key string, val uintptr) Field {
	return Uint64(key, uint64(val))
}

// Int8 Returns a Field with the given key and value.
func Int(key string, val int) Field {
	return Int64(key, int64(val))
}

// Int8 Returns a Field with the given key and value.
func Int8(key string, val int8) Field {
	return Int64(key, int64(val))
}

// Int16 Returns a Field with the given key and value.
func Int16(key string, val int16) Field {
	return Int64(key, int64(val))
}

// Int32 Returns a Field with the given key and value.
func Int32(key string, val int32) Field {
	return Int64(key, int64(val))
}

// Int64 Returns a Field with the given key and value.
func Int64(key string, val int64) Field {
	return Field{Key: key, Val: IntVal(val)}
}

// Float32 Returns a Field with the given key and value.
func Float32(key string, val float32) Field {
	return Float64(key, float64(val))
}

// Float64 Returns a Field with the given key and value.
func Float64(key string, val float64) Field {
	return Field{Key: key, Val: FloatVal(val)}
}

// Complex64 Returns a Field with the given key and value.
func Complex64(key string, val complex64) Field {
	return Complex128(key, complex128(val))
}

// Complex128 Returns a Field with the given key and value.
func Complex128(key string, val complex128) Field {
	return Field{Key: key, Val: ComplexVal(val)}
}

// String Returns a Field with the given key and value.
func String(key, val string) Field {
	return Field{Key: key, Val: StringVal(val)}
}

// Bytes Returns a Field with the given key and value.
func Bytes(key string, val []byte) Field {
	return Field{Key: key, Val: BytesVal(val)}
}

// ByteString Returns a Field with the given key and value.
func ByteString(key string, val []byte) Field {
	return Field{Key: key, Val: ByteStringVal(val)}
}

// Ptr Returns a Field with the given key and value.
// // It outputs the val.String() when logging the val.
func Duration(key string, val time.Duration) Field {
	return String(key, val.String())
}

// Ptr Returns a Field with the given key and value and layout.
// "2006-01-02 15:04:05" is used as the layout when formatting the time.
func Time(key string, val time.Time) Field {
	return Field{Key: key, Val: TimeVal(val)}
}

// Err Returns a Field with the given key and err.
// If the given err is nil, then output "<nil>" when logging the err,
// otherwise output err.Error().
func Err(key string, err error) Field {
	if err == nil {
		return String(key, nilStr)
	}
	return String(key, err.Error())
}

// Object Returns a Field with the given key and value.
func Object(key string, val FieldVal) Field {
	return Field{Key: key, Val: val}
}

// BoolPtr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func BoolPtr(key string, val *bool) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Uint8Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Uint8Ptr(key string, val *uint8) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Uint16Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Uint16Ptr(key string, val *uint16) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Uint32Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Uint32Ptr(key string, val *uint32) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Uint64Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Uint64Ptr(key string, val *uint64) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Int8Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Int8Ptr(key string, val *int8) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Int16Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Int16Ptr(key string, val *int16) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Int32Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Int32Ptr(key string, val *int32) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Int64Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Int64Ptr(key string, val *int64) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Float32Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Float32Ptr(key string, val *float32) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Float64Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Float64Ptr(key string, val *float64) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Complex64Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Complex64Ptr(key string, val *complex64) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Complex128Ptr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func Complex128Ptr(key string, val *complex128) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// StringPtr Returns a Field with the given key and value.
// If the given val is nil, then output "<nil>" when logging the value,
// otherwise output the address of the val instead of the pointer value of the val.
func StringPtr(key string, val *string) Field {
	return Ptr(key, unsafe.Pointer(val))
}

// Ptr Returns a Field with the given key and value.
// It output the address of v in hexadecimal form.
// If the given val is nil, then output "<nil>" when logging the value.
func Ptr(key string, val unsafe.Pointer) Field {
	if val == nil {
		return String(key, nilStr)
	}
	v := uint64(uintptr(val))
	bs := uin64ToHex(v)
	return Field{Key: key, Val: BytesVal(bs)}
}

// Interface Returns a Field with the given key and value.
// It use the type assertion to construct a Field first,
// if necessary, it will construct a Field by fmt.s Sprint method finally.
func Interface(key string, val interface{}) Field {
	//return String(key,"==========================================")
	switch v := val.(type) {
	case bool:
		return Bool(key, v)
	case *bool:
		return BoolPtr(key, v)
	case []bool:
		return Bools(key, v)
	case uint8:
		return Uint8(key, v)
	case *uint8:
		return Uint8Ptr(key, v)
	case []uint8:
		return Uint8s(key, v)
	case int8:
		return Int8(key, v)
	case *int8:
		return Int8Ptr(key, v)
	case []int8:
		return Int8s(key, v)

	case uint16:
		return Uint16(key, v)
	case *uint16:
		return Uint16Ptr(key, v)
	case []uint16:
		return Uint16s(key, v)
	case int16:
		return Int16(key, v)
	case *int16:
		return Int16Ptr(key, v)
	case []int16:
		return Int16s(key, v)

	case uint32:
		return Uint32(key, v)
	case *uint32:
		return Uint32Ptr(key, v)
	case []uint32:
		return Uint32s(key, v)
	case int32:
		return Int32(key, v)
	case *int32:
		return Int32Ptr(key, v)
	case []int32:
		return Int32s(key, v)

	case uint64:
		return Uint64(key, v)
	case *uint64:
		return Uint64Ptr(key, v)
	case []uint64:
		return Uint64s(key, v)
	case int64:
		return Int64(key, v)
	case *int64:
		return Int64Ptr(key, v)
	case []int64:
		return Int64s(key, v)

	case float32:
		return Float32(key, v)
	case *float32:
		return Float32Ptr(key, v)
	case []float32:
		return Float32s(key, v)
	case float64:
		return Float64(key, v)
	case *float64:
		return Float64Ptr(key, v)
	case []float64:
		return Float64s(key, v)

	case complex64:
		return Complex64(key, v)
	case *complex64:
		return Complex64Ptr(key, v)
	case []complex64:
		return Complex64s(key, v)
	case complex128:
		return Complex128(key, v)
	case *complex128:
		return Complex128Ptr(key, v)
	case []complex128:
		return Complex128s(key, v)

	case string:
		return String(key, v)
	case *string:
		return StringPtr(key, v)
	case []string:
		return Strings(key, v)
	case unsafe.Pointer:
		return Ptr(key, v)
	case time.Duration:
		return Duration(key, v)
	case []time.Duration:
		return Durations(key, v)
	case time.Time:
		return Time(key, v)
	case []time.Time:
		return Times(key, v)
	case error:
		return Err(key, v)
	case []error:
		return Errs(key, v)
	}

	if fv, ok := val.(FieldVal); ok {
		return Field{Key: key, Val: fv}
	}

	return Field{Key: key, Val: StringVal(fmt.Sprint(val))}
}
