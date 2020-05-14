package wlog

import (
	"time"
	"sync/atomic"
)

const defaultFlushInterval = 3 * time.Second

type TimingFlushWriter struct {
	Writer
	// interval is the flushing interval.
	// It.s default value is 3 seconds and at least 1 seconds.
	interval    time.Duration
	waitingFlag uint32
	wakeupCh    chan struct{}
	closeCh     chan struct{}
}

func NewTimingFlushWriter(inner Writer, interval time.Duration) *TimingFlushWriter {
	if interval < time.Second {
		interval = defaultFlushInterval
	}
	fw := &TimingFlushWriter{
		Writer:        inner,
		interval: interval,
		wakeupCh: make(chan struct{}, 1),
		closeCh:  make(chan struct{}, 1),
	}
	go fw.flushLoop()
	return fw
}

func (w *TimingFlushWriter) Write(bs []byte) (n int, err error) {
	n, err = w.Writer.Write(bs)
	if n != 0 && atomic.CompareAndSwapUint32(&w.waitingFlag, 0, 1) {
		select {
		case w.wakeupCh <- struct{}{}:
		default:
		}
	}
	return
}

func (w *TimingFlushWriter) Flush() error {
	return w.Writer.Flush()
}

func (w *TimingFlushWriter) Close() error {
	select {
	case w.closeCh <- struct{}{}:
	default:
	}
	return w.Writer.Close()
}

func (w *TimingFlushWriter) flushLoop() {
	t := time.NewTimer(w.interval)
	defer t.Stop()
	for {
		if atomic.LoadUint32(&w.waitingFlag) == 0 {
			select {
			case <-w.wakeupCh:
			case <-w.closeCh:
				return
			}
		}
		t.Reset(w.interval)
		select {
		case <-t.C:
			// Reset the waitingFlag to 0 must be performed before the flushing operation.
			atomic.StoreUint32(&w.waitingFlag, 0)
			w.Flush()
		case <-w.closeCh:
			return
		}

	}
}
