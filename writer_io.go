package wlog

import (
	"io"
)

type IOWriter struct {
	w io.Writer
}

func NewIOWriter(w io.Writer) *IOWriter {
	return &IOWriter{w: w}
}

func (w *IOWriter) Write(bs []byte) (int, error) {
	return w.w.Write(bs)
}

func (w *IOWriter) Flush() error {
	return nil
}

func (w *IOWriter) Close() error {
	return nil
}
