package benchmarks

import (
	"github.com/happyxcj/wlog"
	"io/ioutil"
)

func newWLogLogger() *wlog.Logger {
	w := wlog.NewIOWriter(ioutil.Discard)
	l := wlog.NewLogger(wlog.NewBaseHandler(w, wlog.NewJsonEncoder()))
	return l
}

func fakeWLogFields() []wlog.Field {
	return []wlog.Field{
		wlog.Int("int", _fakeInts[0]),
		wlog.Ints("ints", _fakeInts),
		wlog.Float64("float64", _fakeFloat64s[0]),
		wlog.Float64s("float64s", _fakeFloat64s),
		wlog.String("string", _fakeStrings[0]),
		wlog.Strings("strings", _fakeStrings),
		wlog.Time("time", _fakeTimes[0]),
		wlog.Times("times", _fakeTimes),
		wlog.Duration("duration", _fakeDurations[0]),
		wlog.Durations("durations", _fakeDurations),
		wlog.Err("error", _fakeErrors[0]),
		wlog.Errs("errors", _fakeErrors),
	}
}

func fakeWLogPairs() []interface{} {
	return []interface{}{
		"int", _fakeInts[0],
		"ints", _fakeInts,
		"float64", _fakeFloat64s[0],
		"float64s", _fakeFloat64s,
		"string", _fakeStrings[0],
		"strings", _fakeStrings,
		"time", _fakeTimes[0],
		"times", _fakeTimes,
		"duration", _fakeDurations[0],
		"durations", _fakeDurations,
		"error", _fakeErrors[0],
		"errors", _fakeErrors,
	}
}
