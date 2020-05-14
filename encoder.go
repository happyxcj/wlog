package wlog

import (
	"time"
)

const (
	defaultLineEnding = "\n"
)

// Encoder interface is used to serialize a log message into the structured data.
// Third-party developers can decide how to design it to make it faster, more efficient, more humane, etc.
type Encoder interface {
	ObjEncoder
	// Encode encodes the given entry and fields to the buf.
	Encode(buf *Buffer, entry *Entry, fields ...Field) error
}

type ObjEncoder interface {
	AppendBool(buf *Buffer, val bool)
	AppendByte(buf *Buffer, val byte)

	AppendUint(buf *Buffer, val uint)
	AppendUint8(buf *Buffer, val uint8)
	AppendUint16(buf *Buffer, val uint16)
	AppendUint32(buf *Buffer, val uint32)
	AppendUint64(buf *Buffer, val uint64)

	AppendInt(buf *Buffer, val int)
	AppendInt8(buf *Buffer, val int8)
	AppendInt16(buf *Buffer, val int16)
	AppendInt32(buf *Buffer, val int32)
	AppendInt64(buf *Buffer, val int64)

	AppendFloat32(buf *Buffer, val float32)
	AppendFloat64(buf *Buffer, val float64)

	AppendComplex64(buf *Buffer, val complex64)
	AppendComplex128(buf *Buffer, val complex128)

	AppendString(buf *Buffer, val string)
	AppendBytes(buf *Buffer, val []byte)
	// AppendByteString is used to append the bytes encoded in UTF-8
	AppendByteString(buf *Buffer, val []byte)

	AppendDuration(buf *Buffer, val time.Duration)

	AppendTime(buf *Buffer, val time.Time)

	AppendArray(buf *Buffer, val ArrayEncoder)

	AppendObject(buf *Buffer, val FieldVal)
}

type ArrayEncoder interface {
	// Size return the size of the array.
	Size() int
	// AppendEle encodes the element under the specified i to the buf.
	AppendEle(enc ObjEncoder, buf *Buffer, i int)
}

type TimeEncoder interface {
	Append(buf *Buffer, enc ObjEncoder, t time.Time)
}

// FastTextTimeEncoder is a faster way to encode a time.Time than the LayoutTimeEncoder
// when the LayoutTimeEncoder's layout is "2006-01-02 15:04:05" or "2006-01-02 15:04:05.000“.
type FastTextTimeEncoder struct {
	MillEnabled bool
}

func (e *FastTextTimeEncoder) Append(buf *Buffer, enc ObjEncoder, t time.Time) {
	appendFastTime(buf, t, e.MillEnabled)
}

// FastTextTimeEncoder is a faster way to encode a time.Time than the LayoutTimeEncoder
// when the LayoutTimeEncoder's layout is "2006-01-02 15:04:05" or "2006-01-02 15:04:05.000“.
type FastJsonTimeEncoder struct {
	MillEnabled bool
}

func (e *FastJsonTimeEncoder) Append(buf *Buffer, enc ObjEncoder, t time.Time) {
	buf.AppendByte('"')
	appendFastTime(buf, t, e.MillEnabled)
	buf.AppendByte('"')
}

func appendFastTime(buf *Buffer, t time.Time, milliEnabled bool) {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	ms := t.Nanosecond() / int(time.Millisecond)
	buf.AppendInt(year)
	buf.AppendByte('-')
	if month < 10 {
		buf.AppendByte('0')
	}
	buf.AppendInt(int(month))
	buf.AppendByte('-')
	if day < 10 {
		buf.AppendByte('0')
	}
	buf.AppendInt(day)
	buf.AppendByte(' ')
	if hour < 10 {
		buf.AppendByte('0')
	}
	buf.AppendInt(hour)
	buf.AppendByte(':')
	if minute < 10 {
		buf.AppendByte('0')
	}
	buf.AppendInt(minute)
	buf.AppendByte(':')
	if second < 10 {
		buf.AppendByte('0')
	}
	buf.AppendInt(second)
	if !milliEnabled {
		return
	}
	buf.AppendByte('.')
	if ms < 10 {
		buf.AppendByte('0')
		buf.AppendByte('0')
	} else if ms < 100 {
		buf.AppendByte('0')
	}
	buf.AppendInt(ms)
}

type LayoutTimeEncoder struct {
	Layout string
}

func (e LayoutTimeEncoder) AppendTime(buf *Buffer, enc ObjEncoder, t time.Time) {
	enc.AppendString(buf, t.Format(e.Layout))
}

type DurationEncoder interface {
	Append(buf *Buffer, enc ObjEncoder, dur time.Duration)
}

type StrDurationEncoder struct{}

func (e *StrDurationEncoder) Append(buf *Buffer, enc ObjEncoder, dur time.Duration) {
	enc.AppendString(buf, dur.String())
}

type EncoderOpts struct {
	// colorEnabled is a bool that indicate whether enable the color when encoding a log.
	colorEnabled bool
	// levelLower is a bool that indicate whether outAppend the lowercase of level string when encoding a log.
	levelLower bool
	// colorEnabled is a bool that indicate whether out the created time of the log.
	timeDisabled bool
	// lineEnding is the line ending of every log.
	lineEnding string

	timeEncoder     TimeEncoder

	durationEncoder DurationEncoder
}

type EncoderOpt func(opts *EncoderOpts)

func (o *EncoderOpts) WithOpts(opts ...EncoderOpt) {
	for _, opt := range opts {
		opt(o)
	}
}

// EnableColor enable the color when encoding a log.
func EnableColor() EncoderOpt {
	return func(opts *EncoderOpts) {
		opts.colorEnabled = true
	}
}

// SetLevelLower outAppends the lowercase of level string when encoding a log.
func SetLevelLower() EncoderOpt {
	return func(opts *EncoderOpts) {
		opts.levelLower = true
	}
}

// DisableTime disable the time field when encoding a log.
func DisableTime() EncoderOpt {
	return func(opts *EncoderOpts) {
		opts.timeDisabled = true
	}
}

// SetTimeEncoder sets the time encoder to encode a time.Time.
func SetTimeEncoder(e TimeEncoder) EncoderOpt {
	return func(opts *EncoderOpts) {
		opts.timeEncoder = e
	}
}

// SetDurationEncoder sets the time encoder to encode a time.Duration.
func SetDurationEncoder(e DurationEncoder) EncoderOpt {
	return func(opts *EncoderOpts) {
		opts.durationEncoder = e
	}
}

// SetLineEnding sets the line ending when encoding a log.
func SetLineEnding(ending string) EncoderOpt {
	return func(opts *EncoderOpts) {
		if ending == "" {
			opts.lineEnding = defaultLineEnding
			return
		}
		opts.lineEnding = ending
	}
}

type BasicObjEncoder struct {
}

func (e *BasicObjEncoder) AppendBool(buf *Buffer, val bool) {
	buf.AppendBool(val)
}

func (e *BasicObjEncoder) AppendByte(buf *Buffer, val byte) {
	buf.AppendByte(val)
}

func (e *BasicObjEncoder) AppendUint(buf *Buffer, val uint) {
	buf.AppendUint(val)
}

func (e *BasicObjEncoder) AppendUint8(buf *Buffer, val uint8) {
	buf.AppendUint8(val)
}

func (e *BasicObjEncoder) AppendUint16(buf *Buffer, val uint16) {
	buf.AppendUint16(val)
}

func (e *BasicObjEncoder) AppendUint32(buf *Buffer, val uint32) {
	buf.AppendUint32(val)
}

func (e *BasicObjEncoder) AppendUint64(buf *Buffer, val uint64) {
	buf.AppendUint64(val)
}

func (e *BasicObjEncoder) AppendInt(buf *Buffer, val int) {
	buf.AppendInt(val)
}

func (e *BasicObjEncoder) AppendInt8(buf *Buffer, val int8) {
	buf.AppendInt8(val)
}

func (e *BasicObjEncoder) AppendInt16(buf *Buffer, val int16) {
	buf.AppendInt16(val)
}

func (e *BasicObjEncoder) AppendInt32(buf *Buffer, val int32) {
	buf.AppendInt32(val)
}

func (e *BasicObjEncoder) AppendInt64(buf *Buffer, val int64) {
	buf.AppendInt64(val)
}

func (e *BasicObjEncoder) AppendFloat32(buf *Buffer, val float32) {
	buf.AppendFloat32(val)
}

func (e *BasicObjEncoder) AppendFloat64(buf *Buffer, val float64) {
	buf.AppendFloat64(val)
}

func (e *BasicObjEncoder) AppendComplex64(buf *Buffer, val complex64) {
	buf.AppendComplex64(val)
}

func (e *BasicObjEncoder) AppendComplex128(buf *Buffer, val complex128) {
	buf.AppendComplex128(val)
}

func (e *BasicObjEncoder) AppendString(buf *Buffer, val string) {
	buf.AppendString(val)
}

func (e *BasicObjEncoder) AppendBytes(buf *Buffer, val []byte) {
	buf.AppendBytes(val)
}

func (e *BasicObjEncoder) AppendByteString(buf *Buffer, val []byte) {
	buf.AppendBytes(val)
}


