package wlog

import (
	"sync"
	"fmt"
	"io"
	"os"
	"errors"
	"time"
)

const (
	defaultBufMinSize = 4096
	defaultBufMaxSize = 500 * 1 << 20

	initBuffersSize = 5
)

// BufWriter implements buffering for an Writer object.
// After all data has been written, the client should call the Flush or Close method
// to guarantee all data has been forwarded to the underlying Writer.
// It's safe to call the Write, Flush and Close methods concurrently.
//
// Note that it will discard the later messages directly if the buffered data overflows,
// so we force the BufWriter's maxSize to be no less than "500 * 1 << 20" by default,
// you can increase it according to the actual situation.
//
// TODO: Is it necessary to Wait for a specified timeout before discarding the later messages when the buffered data overflows.
type BufWriter struct {
	Writer
	// buf is the current buffer to cache the data to be appended to the "buffers".
	buf []byte
	// nextBuf is used to get the next buffer first when the current buffer is full.
	nextBuf []byte
	// nextBufBusy indicates whether the "nextBuf" is busy.
	nextBufBusy bool
	// buffers is used to cache the pending buffers to be written to the underlying Writer.
	buffers [][]byte
	// minSize is the minimum buffer size before writing the buffered data to the underlying Writer.
	minSize int
	// maxSize is the maximum size of all buffered data.
	maxSize int
	// bufferedSize is the current size of all buffered data
	bufferedSize int
	mu           sync.Mutex
	cond         *sync.Cond
	// condWaiting indicates the "writeLoop" is waiting on condition wait.
	condWaiting bool
	// isClosed indicates whether the BufWriter has been closed.
	isClosed bool
	// errW is used to output the interval error when writing data to the underlying "w".
	// os.Stderr is the default io writer.
	errW io.Writer
	wg   sync.WaitGroup
}

type BufWriterOpt func(w *BufWriter)

// SetBufMinSize sets the minimum buffer size of the BufWriter.
func SetBufMinSize(minSize int) BufWriterOpt {
	return func(w *BufWriter) {
		if minSize <= 0 {
			w.minSize = defaultBufMinSize
			return
		}
		w.minSize = minSize
	}
}

// SetBufMaxSize sets the maximum buffered size of the BufWriter.
// If the given maxSize is less than 500 * 1 << 20,
// it sets the maximum buffer size to 500 * 1 << 20.
func SetBufMaxSize(maxSize int) BufWriterOpt {
	return func(w *BufWriter) {
		if maxSize <= defaultBufMaxSize {
			w.maxSize = defaultBufMaxSize
			return
		}
		w.maxSize = maxSize
	}
}

// SetBufErrW sets the underlying io.Writer of the BufWriter to output the internal error.
func SetBufErrW(errW io.Writer) BufWriterOpt {
	return func(w *BufWriter) {
		w.errW = errW
	}
}

func NewBufWriter(inner Writer, opts ...BufWriterOpt) *BufWriter {
	w := &BufWriter{
		Writer:  inner,
		minSize: defaultBufMinSize,
		maxSize: defaultBufMaxSize,
		errW:    os.Stderr,
	}
	for _, opt := range opts {
		opt(w)
	}
	w.buf = w.makeBuffer()
	w.nextBuf = w.makeBuffer()
	w.buffers = make([][]byte, 0, initBuffersSize)
	w.cond = sync.NewCond(&w.mu)
	w.wg.Add(1)
	go w.writeLoop()
	return w
}

// newBuffer returns a new buffer.
// In most case, the size of buffer will be greater than w's "minSize" after appending the last
// message to the buffer, in order to ensure enough space to cache the last message,
// initialize the capacity of the buffer to "2*minSize" instead of "minSize".
func (w *BufWriter) makeBuffer() []byte {
	return make([]byte, 0, 2*w.minSize)
}

func (w *BufWriter) Write(bs []byte) (int, error) {
	n := len(bs)
	if n == 0 {
		return 0, nil
	}
	w.mu.Lock()
	if w.isClosed {
		w.mu.Unlock()
		return 0, errors.New("the BufWriter had been closed")
	}
	if w.bufferedSize >= w.maxSize {
		w.mu.Unlock()
		return 0, errors.New("the BufWriter's buffered data overflow")
	}
	// The buffered size may exceed "w.minSize" after caching the given bs,
	// in order to ensure the integrity of the given bs, we do not limit the length
	// of "w.buf" to "w.minSize".
	w.buf = append(w.buf, bs...)
	nBuf := len(w.buf)
	if nBuf >= w.minSize {
		w.bufferedSize += nBuf
		w.flushBuf()
	}
	w.mu.Unlock()
	return n, nil
}

func (w *BufWriter) Flush() error {
	w.mu.Lock()
	if len(w.buf) > 0 {
		w.flushBuf()
	}
	w.mu.Unlock()
	return nil
}

func (w *BufWriter) Close() error {
	w.mu.Lock()
	if w.isClosed {
		w.mu.Unlock()
		return nil
	}
	if len(w.buf) > 0 {
		w.flushBuf()
	}
	w.isClosed = true
	w.mu.Unlock()
	w.wg.Wait()
	return w.Writer.Close()
}

// flushBuf flushes the buffered data to the underlying Writer.
func (w *BufWriter) flushBuf() {
	w.buffers = append(w.buffers, w.buf)
	if w.condWaiting {
		w.cond.Signal()
	}
	// Reset the current buffer to cache later data.
	if !w.nextBufBusy {
		w.buf = w.nextBuf
		w.nextBufBusy = true
		return
	}
	w.buf = w.makeBuffer()
}

func (w *BufWriter) writeLoop() {
	defer w.wg.Done()
	newBuffers := make([][]byte, 0, initBuffersSize)
	newNextBuf := w.makeBuffer()
	for {
		w.mu.Lock()
		if len(w.buffers) <= 0 {
			// In order to write all buffered data to underlying writer after "Close" method is called.
			// do not exit the loop even if it is detected that the w has been closed when the "buf" is not empty.
			if w.isClosed {
				w.mu.Unlock()
				return
			}
			w.condWaiting = true
			w.cond.Wait()
			w.condWaiting = false
		}
		toWriteBuffers := w.buffers
		w.buffers = newBuffers
		w.nextBuf = newNextBuf
		w.nextBufBusy = false
		w.bufferedSize = 0
		w.mu.Unlock()
		var err error
		for _, data := range toWriteBuffers {
			_, err = w.Writer.Write(data)
			if err != nil {
				t := time.Now().Format("2006-01-02 15:04:05")
				fmt.Fprintf(w.errW, "BufWriter: unable to Write data to underlying writer at time: %v, error: %v\n", t, err)
			}
		}
		newNextBuf = toWriteBuffers[0][:0]
		// Reset the next buffers and check to release the useless memory.
		if len(toWriteBuffers) < 25 && cap(toWriteBuffers) > 50 {
			newBuffers = make([][]byte, 0, initBuffersSize)
		} else {
			newBuffers = toWriteBuffers[:0]
		}
	}
}
