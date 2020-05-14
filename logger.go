package wlog

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Logger contains all common data needed for logging and contains methods used to log messages.
type Logger struct {
	// minLvl is the minimum level allowed to log a message.
	// It's default value is "DebugLvl".
	minLvl Level
	h      Handler
	// errW is used to output the internal error when logging the message.
	// os.Stderr is the default io writer.
	errW io.Writer
}

type LoggerOpt func(l *Logger)

// SetLogMinLvl sets the minimum logging level of the Logger
func SetLogMinLvl(lvl Level) LoggerOpt {
	return func(l *Logger) {
		l.minLvl = lvl
	}
}

// SetLogHandler sets the underlying Handler level of the Logger
func SetLogHandler(h Handler) LoggerOpt {
	return func(l *Logger) {
		l.h = h
	}
}

// WrapLogHandler wraps the underlying Handler level of the Logger
func WrapLogHandler(f func(Handler) Handler) LoggerOpt {
	return func(l *Logger) {
		l.h = f(l.h)
	}
}

// SetLogErrW sets the underlying io.Writer of the Logger to output the internal error.
func SetLogErrW(errW io.Writer) LoggerOpt {
	return func(l *Logger) {
		l.errW = errW
	}
}

func NewLogger(h Handler, opts ...LoggerOpt) *Logger {
	l := &Logger{
		minLvl: DebugLvl,
		h:      h,
		errW:   os.Stderr,
	}
	return l.WithOpts(opts...)
}

// WithOpts returns a new Logger by cloning l, and then applies the given opts for it.
func (l *Logger) WithOpts(opts ...LoggerOpt) *Logger {
	if len(opts) == 0 {
		return l
	}
	clone := l.clone()
	for _, opt := range opts {
		opt(clone)
	}
	return clone
}

// With returns a new Logger by cloning l, and then adds the given fields to it.
func (l *Logger) With(fields ...Field) *Logger {
	if len(fields) == 0 {
		return l
	}
	clone := l.clone()
	clone.h = l.h.With(fields...)
	return clone
}

// Withp returns a new Logger by cloning l, and then adds the given key-value pairs to it.
// In addition to key-value pairs, the pairs can also contains any independent fields of type *Field.
func (l *Logger) Withp(pairs ...interface{}) *Logger {
	fields := pairsToFields(pairs...)
	return l.With(fields...)
}

// Debug logs a message to be constructed at debug-level. It uses Sprint(...interface{}) to construct the args.
//
// Note that: l.With(...*Field).Debug(...interface{}) is the equivalent of l.Debugw(string, ...*Field),
// in most case, the first method may be more convenient to use, but if the debug-level is disabled,
// the second method is more efficient.
func (l *Logger) Debug(args ...interface{}) {
	l.Debugw(fmt.Sprint(args...))
}

// Info logs a message to be constructed at info-level. It uses Sprint(...interface{}) to construct the args.
//
// Note that: l.With(...*Field).Info(...interface{}) is the equivalent of l.Infow(string, ...*Field),
// in most case, the first method may be more convenient to use, but if the info-level is disabled,
// the second method is more efficient.
func (l *Logger) Info(args ...interface{}) {
	l.Infow(fmt.Sprint(args...))
}

// Warn logs a message to be constructed at warn-level. It uses Sprint(...interface{}) to construct the args.
//
// Note that: l.With(...*Field).Warn(...interface{}) is the equivalent of l.Warnw(string, ...*Field),
// in most case, the first method may be more convenient to use, but if the warn-level is disabled,
// the second method is more efficient.
func (l *Logger) Warn(args ...interface{}) {
	l.Warnw(fmt.Sprint(args...))
}

// Error logs a message to be constructed at error-level. It uses Sprint(...interface{}) to construct the args.
//
// Note that: l.With(...*Field).Error(...interface{}) is the equivalent of l.Errorw(string, ...*Field),
// in most case, the first method may be more convenient to use, but if the error-level is disabled,
// the second method is more efficient.
func (l *Logger) Error(args ...interface{}) {
	l.Errorw(fmt.Sprint(args...))
}

// Fatal logs a message to be constructed at fatal-level. It uses Sprint(...interface{}) to construct the args.
//
// Note that: l.With(...*Field).Fatal(...interface{}) is the equivalent of l.Fatalw(string, ...*Field),
// in most case, the first method may be more convenient to use, but if the fatal-level is disabled,
// the second method is more efficient.
func (l *Logger) Fatal(args ...interface{}) {
	l.Fatalw(fmt.Sprint(args...))
}

// Panic logs a message to be constructed at panic-level. It uses Sprint(...interface{}) to construct the args.
//
// Note that: l.With(...*Field).Panic(...interface{}) is the equivalent of l.Panicw(string, ...*Field),
// in most case, the first method may be more convenient to use, but if the panic-level is disabled,
// the second method is more efficient.
func (l *Logger) Panic(args ...interface{}) {
	l.Panicw(fmt.Sprint(args...))
}

// Debugf logs a message to be formatted at debug-level.
// It uses fmt.Sprintf(string, ...interface{}) to format the args.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debugw(fmt.Sprintf(format, args...))
}

// Infof logs a message at to be formatted info-level.
// It uses fmt.Sprintf(string, ...interface{}) to format the args.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Infow(fmt.Sprintf(format, args...))
}

// Warnf logs a message to be formatted at warn-level.
// It uses fmt.Sprintf(string, ...interface{}) to format the args.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warnw(fmt.Sprintf(format, args...))
}

// Error logs a message to be formatted at error-level.
// It uses fmt.Sprintf(string, ...interface{}) to format the args.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Errorw(fmt.Sprintf(format, args...))
}

// Fatal logs a message at to be formatted fatal-level.
// It uses fmt.Sprintf(string, ...interface{}) to format the args.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatalw(fmt.Sprintf(format, args...))
}

// Panic logs a message at to be formatted panic-level.
// It uses fmt.Sprintf(string, ...interface{}) to format the args.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Panicw(fmt.Sprintf(format, args...))
}

// Debugw logs a message at debug-level with any fields.
func (l *Logger) Debugw(msg string, fields ...Field) {
	l.output(DebugLvl, msg, fields...)
}

// Infow logs a message at info-level with any fields.
func (l *Logger) Infow(msg string, fields ...Field) {
	l.output(InfoLvl, msg, fields...)
}

// Warnw logs a message at warn-level with any fields.
func (l *Logger) Warnw(msg string, fields ...Field) {
	l.output(WarnLvl, msg, fields...)
}

// Errorw logs a message at error-level with any fields.
func (l *Logger) Errorw(msg string, fields ...Field) {
	l.output(ErrorLvl, msg, fields...)
}

// Fatalw logs a message at fatal-level with any fields.
func (l *Logger) Fatalw(msg string, fields ...Field) {
	l.output(FatalLvl, msg, fields...)
}

// Panicw logs a message at panic-level with any fields.
func (l *Logger) Panicw(msg string, fields ...Field) {
	l.output(PanicLvl, msg, fields...)
}

// Debugp logs a message at debug-level with any key-value pairs.
// In addition to key-value pairs, the pairs can also contains any independent fields of type *Field.
func (l *Logger) Debugp(msg string, pairs ...interface{}) {
	l.outputPairs(DebugLvl, msg, pairs...)
}

// Infop logs a message at info-level with any key-value pairs.
// In addition to key-value pairs, the pairs can also contains any independent fields of type *Field.
func (l *Logger) Infop(msg string, pairs ...interface{}) {
	l.outputPairs(InfoLvl, msg, pairs...)
}

// Warnp logs a message at warn-level with any key-value pairs.
// In addition to key-value pairs, the pairs can also contains any independent fields of type *Field.
func (l *Logger) Warnp(msg string, pairs ...interface{}) {
	l.outputPairs(WarnLvl, msg, pairs...)
}

// Errorp logs a message at error-level with any key-value pairs.
// In addition to key-value pairs, the pairs can also contains any independent fields of type *Field.
func (l *Logger) Errorp(msg string, pairs ...interface{}) {
	l.outputPairs(ErrorLvl, msg, pairs...)
}

// Fatalp logs a message at fatal-level with any key-value pairs.
// In addition to key-value pairs, the pairs can also contains any independent fields of type *Field.
func (l *Logger) Fatalp(msg string, pairs ...interface{}) {
	l.outputPairs(FatalLvl, msg, pairs...)
}

// Panicp logs a message at panic-level with any key-value pairs.
// In addition to key-value pairs, the pairs can also contains any independent fields of type *Field.
func (l *Logger) Panicp(msg string, pairs ...interface{}) {
	l.outputPairs(PanicLvl, msg, pairs...)
}

// Flush flushes any buffered logs to the disk.
// It actually calls internal Handler's Flush method.
func (l *Logger) Flush() error {
	return l.h.Flush()
}

// Close closes the logger after flushing any buffered logs to the disk.
// It actually calls internal Handler's Close method.
func (l *Logger) Close() error {
	l.h.Flush()
	return l.h.Close()
}

func (l *Logger) output(lvl Level, msg string, fields ...Field) {
	if lvl < l.minLvl {
		return
	}
	e := getEntry()
	e.Set(lvl, msg)
	err := l.h.Write(e, fields...)
	putEntry(e)
	if err != nil {
		t:=time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(l.errW,"Logger: unable to write at time: %v, error: %v\n",t,err)
	}
	// Finally, regardless of the handling result, call os.Exit(1) or panic if necessary.
	if lvl < FatalLvl {
		return
	}
	l.Close()
	if lvl == FatalLvl {
		os.Exit(1)
	} else {
		panic(msg)
	}
}

func (l *Logger) outputPairs(lvl Level, msg string, pairs ...interface{}) {
	if len(pairs) == 0 {
		l.output(lvl, msg)
		return
	}
	fields := pairsToFields(pairs...)
	l.output(lvl, msg, fields...)
}

func pairsToFields(pairs ...interface{}) []Field {
	n := len(pairs)
	// Allocate enough space for the worst case.
	fields := make([]Field, 0, n)
	if n == 0 {
		return fields
	}
	var key, value interface{}
	for i := 0; i < n; {
		key = pairs[i]
		if field, ok := key.(Field); ok {
			fields = append(fields, field)
			i++
			continue
		}
		value = pairs[i+1]
		strKey, ok := key.(string)
		if !ok {
			strKey = fmt.Sprint(key)
		}
		fields = append(fields, Interface(strKey, value))
		i += 2
	}
	return fields
}

func (l *Logger) clone() *Logger {
	clone := *l
	return &clone
}
