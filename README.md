# wlog

A fast, highly scalable log library in Golang, it support different log levels: "debug", "info", "earn", "error", "fatal", "panic".

## Features

+ Easy to use (See [Quick Start](#quick-start)).
+ High efficiency (See [Performance](#performance)).
+ Support the file rotation by size and day.
+ Provide many highly scalable interfaces such as Encoder, Writer, Handler and etc.

## Install

```
go get -u github.com/happyxcj/wlog
```

## Quick Start

```go
package main

import "github.com/happyxcj/wlog"

func main() {
	wlog.UseGlobalLogger(nil)
	defer wlog.Close()
	wlog.Infow("failed to login account",
		wlog.String("username", "root"),
		wlog.String("password", "root"))
}
```

## Performance

The follow benchmarks are mainly to compare the performance of different log libraries. the log libraries are as  follow: 

​		github.com/happyxcj/wlog

​		go.uber.org/zap

​		standard log (fmt.Println)

​		github.com/Sirupsen/logrus

For more details and more benchmarks , see internal [benchmarks](https://github.com/happyxcj/wlog/blob/master/benchmarks).



Benchmark of logging a message with 12 fields:

```
goos: windows
goarch: amd64
BenchmarkWithFields/wlog-4               1000000   1320 ns/op  816 B/op   32 allocs/op
BenchmarkWithFields/uber-go/zap-4        1000000   2187 ns/op  985 B/op   18 allocs/op
BenchmarkWithFields/fmt.Println-4         300000   4906 ns/op  3785 B/op  129 allocs/op
BenchmarkWithFields/Sirupsen/logrus-4     200000   6921 ns/op  9201 B/op  99 allocs/op
```



Benchmark of logging a message with 12 key-value pairs:

```
goos: windows
goarch: amd64
BenchmarkWithPairs/wlog-4                1000000   1875 ns/op  1562 B/op  54 allocs/op
BenchmarkWithPairs/uber-go/zap-4          500000   3017 ns/op  3325 B/op  37 allocs/op
BenchmarkWithPairs/fmt.Println-4 	      300000   5006 ns/op  3785 B/op  129 allocs/op
BenchmarkWithPairs/Sirupsen/logrus-4      200000   6891 ns/op  9202 B/op  99 allocs/op
```



## Examples

Create a Logger with any expected Encoder, Writer and Handler.

Note that If logging a string message with any fields, method  Infow(string, ...Field) is a better choice than method With(...Field).Info(...Interface{}).

```go
package main

import "github.com/happyxcj/wlog"

func main() {
	w := wlog.NewIOWriter(os.Stdout)
	bw := wlog.NewBufWriter(w, wlog.SetBufMaxSize(50*1<<20))
	fw := wlog.NewTimingFlushWriter(bw, 5*time.Second)
	h := wlog.NewBaseHandler(fw, wlog.NewTextEncoder(wlog.DisableTime()))
	logger := wlog.NewLogger(h)
	defer logger.Close()
    
   logger.With(wlog.String("username", "root"),
		wlog.String("password", "root"),
		wlog.Int("retry", 5)).
		Info("failed to login account")
    
	logger.Infow("failed to login account",
		wlog.String("username", "root"),
		wlog.String("password", "root"),
		wlog.Int("retry", 5))
}

```



For more simple usage, see the [test file](https://github.com/happyxcj/wlog/blob/master/logger_example_test.log)

For more advanced usage such as to configure file rotation, configure multiple loggers and etc, see the [examples](https://github.com/happyxcj/wlog/blob/master/examples)



## References

- github.com/Sirupsen/logrus
- go.uber.org/zap