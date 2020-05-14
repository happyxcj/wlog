package wlog

import "fmt"

// globalLogger is the logger that can be conveniently used in all packages.
var globalLogger *Logger

// UseGlobalLogger initializes the global logger, if the given logger is nil
// use method NewConsoleConfig().Create(...Option) to initialize it.
func UseGlobalLogger(logger *Logger) {
	if logger != nil {
		globalLogger = logger
		return
	}
	var err error
	globalLogger, err = NewConsoleConfig().Create()
	if err != nil {
		panic(fmt.Sprint("failed to initialize the global logger, error: ", err.Error()))
	}
}

// WithOpts is the WithOpts method of a Logger that can be conveniently used in all packages.
func WithOpts(opts ...LoggerOpt) *Logger {
	return globalLogger.WithOpts(opts...)
}

// With is the With method of a Logger that can be conveniently used in all packages.
func With(fields ...Field) *Logger {
	return globalLogger.With(fields...)
}

// Withp is the Withp method of a Logger that can be conveniently used in all packages.
func Withp(pairs ...interface{}) *Logger {
	return globalLogger.Withp(pairs...)
}

// Debug is the Debug method of a Logger that can be conveniently used in all packages.
func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

// Info is the Info method of a Logger that can be conveniently used in all packages.
func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

// WithOpts is the WithOpts method of a Logger that can be conveniently used in all packages.
func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

// Error is the Error method of a Logger that can be conveniently used in all packages.
func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

// Fatal is the Fatal method of a Logger that can be conveniently used in all packages.
func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

// Panic is the Panic method of a Logger that can be conveniently used in all packages.
func Panic(args ...interface{}) {
	globalLogger.Panic(args...)
}

// Debugf is the Debugf method of a Logger that can be conveniently used in all packages.
func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

// Infof is the Infof method of a Logger that can be conveniently used in all packages.
func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

// Warnf is the Warnf method of a Logger that can be conveniently used in all packages.
func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

// Errorf is the Errorf method of a Logger that can be conveniently used in all packages.
func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

// Fatalf is the Fatalf method of a Logger that can be conveniently used in all packages.
func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

// Panicf is the Panicf method of a Logger that can be conveniently used in all packages.
func Panicf(format string, args ...interface{}) {
	globalLogger.Panicf(format, args...)
}

// Debugw is the Debugw method of a Logger that can be conveniently used in all packages.
func Debugw(msg string, fields ...Field) {
	globalLogger.Debugw(msg, fields...)
}

// Infow is the Infow method of a Logger that can be conveniently used in all packages.
func Infow(msg string, fields ...Field) {
	globalLogger.Infow(msg, fields...)
}

// Warnw is the Warnw method of a Logger that can be conveniently used in all packages.
func Warnw(msg string, fields ...Field) {
	globalLogger.Warnw(msg, fields...)
}

// Errorw is the Errorw method of a Logger that can be conveniently used in all packages.
func Errorw(msg string, fields ...Field) {
	globalLogger.Errorw(msg, fields...)
}

// Fatalw is the Fatalw method of a Logger that can be conveniently used in all packages.
func Fatalw(msg string, fields ...Field) {
	globalLogger.Fatalw(msg, fields...)
}

// Panicw is the Panicw method of a Logger that can be conveniently used in all packages.
func Panicw(msg string, fields ...Field) {
	globalLogger.Panicw(msg, fields...)
}

// Debugp is the Debugp method of a Logger that can be conveniently used in all packages.
func Debugp(msg string, pairs ...interface{}) {
	globalLogger.Debugp(msg, pairs...)
}

// Infop is the Infop method of a Logger that can be conveniently used in all packages.
func Infop(msg string, pairs ...interface{}) {
	globalLogger.Infop(msg, pairs...)
}

// Warnp is the Warnp method of a Logger that can be conveniently used in all packages.
func Warnp(msg string, pairs ...interface{}) {
	globalLogger.Warnp(msg, pairs...)
}

// Errorp is the Errorp method of a Logger that can be conveniently used in all packages.
func Errorp(msg string, pairs ...interface{}) {
	globalLogger.Errorp(msg, pairs...)
}

// Fatalp is the Fatalp method of a Logger that can be conveniently used in all packages.
func Fatalp(msg string, pairs ...interface{}) {
	globalLogger.Fatalp(msg, pairs...)
}

// Panicp is the Panicp method of a Logger that can be conveniently used in all packages.
func Panicp(msg string, pairs ...interface{}) {
	globalLogger.Panicp(msg, pairs...)
}

// Flush is the Flush method of a Logger that can be conveniently used in all packages.
func Flush() error{
	return globalLogger.Flush()
}

// Close is the Close method of a Logger that can be conveniently used in all packages.
func Close() error{
	return globalLogger.Close()
}
