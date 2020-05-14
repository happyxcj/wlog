package wlog

import "io"

type Writer interface {
	io.Writer

	// Flush flushes any buffered logs to the disk.
	Flush() error

	io.Closer
}
