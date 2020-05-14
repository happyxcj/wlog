package wlog

import (
	"testing"
	"fmt"
	"sync"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
)

func TestLogger(t *testing.T) {
	w := NewIOWriter(ioutil.Discard)
	logger := NewLogger(NewBaseHandler(w, NewTextEncoder()))
	logger.Debug("test logger")

}

func TestLoggef(t *testing.T) {
	w := NewIOWriter(ioutil.Discard)
	logger := NewLogger(NewBaseHandler(w, NewTextEncoder()))
	logger.Infof("%v", "test logger")
}

func TestLoggerw(t *testing.T) {
	w := NewIOWriter(ioutil.Discard)
	logger := NewLogger(NewBaseHandler(w, NewTextEncoder()))
	logger.Warnw("test logger", String("name", "xcj"), Int("age", 10))

}

func TestLoggerp(t *testing.T) {
	w := NewIOWriter(ioutil.Discard)
	logger := NewLogger(NewBaseHandler(w, NewTextEncoder()))
	logger.Errorp("test logger", "name", "xcj", Int("age", 10))
}

func TestLoggerWith(t *testing.T) {
	w := NewIOWriter(ioutil.Discard)
	logger := NewLogger(NewBaseHandler(w, NewTextEncoder()))
	logger.With(String("name", "xcj"), Int("age", 10)).Debugw("test logger", Int("age", 10))
}

func TestLoggerWithp(t *testing.T) {
	w := NewIOWriter(ioutil.Discard)
	logger := NewLogger(NewBaseHandler(w, NewTextEncoder()))
	logger.Withp("name", "xcj", Int("age", 10)).Debugp("test logger", "age", 10)

}

func TestLoggerConcurrency(t *testing.T) {
	w := NewIOWriter(ioutil.Discard)
	logger := NewLogger(NewBaseHandler(w, NewTextEncoder()))
	var wg sync.WaitGroup
	wg.Add(5)
	for g := 0; g < 5; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				logger.Errorp("test logger", "name", "xcj", Int("age", 10))
			}
		}()
	}
	wg.Wait()
}

func TestLoggerPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover panic error: ", err)
		}
	}()
	w := NewIOWriter(ioutil.Discard)
	logger := NewLogger(NewBaseHandler(w, NewTextEncoder()))
	assert.Panics(t, func() {
		logger.Panicw("test logger", String("name", "xcj"), Int("age", 10))
	}, "Expected panic")
}

func TestLoggerFatal(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		w := NewIOWriter(ioutil.Discard)
		logger := NewLogger(NewBaseHandler(w, NewTextEncoder()))
		logger.Fatalw("test logger", String("name", "xcj"), Int("age", 10))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestLoggerFatal")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
