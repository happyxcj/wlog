package wlog_test

import (
	//"github.com/happyxcj/wlog"
	//"os"
	//"time"
	"github.com/happyxcj/wlog"
	"os"
	"time"
	"io/ioutil"
)

func createTimeDisabledLogger() *wlog.Logger {
	c := wlog.NewConsoleConfig()
	c.EncoderConfig.TimeDisabled = true
	logger, _ := c.Create()
	return logger
}

func ExampleConfigCreateLogger() {
	logger := createTimeDisabledLogger()
	defer logger.Close()
	logger.Infow("failed to login account",
		wlog.String("username", "root"),
		wlog.String("password", "root"),
		wlog.Int("retry", 5))
	// Output:
	// [INFO]  failed to login account  [username=root password=root retry=5]
}

func ExampleNewLogger() {
	w := wlog.NewIOWriter(os.Stdout)
	bw := wlog.NewBufWriter(w, wlog.SetBufMaxSize(50*1<<20))
	fw := wlog.NewTimingFlushWriter(bw, 5*time.Second)
	h := wlog.NewBaseHandler(fw, wlog.NewTextEncoder(wlog.DisableTime()))
	logger := wlog.NewLogger(h)
	defer logger.Close()
	logger.Infop("failed to login account",
		"username", "root",
		"password", "root",
		"retry", 5)
	// Output:
	// [INFO]  failed to login account  [username=root password=root retry=5]
}

func ExampleMultiWriterLogger() {
	c := wlog.NewMultiWriterConfig("stdout", "stdout")
	c.EncoderConfig.TimeDisabled = true
	logger, _ := c.Create()
	defer logger.Close()
	logger.Infow("failed to login account",
		wlog.String("username", "root"),
		wlog.String("password", "root"),
		wlog.Int("retry", 5))
	// Output:
	// [INFO]  failed to login account  [username=root password=root retry=5]
	// [INFO]  failed to login account  [username=root password=root retry=5]
}

func ExampleUseGlobalLogger() {
	wlog.UseGlobalLogger(createTimeDisabledLogger())
	defer wlog.Close()
	wlog.Infow("failed to login account",
		wlog.String("username", "root"),
		wlog.String("password", "root"),
		wlog.Int("retry", 5))
	// Output:
	// [INFO]  failed to login account  [username=root password=root retry=5]
}

func ExampleWithFirst() {
	wlog.UseGlobalLogger(createTimeDisabledLogger())
	defer wlog.Close()
	wlog.Infow("no fields")
	wlog.With(wlog.String("belong", "xcj")).Infow("have fields")
	wlog.Infow("still no fields")
	// Output:
	// [INFO]  no fields
	// [INFO]  have fields  [belong=xcj]
	// [INFO]  still no fields
}

func ExampleReplaceMinLvl() {
	wlog.UseGlobalLogger(createTimeDisabledLogger())
	defer wlog.Close()
	wlog.Infow("working")
	wlog.WithOpts(wlog.SetLogMinLvl(wlog.ErrorLvl)).Infow("no working")
	wlog.Infow("still working")
	// Output:
	// [INFO]  working
	// [INFO]  still working
}

func ExampleWrapHandler() {
	wlog.UseGlobalLogger(createTimeDisabledLogger())
	defer wlog.Close()
	wlog.Infow("no fields")
	opt := wlog.WrapLogHandler(func(handler wlog.Handler) wlog.Handler {
		return handler.With(wlog.String("name", "xcj"))
	})
	wlog.WithOpts(opt).Infow("have fields")
	wlog.Infow("still no fields")
	// Output:
	// [INFO]  no fields
	// [INFO]  have fields  [name=xcj]
	// [INFO]  still no fields
}

func ExampleSetHandler() {
	wlog.UseGlobalLogger(createTimeDisabledLogger())
	defer wlog.Close()
	wlog.Infow("working")
	opt := wlog.SetLogHandler(wlog.NewBaseHandler(wlog.NewIOWriter(ioutil.Discard), wlog.NewTextEncoder()))
	wlog.WithOpts(opt).Infow("discard logs")
	wlog.Infow("still working")
	// Output:
	// [INFO]  working
	// [INFO]  still working
}
