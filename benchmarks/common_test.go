package benchmarks

import (
	"time"
	"errors"
	"fmt"
)

var (

	_fakeMsg       = "message for benchmark"
	_fakeMsgFormat = "formatted message for benchmarks: %v"
	_fakeArg       = "happyxcj/wlog"

	//_fakeStdMsgformatWithFields = intiFakeStdMsgFormatWithFields(12)
	//
	//_messages     = initFakeMessages(10)
	_fakeInts     = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	_fakeFloat64s = []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	_fakeStrings  = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	_fakeDurations = []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second,
		5 * time.Second, 6 * time.Second, 7 * time.Second, 8 * time.Second, 9 * time.Second, 10 * time.Second}

	_fakeTimes = []time.Time{time.Now(), time.Now(), time.Now(), time.Now(), time.Now(),
		time.Now(), time.Now(), time.Now(), time.Now(), time.Now()}

	_fakeError  = errors.New("fake error")
	_fakeErrors = []error{_fakeError, _fakeError, _fakeError, _fakeError,
		_fakeError, _fakeError, _fakeError, _fakeError, _fakeError, _fakeError}
)

