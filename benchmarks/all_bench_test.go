package benchmarks

import (
	"testing"
	"go.uber.org/zap"
)



func BenchmarkMsg(b *testing.B) {
	b.Log("benchmark message")

	b.Run("wlog", func(b *testing.B) {
		logger := newWLogLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow(_fakeMsg)
			}
		})
	})

	b.Run("uber-go/zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_fakeMsg)
			}
		})
	})

	b.Run("fmt.Println", func(b *testing.B) {
		logger := newStdLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Println(_fakeMsg)
			}
		})
	})

	b.Run("Sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_fakeMsg)
			}
		})
	})
}

func BenchmarkFormat(b *testing.B) {
	b.Log("benchmark format")

	b.Run("wlog", func(b *testing.B) {
		logger := newWLogLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof(_fakeMsgFormat,_fakeArg)
			}
		})
	})

	b.Run("uber-go/zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof(_fakeMsgFormat,_fakeArg)
			}
		})
	})

	b.Run("fmt.Println", func(b *testing.B) {
		logger := newStdLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Printf(_fakeMsgFormat,_fakeArg)
			}
		})
	})

	b.Run("Sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof(_fakeMsgFormat,_fakeArg)
			}
		})
	})
}

func BenchmarkWithFields(b *testing.B) {
	b.Log("benchmark with fields")

	b.Run("wlog", func(b *testing.B) {
		logger := newWLogLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow(_fakeMsg,fakeWLogFields()...)
			}
		})
	})

	b.Run("uber-go/zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(_fakeMsg,fakeZapFields()...)
			}
		})
	})

	b.Run("fmt.Println", func(b *testing.B) {
		logger := newStdLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Println(fakeStdMsgFields()...)
			}
		})
	})

	b.Run("Sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeLogrusFields()).Info(_fakeMsg)
			}
		})
	})
}

func BenchmarkWithPairs(b *testing.B) {
	b.Log("benchmark with pairs")

	b.Run("wlog", func(b *testing.B) {
		logger := newWLogLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infop(_fakeMsg,fakeWLogPairs()...)
			}
		})
	})

	b.Run("uber-go/zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.With(fakeZapPairs()...).Info(_fakeMsg)
			}
		})
	})

	b.Run("fmt.Println", func(b *testing.B) {
		logger := newStdLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Println(fakeStdMsgFields()...)
			}
		})
	})

	b.Run("Sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeLogrusFields()).Info(_fakeMsg)
			}
		})
	})
}

func BenchmarkFormatWithFields(b *testing.B) {
	b.Log("benchmark format with fields")

	b.Run("wlog", func(b *testing.B) {
		logger := newWLogLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.With(fakeWLogFields()...).Infof(_fakeMsgFormat,_fakeArg)
			}
		})
	})

	b.Run("uber-go/zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.With(fakeZapPairs()...).Infof(_fakeMsgFormat,_fakeArg)
			}
		})
	})

	b.Run("Sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeLogrusFields()).Infof(_fakeMsgFormat,_fakeArg)
			}
		})
	})
}


func BenchmarkFormatWithPairs(b *testing.B) {
	b.Log("benchmark format with pairs")

	b.Run("wlog", func(b *testing.B) {
		logger := newWLogLogger()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Withp(fakeWLogPairs()...).Infof(_fakeMsgFormat,_fakeArg)
			}
		})
	})

	b.Run("uber-go/zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.With(fakeZapPairs()...).Infof(_fakeMsgFormat,_fakeArg)
			}
		})
	})

	b.Run("Sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeLogrusFields()).Infof(_fakeMsgFormat,_fakeArg)
			}
		})
	})
}