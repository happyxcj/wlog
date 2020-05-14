package wlog

import (
	"sync"
	"strconv"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &Buffer{buf: make([]byte, 1024)}
	},
}

// GetBuf returns a Buffer from the buffer pool.
func GetBuf() *Buffer {
	buf := bufferPool.Get().(*Buffer)
	buf.Reset()
	return buf
}

// AppendBuf puts the buf to the buffer pool.
func PutBuf(buf *Buffer) {
	bufferPool.Put(buf)
}

type Buffer struct {
	buf []byte
}

// Reset resets the internal byte slice for next reuse.
func (b *Buffer) Reset() {
	b.buf = b.buf[:0]
}

// Bytes returns the length of to the internal byte slice.
func (b *Buffer) Len() int {
	return len(b.buf)
}

// Bytes returns the capacity of to the internal byte slice.
func (b *Buffer) Cap() int {
	return cap(b.buf)
}

// Bytes returns the reference to the internal byte slice.
func (b *Buffer) Bytes() []byte {
	return b.buf
}

// String returns a copy of the internal byte slice.
func (b *Buffer) String() string {
	return string(b.buf)
}

// AppendBool appends a string representation of the bool to the buffer.
func (b *Buffer) AppendBool(v bool) {
	if v {
		b.AppendString("true")
	} else {
		b.AppendString("false")
	}
}

// AppendByte appends a string representation of the byte to the buffer.
func (b *Buffer) AppendByte(v byte) {
	b.buf = append(b.buf, v)
}

// AppendUint8 appends a string representation of the uint to the buffer.
func (b *Buffer) AppendUint(v uint) {
	b.AppendUint64(uint64(v))
}

// AppendUint8 appends a string representation of the uint8 to the buffer.
func (b *Buffer) AppendUint8(v uint8) {
	b.AppendUint64(uint64(v))
}

// AppendUint16 appends a string representation of the uint16 to the buffer.
func (b *Buffer) AppendUint16(v uint16) {
	b.AppendUint64(uint64(v))
}

// AppendUint32 appends a string representation of the uint32 to the buffer.
func (b *Buffer) AppendUint32(v uint32) {
	b.AppendUint64(uint64(v))
}

// AppendUint64 appends a string representation of the uint64 to the buffer.
func (b *Buffer) AppendUint64(v uint64) {
	//str := strconv.FormatUint(v, 10)
	//b.AppendString(str)
	b.buf = strconv.AppendUint(b.buf, v, 10)
}

// AppendInt8 appends a string representation of the int to the buffer.
func (b *Buffer) AppendInt(v int) {
	b.AppendInt64(int64(v))
}

// AppendInt8 appends a string representation of the int8 to the buffer.
func (b *Buffer) AppendInt8(v int8) {
	b.AppendInt64(int64(v))
}

// AppendInt16 appends a string representation of the int16 to the buffer.
func (b *Buffer) AppendInt16(v int16) {
	b.AppendInt64(int64(v))
}

// AppendInt32 appends a string representation of the int32 to the buffer.
func (b *Buffer) AppendInt32(v int32) {
	b.AppendInt64(int64(v))
}

// AppendInt64 appends a string representation of the int64 to the buffer.
func (b *Buffer) AppendInt64(v int64) {
	//str := strconv.FormatInt(v, 10)
	//b.AppendString(str)
	b.buf = strconv.AppendInt(b.buf, v, 10)
}

// AppendFloat32 appends a string representation of the float32 to the buffer.
func (b *Buffer) AppendFloat32(v float32) {
	b.AppendFloat64(float64(v))
}

// AppendFloat64 appends a string representation of the float64 to the buffer.
func (b *Buffer) AppendFloat64(v float64) {
	b.buf = strconv.AppendFloat(b.buf, v, 'f', -1, 64)
}

// AppendFloat64 appends a string representation of the complex64 to the buffer.
func (b *Buffer) AppendComplex64(v complex64) {
	b.AppendComplex128(complex128(v))
}

// AppendFloat64 appends a string representation of the complex128 to the buffer.
func (b *Buffer) AppendComplex128(v complex128) {
	val := complex128(v)
	r, i := float64(real(val)), float64(imag(val))
	b.AppendByte('(')
	b.AppendFloat64(r)
	b.AppendFloat64(i)
	b.AppendString("i)")
}

// AppendString appends a string to the buffer.
func (b *Buffer) AppendString(v string) {
	b.buf = append(b.buf, v...)
}

// AppendBytes appends a byte slice to the buffer.
func (b *Buffer) AppendBytes(v []byte) {
	b.buf = append(b.buf, v...)
}

// AppendBytes appends a byte slice to the buffer.
//func (b *Buffer) AppendTime(v int64) {
//	//b.AppendInt64(v)
//	t:=time.Unix(0,v)
//	b.buf = t.AppendFormat(b.buf,"2016-01-02 15:04:05.000")
//}
//
//func (b *Buffer) AppendTimeNew(v time.Time) {
//	b.buf = v.AppendFormat(b.buf,"2016-01-02 15:04:05.000")
//}