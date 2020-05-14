package benchmarks

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
)

type zapDiscarder struct {
}

func (d *zapDiscarder) Write(b []byte) (int, error) {
	return ioutil.Discard.Write(b)
}

func (d *zapDiscarder) Sync() error {
	return nil
}

func newZapLogger(lvl zapcore.Level) *zap.Logger {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeDuration = zapcore.StringDurationEncoder
	ec.EncodeTime = zapcore.RFC3339TimeEncoder
	enc := zapcore.NewJSONEncoder(ec)
	return zap.New(zapcore.NewCore(
		enc,
		&zapDiscarder{},
		lvl,
	))
}

func fakeZapFields() []zap.Field {
	return []zap.Field{
		zap.Int("int", _fakeInts[0]),
		zap.Ints("ints", _fakeInts),
		zap.Float64("float64", _fakeFloat64s[0]),
		zap.Float64s("float64s", _fakeFloat64s),
		zap.String("string", _fakeStrings[0]),
		zap.Strings("strings", _fakeStrings),
		zap.Time("time", _fakeTimes[0]),
		zap.Times("times", _fakeTimes),
		zap.Duration("duration", _fakeDurations[0]),
		zap.Durations("durations", _fakeDurations),
		zap.Error(_fakeErrors[0]),
		zap.Errors("errors", _fakeErrors),
	}
}

func fakeZapPairs() []interface{} {
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
