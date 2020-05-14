package benchmarks

import (
	"log"
	"io/ioutil"
)

func newStdLogger() *log.Logger {
	logger := log.New(ioutil.Discard, "", log.LstdFlags)
	return logger
}

func fakeStdMsgFields() []interface{} {
	return []interface{}{
		_fakeMsg,
		" int:", _fakeInts[0],
		" ints:", _fakeInts,
		" float64:", _fakeFloat64s[0],
		" float64s:", _fakeFloat64s,
		" string:", _fakeStrings[0],
		" strings:", _fakeStrings,
		" time:", _fakeTimes[0],
		" times:", _fakeTimes,
		" duration:", _fakeDurations[0],
		" durations:", _fakeDurations,
		" error:", _fakeErrors[0],
		" errors:", _fakeErrors,
	}
}

