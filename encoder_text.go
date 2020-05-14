package wlog

import "time"

type TextEncoder struct {
	BasicObjEncoder
	EncoderOpts
}

func NewTextEncoder(opts ...EncoderOpt) *TextEncoder {
	e := &TextEncoder{}
	e.timeEncoder = &FastTextTimeEncoder{}
	e.durationEncoder = &StrDurationEncoder{}
	e.lineEnding = defaultLineEnding
	e.WithOpts(opts...)
	return e
}

func (e *TextEncoder) AppendDuration(buf *Buffer, val time.Duration) {
	e.durationEncoder.Append(buf, e, val)
}

func (e *TextEncoder) AppendTime(buf *Buffer, val time.Time) {
	e.timeEncoder.Append(buf, e, val)
}

func (e *TextEncoder) AppendArray(buf *Buffer, val ArrayEncoder) {
	buf.AppendByte('[')
	for i, size := 0, val.Size(); i < size; i++ {
		if i > 0 {
			buf.AppendByte(' ')
		}
		val.AppendEle(e, buf, i)
	}
	buf.AppendByte(']')
}

func (e *TextEncoder) AppendObject(buf *Buffer, val FieldVal) {
	val.Encode(e, buf)
}

func (e *TextEncoder) Encode(buf *Buffer, entry *Entry, fields ...Field) error {
	// Encode message level.
	buf.AppendByte('[')
	if e.colorEnabled {
		e.AppendString(buf, entry.Level.ColorfulStr(e.levelLower))
	} else {
		e.AppendString(buf, entry.Level.Str(e.levelLower))
	}
	buf.AppendByte(']')
	// Encode message time.
	if !e.timeDisabled {
		e.encodeSep(buf)
		e.AppendTime(buf, entry.Time)
		//e.AppendString(buf, entry.Time.Format(e.timeLayout))
	}
	// Encode message text.
	if len(entry.Msg) > 0 {
		e.encodeSep(buf)
		e.AppendString(buf, entry.Msg)
	}
	n := len(fields)
	// Encode message fields.
	if n > 0 {
		e.encodeSep(buf)
		buf.AppendByte('[')
		for i, field := range fields {
			if i > 0 {
				buf.AppendByte(' ')
			}
			e.AppendString(buf, field.Key)
			buf.AppendByte('=')
			field.Val.Encode(e, buf)
		}
		buf.AppendByte(']')
	}
	// Encode the line ending.
	buf.AppendString(e.lineEnding)
	return nil
}

// encodeSep encodes a specified separator to the buf between two independent fields.
// The independent fields are as follows: "Level", "Time", "Msg" and "Field".
func (e *TextEncoder) encodeSep(buf *Buffer) {
	buf.AppendByte(' ')
	buf.AppendByte(' ')
}
