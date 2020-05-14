package wlog

import (
	"strings"
)

// log level enum
const (
	// DebugLvl is usually used in development to print some test information, find any bugs, etc.
	DebugLvl Level = iota
	// InfoLvl is the default logging level in production.
	// It is usually used to record some necessary information.
	InfoLvl
	// WarnLvl is usually used to record some more important information than DebugLvl.
	WarnLvl
	// WarnLvl is usually used to record some application error that will not occur
	// if the application runs well.
	ErrorLvl
	// FatalLvl always logs a message, then the application calls os.Exit(1).
	FatalLvl
	// PanicLvl always logs a message, then the application panics.
	PanicLvl
)

const (
	levelNum = 6
)

// Level is the type defined for log level.
type Level uint8

// levels is a array consist of all log levels.
var levels = [6]Level{DebugLvl, InfoLvl, WarnLvl, ErrorLvl, FatalLvl, PanicLvl}

// levelColors contains all colors corresponding to all log levels.
var levelColors = [levelNum]Color{Green, Blue, Yellow, Red, Red, Red}

// levelColors contains all strings corresponding to all log levels.
var levelStrings = [levelNum]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL", "PANIC"}

// levelUpperColorfulPrefixes contains all upper colorful strings corresponding to all log levels.
var levelUpperColorfulStrings [levelNum]string

// levelUpperColorfulPrefixes contains all lower colorful strings corresponding to all log levels.
var levelLowerColorfulStrings [levelNum]string

func init() {
	for i, color := range levelColors {
		levelUpperColorfulStrings[i] = color.With(levels[i].UpperStr())
		levelLowerColorfulStrings[i] = color.With(levels[i].LowerStr())
	}
}

// UpperStr returns a uppercase string of the log level.
func (l Level) UpperStr() string {
	return levelStrings[l]
}

// LowerStr returns a lowercase string of the log level.
func (l Level) LowerStr() string {
	return strings.ToLower(levelStrings[l])
}

// UpperColorfulStr returns a uppercase colorful string of the log level.
func (l Level) UpperColorfulStr() string {
	return levelUpperColorfulStrings[l]
}

// LowerColorfulStr returns a lowercase colorful string of the log level.
func (l Level) LowerColorfulStr() string {
	return levelLowerColorfulStrings[l]
}

// Str returns a lowercase string of the log level if the given isLower is true,
// otherwise returns a uppercase string of the log level.
func (l Level) Str(isLower bool) string {
	if isLower{
		return l.LowerStr()
	}
	return l.UpperStr()
}

// ColorfulStr returns a lowercase colorful string of the log level if the given isLower is true,
// otherwise returns a uppercase colorful string of the log level.
func (l Level) ColorfulStr(isLower bool) string {
	if isLower{
		return l.LowerColorfulStr()
	}
	return l.UpperColorfulStr()
}
