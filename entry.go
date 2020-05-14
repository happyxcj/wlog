package wlog

import (
	"sync"
	"time"
)

var entryPool = sync.Pool{
	New: func() interface{} {
		return &Entry{}
	},
}

// getEntry returns a Entry form the entry pool.
func getEntry() *Entry {
	e, ok := entryPool.Get().(*Entry)
	if ok {
		return e
	}
	return &Entry{}
}

// getEntry puts the given e to the entry pool.
func putEntry(e *Entry) {
	entryPool.Put(e)
}

type Entry struct {
	Level Level
	Time  time.Time
	Msg   string
}

// Set sets the partial fields of the entry with the given lvl and msg
// and sets the log time to the current time by default.
func (e *Entry) Set(lvl Level, msg string) {
	e.Level = lvl
	e.Time = time.Now()
	e.Msg = msg
}
