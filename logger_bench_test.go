package wlog

import (
	"testing"
	"io/ioutil"
	"time"
	"errors"
)

func doBenchmark(b *testing.B, f func(*Logger)) {
	w := NewIOWriter(ioutil.Discard)
	logger := NewLogger(NewBaseHandler(w, NewTextEncoder()))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f(logger)
		}
	})
}

func BenchmarkMsg(b *testing.B) {
	doBenchmark(b, func(logger *Logger) {
		logger.Debugw("log pure message")
	})
}

func BenchmarkArgs(b *testing.B) {
	doBenchmark(b, func(logger *Logger) {
		logger.Debug("log one argument", "arg1")
	})
}

func BenchmarkFormattedMsg(b *testing.B) {
	doBenchmark(b, func(logger *Logger) {
		logger.Debugf("log formatted message: %v=%v", "author", "xcj")
	})
}

func BenchmarkMsgWithFields(b *testing.B) {
	err := errors.New("fake error of wlog")
	t := time.Now()
	doBenchmark(b, func(logger *Logger) {
		logger.Debugw("log pure message with fields",
			Bool("bool", true),
			Int("int", 1),
			Float64("float64", 2.0),
			String("string", "xcj"),
			Time("time", t),
			Err("error", err), )
	})
}

func BenchmarkFormattedMsgWithFields(b *testing.B) {
	err := errors.New("fake error of wlog")
	t := time.Now()
	doBenchmark(b, func(logger *Logger) {
		logger.With(Bool("bool", true),
			Int("int", 1),
			Float64("float64", 2.0),
			String("string", "xcj"),
			Time("time", t),
			Err("error", err), ).
			Debugf("log formatted message: %v=%v", "author", "xcj")
	})
}

func BenchmarkMsgWithPairs(b *testing.B) {
	err := errors.New("fake error of wlog")
	t := time.Now()
	doBenchmark(b, func(logger *Logger) {
		logger.Debugp("log pure message with fields",
			"bool", true,
			"int", 1,
			"float64", 2.0,
			"string", "xcj",
			"time", t,
			"timestamp", t,
			"error", err, )
	})
}

func BenchmarkFormattedMsgWithPairs(b *testing.B) {
	err := errors.New("fake error of wlog")
	t := time.Now()
	doBenchmark(b, func(logger *Logger) {
		logger.Withp("bool", true,
			"int", 1,
			"float64", 2.0,
			"string", "xcj",
			"time", t,
			"error", err, ).
			Debugf("log formatted message: %v=%v", "author", "xcj")
	})
}

func BenchmarkMsgWith50Fields(b *testing.B) {
	w := NewIOWriter(ioutil.Discard)
	logger := NewLogger(NewBaseHandler(w, NewTextEncoder()))
	fields := make([]Field, 50)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 50; i++ {
			if i < 25 {
				fields[i] = Int("int", i)
			} else {
				fields[i] = String("string", "same")
			}
		}
		logger.Debugw("log pure message with 50 fields", fields...)
	}

}
