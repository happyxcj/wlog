package benchmarks

import (
	"io/ioutil"

	"github.com/Sirupsen/logrus"
)

func newLogrus() *logrus.Logger {
	return &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}
}

func fakeLogrusFields() logrus.Fields {
	return logrus.Fields{
		"int":       _fakeInts[0],
		"ints":      _fakeInts,
		"float64":   _fakeFloat64s[0],
		"float64s":  _fakeFloat64s,
		"string":    _fakeFloat64s[0],
		"strings":   _fakeFloat64s,
		"time":      _fakeTimes[0],
		"times":     _fakeTimes,
		"duration":  _fakeDurations[0],
		"durations": _fakeDurations,
		"error":     _fakeErrors[0],
		"errors":    _fakeErrors,
	}
}
