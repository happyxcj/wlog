package wlog

import (
	"time"
)

const (
	KeyLevel = "level"
	KeyTime  = "time"
	KeyMsg   = "msg"
)

type JsonEncoder struct {
	BasicObjEncoder
	EncoderOpts
}

func NewJsonEncoder(opts ...EncoderOpt) *JsonEncoder {
	e := &JsonEncoder{}
	e.durationEncoder = &StrDurationEncoder{}
	e.timeEncoder = &FastJsonTimeEncoder{}
	e.levelLower = true
	e.lineEnding = defaultLineEnding
	e.WithOpts(opts...)
	return e
}

func (e *JsonEncoder) AppendString(buf *Buffer, val string) {
	buf.AppendByte('"')
	buf.AppendString(val)
	buf.AppendByte('"')
}

func (e *JsonEncoder) AppendByeString(buf *Buffer, val []byte) {
	buf.AppendByte('"')
	buf.AppendBytes(val)
	buf.AppendByte('"')
}

func (e *JsonEncoder) AppendDuration(buf *Buffer, val time.Duration) {
	e.durationEncoder.Append(buf, e, val)
}

func (e *JsonEncoder) AppendTime(buf *Buffer, val time.Time) {
	e.timeEncoder.Append(buf, e, val)
}

func (e *JsonEncoder) AppendArray(buf *Buffer, val ArrayEncoder) {
	buf.AppendByte('[')
	for i, size := 0, val.Size(); i < size; i++ {
		if i > 0 {
			buf.AppendByte(',')
		}
		val.AppendEle(e, buf, i)
	}
	buf.AppendByte(']')
}

func (e *JsonEncoder) AppendObject(buf *Buffer, val FieldVal) {
	val.Encode(e, buf)
}

func (e *JsonEncoder) Encode(buf *Buffer, entry *Entry, fields ...Field) error {
	buf.AppendByte('{')
	// Encode message level.
	e.encodeKey(buf, KeyLevel)
	if e.colorEnabled {
		e.AppendString(buf, entry.Level.ColorfulStr(e.levelLower))
	} else {
		e.AppendString(buf, entry.Level.Str(e.levelLower))
	}
	// Encode message time.
	if !e.timeDisabled {
		buf.AppendByte(',')
		e.encodeKey(buf, KeyTime)
		e.AppendTime(buf, entry.Time)
	}
	// Encode message text.
	buf.AppendByte(',')
	e.encodeKey(buf, KeyMsg)
	e.AppendString(buf, entry.Msg)
	// Encode message fields.
	if len(fields) > 0 {
		buf.AppendByte(',')
		for i, field := range fields {
			if i > 0 {
				buf.AppendByte(',')
			}
			e.encodeKey(buf, field.Key)
			field.Val.Encode(e, buf)
		}
	}
	buf.AppendByte('}')
	// Encode the line ending.
	buf.AppendString(e.lineEnding)
	return nil
}

func (e *JsonEncoder) encodeKey(buf *Buffer, key string) {
	e.AppendString(buf, key)
	buf.AppendByte(':')
}
