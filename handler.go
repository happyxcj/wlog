package wlog

// Handler interface describes how to handle a complete log message with any fields.
type Handler interface {
	// With Adds structured information outside the log message body.
	With(fields ...Field) Handler

	// Write encodes the entry and any fields and writes them to the specified writer.
	Write(entry *Entry, fields ...Field) error

	// Flush flushes any buffered logs to the disk.
	Flush() error

	// Close closes the handler.
	Close() error
}
