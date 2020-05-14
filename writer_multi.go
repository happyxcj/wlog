package wlog

import (
	"fmt"
	"errors"
	"math"
)

type MultiWriter struct {
	ws []Writer
}

func NewMultiWriter(ws []Writer) *MultiWriter {
	return &MultiWriter{ws: ws}
}

// Writer returns the Writer under the specified index.
func (w *MultiWriter) Writer(index int) (Writer, error) {
	if index < 0 || index >= len(w.ws) {
		return nil, errors.New("MultiWriter: the index is invalid")
	}
	return w.ws[index], nil
}

// Write returns the minimum number of bytes written from bs (0 <= n <= len(bs)) in every underlying writer.
func (w *MultiWriter) Write(bs []byte) (int, error) {
	var tarErr error
	tarNum := math.MaxInt64
	for _, w := range w.ws {
		n, err := w.Write(bs)
		if n < tarNum {
			tarNum = n
		}
		tarErr = multiErr(tarErr, err)
	}
	return tarNum, tarErr
}

func (w *MultiWriter) Flush() error {
	var err error
	for _, w := range w.ws {
		err = multiErr(err, w.Flush())
	}
	return err
}

func (w *MultiWriter) Close() error {
	var err error
	for _, w := range w.ws {
		err = multiErr(err, w.Close())
	}
	return err
}

func multiErr(err1, err2 error) error {
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}
	return errors.New(fmt.Sprint(err1.Error(), ";", err2.Error()))
}
