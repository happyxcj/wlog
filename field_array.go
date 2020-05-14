package wlog

import "time"

type bools []bool

func (vs bools) Size() int {
	return len(vs)
}

func (vs bools) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendBool(buf, vs[i])
}

type uints []uint

func (vs uints) Size() int {
	return len(vs)
}
func (vs uints) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendUint(buf, vs[i])
}

type uint8s []uint8

func (vs uint8s) Size() int {
	return len(vs)
}
func (vs uint8s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendUint8(buf, vs[i])
}

type uint16s []uint16

func (vs uint16s) Size() int {
	return len(vs)
}
func (vs uint16s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendUint16(buf, vs[i])
}

type uint32s []uint32

func (vs uint32s) Size() int {
	return len(vs)
}
func (vs uint32s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendUint32(buf, vs[i])
}

type uint64s []uint64

func (vs uint64s) Size() int {
	return len(vs)
}
func (vs uint64s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendUint64(buf, vs[i])
}

type ints []int

func (vs ints) Size() int {
	return len(vs)
}
func (vs ints) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendInt(buf, vs[i])
}

type int8s []int8

func (vs int8s) Size() int {
	return len(vs)
}
func (vs int8s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendInt8(buf, vs[i])
}

type int16s []int16

func (vs int16s) Size() int {
	return len(vs)
}
func (vs int16s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendInt16(buf, vs[i])
}

type int32s []int32

func (vs int32s) Size() int {
	return len(vs)
}
func (vs int32s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendInt32(buf, vs[i])
}

type int64s []int64

func (vs int64s) Size() int {
	return len(vs)
}
func (vs int64s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendInt64(buf, vs[i])
}

type float32s []float32

func (vs float32s) Size() int {
	return len(vs)
}
func (vs float32s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendFloat32(buf, vs[i])
}

type float64s []float64

func (vs float64s) Size() int {
	return len(vs)
}
func (vs float64s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendFloat64(buf, vs[i])
}

type complex64s []complex64

func (vs complex64s) Size() int {
	return len(vs)
}
func (vs complex64s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendComplex64(buf, vs[i])
}

type complex128s []complex128

func (vs complex128s) Size() int {
	return len(vs)
}
func (vs complex128s) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendComplex128(buf, vs[i])
}

type strs []string

func (vs strs) Size() int {
	return len(vs)
}
func (vs strs) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendString(buf, vs[i])
}

type bytess [][]byte

func (vs bytess) Size() int {
	return len(vs)
}
func (vs bytess) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendBytes(buf, vs[i])
}

type byteStrings [][]byte

func (vs byteStrings) Size() int {
	return len(vs)
}
func (vs byteStrings) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendByteString(buf, vs[i])
}

type durations []time.Duration

func (vs durations) Size() int {
	return len(vs)
}

func (vs durations) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendDuration(buf, vs[i])
}

type times []time.Time

func (vs times) Size() int {
	return len(vs)
}

func (vs times) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	enc.AppendTime(buf, vs[i])
}

type errs []error

func (vs errs) Size() int {
	return len(vs)
}

func (vs errs) AppendEle(enc ObjEncoder, buf *Buffer, i int) {
	err := vs[i]
	if err != nil {
		enc.AppendString(buf, err.Error())
	} else {
		enc.AppendString(buf, nilStr)
	}
}

type ArrayVal struct {
	Val ArrayEncoder
}

func (v ArrayVal) Encode(enc ObjEncoder, buf *Buffer) {
	enc.AppendArray(buf, v.Val)
}

// Bools Returns a Field with the given key and values.
func Bools(key string, vs []bool) Field {
	return Array(key,bools(vs))
}

// Uints Returns a Field with the given key and values.
func Uints(key string, vs []uint) Field {
	return Array(key,uints(vs))
}

// Uint8s Returns a Field with the given key and values.
func Uint8s(key string, vs []uint8) Field {
	return Array(key,uint8s(vs))
}

// Uint16s Returns a Field with the given key and values.
func Uint16s(key string, vs []uint16) Field {
	return Array(key,uint16s(vs))
}

// Uint32s Returns a Field with the given key and values.
func Uint32s(key string, vs []uint32) Field {
	return Array(key,uint32s(vs))
}

// Uint64s Returns a Field with the given key and values.
func Uint64s(key string, vs []uint64) Field {
	return Array(key,uint64s(vs))
}

// Ints Returns a Field with the given key and values.
func Ints(key string, vs []int) Field {
	return Array(key,ints(vs))
}

// Int8s Returns a Field with the given key and values.
func Int8s(key string, vs []int8) Field {
	return Array(key,int8s(vs))
}

// Int16s Returns a Field with the given key and values.
func Int16s(key string, vs []int16) Field {
	return Array(key,int16s(vs))
}

// Int32s Returns a Field with the given key and values.
func Int32s(key string, vs []int32) Field {
	return Array(key,int32s(vs))
}

// Int64s Returns a Field with the given key and values.
func Int64s(key string, vs []int64) Field {
	return Array(key,int64s(vs))
}

// Float32s Returns a Field with the given key and values.
func Float32s(key string, vs []float32) Field {
	return Array(key,float32s(vs))
}

// Float64s Returns a Field with the given key and values.
func Float64s(key string, vs []float64) Field {
	return Array(key,float64s(vs))
}

// Complex64s Returns a Field with the given key and values.
func Complex64s(key string, vs []complex64) Field {
	return Array(key,complex64s(vs))
}

// Complex128s Returns a Field with the given key and values.
func Complex128s(key string, vs []complex128) Field {
	return Array(key,complex128s(vs))
}

// Strings Returns a Field with the given key and values.
func Strings(key string, vs []string) Field {
	return Array(key,strs(vs))
}

// Bytess Returns a Field with the given key and values.
func Bytess(key string, vs [][]byte) Field {
	return Array(key,bytess(vs))
}

// ByteStrings Returns a Field with the given key and values.
func ByteStrings(key string, vs [][]byte) Field {
	return Array(key,byteStrings(vs))
}

// Strings Returns a Field with the given key and values.
func Durations(key string, vs []time.Duration) Field {
	return Array(key,durations(vs))
}

// Strings Returns a Field with the given key and values.
func Times(key string, vs []time.Time) Field {
	return Array(key,times(vs))
}

// Strings Returns a Field with the given key and values.
func Errs(key string, vs []error) Field {
	return Array(key,errs(vs))
}

// Strings Returns a Field with the given key and values.
func Array(key string, val ArrayEncoder) Field {
	return Field{Key: key, Val: ArrayVal{Val: val}}
}
