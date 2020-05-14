package wlog

import (
	"os"
	"time"
	"fmt"
	"strings"
	"path/filepath"
	"sort"
	"io"
)

const (
	baseFileNameFormat        = "%v_%v%v"
	dailySizeFileNameFormat   = "%v_%v-%v%v"
	dayLayout                 = "20060102"
	defaultFileExtension      = ".log"
	defaultFileMaxSize        = 100 * 1 << 20
	defaultFileMaxRotatedSize = 10 * 1 << 30
	defaultFileMaxRotatedDays = 100
)

type FileWriter struct {
	file *os.File
	// fileName is the file name.
	fileName string
	// basicName is the basic name of the file.
	basicName string
	// extension is the extension of the file.
	// It's default value is "log".
	extension string

	// maxSize is the size threshold for log rotation.
	// It's default value is "100*1<<20".
	maxSize int64
	// currSize is the current size of the file.
	currSize int64
	// currRotatedSize is the current size of all rotated files.
	// It's default value is "10*1<<30".
	currRotatedSize int64
	// maxRotatedSize is the maximum size of all rotated files.
	maxRotatedSize int64

	// isDaily is a tag indicates whether rotates the file daily.
	isDaily bool
	// dayRotateAt is the unix time in seconds for next day rotation.
	dayRotateAt int64
	// dayRotatedCount is the count of files that have been rotated today.
	// It will be reset if the "isDaily" is set to true.
	dayRotatedCount int
	// currRotatedDays is the current days of all rotated files.
	currRotatedDays int
	// maxRotatedDays is the maximum days of all rotated files.
	// It's default value is "100".
	maxRotatedDays int

	// rotatedFiles records all rotated file names.
	// It records files in ascending order of time.
	rotatedFiles []string
	// errW is used to output the interval error.
	// os.Stderr is the default io writer.
	errW io.Writer
}

type FileWriterOpt func(w *FileWriter)

// SetFileMaxSize sets the maximum file size of the FileWriter.
func SetFileMaxSize(maxSize int64) FileWriterOpt {
	return func(w *FileWriter) {
		if maxSize <= 0 {
			w.maxSize = defaultFileMaxSize
			return
		}
		w.maxSize = maxSize
	}
}

// SetFileMaxSize sets the the maximum rotated size of the FileWriter.
func SetFileMaxRotatedSize(maxRotatedSize int64) FileWriterOpt {
	return func(w *FileWriter) {
		if maxRotatedSize <= 0 {
			w.maxRotatedSize = defaultFileMaxRotatedSize
			return
		}
		w.maxRotatedSize = maxRotatedSize
	}
}

// SetFileMaxSize sets the the maximum rotated days of the FileWriter.
func SetFileMaxRotatedDays(maxRotatedDays int) FileWriterOpt {
	return func(w *FileWriter) {
		if maxRotatedDays <= 0 {
			w.maxRotatedDays = defaultFileMaxRotatedDays
			return
		}
		w.maxRotatedDays = maxRotatedDays
	}
}

// SetFileMaxSize sets the "isDaily" of the FileWriter to false.
func DisableFileDaily() FileWriterOpt {
	return func(w *FileWriter) {
		w.isDaily = true
	}
}

// SetFileErrW sets the underlying io.Writer of the FileWriter to output the internal error.
func SetFileErrW(errW io.Writer) FileWriterOpt {
	return func(w *FileWriter) {
		w.errW = errW
	}
}

func NewFileWriter(fileName string, opts ...FileWriterOpt) *FileWriter {
	w := &FileWriter{
		fileName:       fileName,
		maxSize:        defaultFileMaxSize,
		maxRotatedSize: defaultFileMaxRotatedSize,
		isDaily:        true,
		maxRotatedDays: defaultFileMaxRotatedDays,
		rotatedFiles:   make([]string, 0, 20),
		errW:           os.Stderr,
	}
	index := strings.LastIndex(fileName, ".")
	if index == -1 {
		w.basicName = fileName
		w.extension = defaultFileExtension
	} else {
		w.basicName = fileName[:index]
		w.extension = fileName[index:]
	}
	for _, opt := range opts {
		opt(w)
	}
	w.resetRotateAt()
	w.resetRotatedFiles()
	err := w.resetCurrFile()
	if err != nil {
		panic(fmt.Sprintf("unable to open the file: %v, error: %v", w.fileName, err))
	}
	return w
}

func (w *FileWriter) resetRotateAt() {
	next := time.Now().Add(24 * time.Hour)
	rotateAt := time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
	w.dayRotateAt = rotateAt.Unix()
}

func (w *FileWriter) resetRotatedFiles() {
	matchNames, _ := filepath.Glob(fmt.Sprint(w.basicName, "_*"))
	n := len(matchNames)
	if n == 0 {
		return
	}
	//sort file name by mode time
	sort.SliceStable(matchNames, func(i, j int) bool {
		leftInfo, _ := os.Stat(matchNames[i])
		rightInfo, _ := os.Stat(matchNames[j])
		return leftInfo.ModTime().Sub(rightInfo.ModTime()) < 0
	})
	// Cache all rotated files.
	w.rotatedFiles = matchNames
	if !w.isDaily {
		for _, matchName := range matchNames {
			w.currRotatedSize += w.getFileSize(matchName)
		}
		return
	}
	var currRotatedSize int64
	var currRotatedDays, todayCount int
	var lastDay string
	today := time.Now().Format(dayLayout)
	for _, matchName := range matchNames {
		info, err := os.Stat(matchName)
		if err != nil {
			continue
		}
		currRotatedSize += info.Size()
		day := info.ModTime().Format(dayLayout)
		if day == today {
			todayCount++
			continue
		}
		if day != lastDay {
			currRotatedDays++
			lastDay = day
		}
	}
	w.currRotatedSize = currRotatedSize
	w.currRotatedDays = currRotatedDays
	w.dayRotatedCount = todayCount
}

func (w *FileWriter) resetCurrFile() error {
	if w.file != nil {
		w.file.Close()
		w.file = nil
	}
	var err error
	w.file, err = os.OpenFile(w.fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return err
	}
	f, err := w.file.Stat()
	if err != nil {
		return err
	}
	w.currSize = f.Size()
	// Check to start with a new line.
	if w.currSize > 0 {
		w.file.Write([]byte{'\n'})
	}
	return err
}

func (w *FileWriter) Write(bs []byte) (n int, err error) {
	w.doSizeRotation()
	w.doDayRotation()
	n, err = w.file.Write(bs)
	w.currSize += int64(n)
	return
}

func (w *FileWriter) Flush() error {
	return nil
}

func (w *FileWriter) Close() error {
	return w.file.Close()
}

func (w *FileWriter) doSizeRotation() {
	if w.currSize < w.maxSize {
		return
	}
	w.rotateFile()
	if w.currRotatedSize > w.maxRotatedSize {
		w.deleteExpiredSize()
	}
}

func (w *FileWriter) deleteExpiredSize() {
	if len(w.rotatedFiles) == 0 {
		return
	}
	// Regardless of the result of delete operation, update some information.
	deleteName := w.rotatedFiles[0]
	w.rotatedFiles = w.rotatedFiles[1:]
	w.currRotatedSize -= w.getFileSize(deleteName)
	if w.isDaily {
		// The rotating date is the day before the date represented by "dayRotateAt".
		today := time.Unix(w.dayRotateAt, 0).Add(-time.Hour).Format(dayLayout)
		if w.dayRotatedCount > 0 && strings.HasPrefix(deleteName, fmt.Sprint(w.basicName, "_", today)) {
			w.dayRotatedCount -= 1
		}
	}
	w.deleteFiles(deleteName)
}

func (w *FileWriter) doDayRotation() {
	if !w.isDaily {
		return
	}
	now := time.Now()
	if now.Unix() < w.dayRotateAt {
		return
	}
	if w.currSize > 0 {
		w.rotateFile()
	}
	w.resetRotateAt()
	// Reset the next rotated count of files for a new day.
	w.dayRotatedCount = 0
	// Even if the file has no content, the number of rotated days is incremented by 1.
	w.currRotatedDays += 1
	if w.currRotatedDays > w.maxRotatedDays {
		w.deleteExpiredDay()
	}
}

func (w *FileWriter) deleteExpiredDay() {
	n := len(w.rotatedFiles)
	if n == 0 {
		return
	}
	oldestName := w.rotatedFiles[0]
	w.currRotatedSize -= w.getFileSize(oldestName)

	var dayPrefix string
	// prefix format for same day: w.basicName_20060102.
	prefixLen := len(w.basicName) + 1 + len(dayLayout)
	if len(oldestName) <= prefixLen {
		dayPrefix = oldestName
	} else {
		dayPrefix = oldestName[:prefixLen]
	}
	var dayNum int
	for dayNum = 1; dayNum < n; dayNum++ {
		name := w.rotatedFiles[dayNum]
		if !strings.HasPrefix(name, dayPrefix) {
			break
		}
		w.currRotatedSize -= w.getFileSize(name)
	}
	deleteNames := w.rotatedFiles[:dayNum]
	w.rotatedFiles = w.rotatedFiles[dayNum:]
	w.deleteFiles(deleteNames...)
}

func (w *FileWriter) rotateFile() {
	n := len(w.rotatedFiles)
	var newPath string
	if !w.isDaily {
		for i := 0; i < n; i++ {
			newPath = fmt.Sprintf(baseFileNameFormat, w.basicName, n-i+1, w.extension)
			os.Rename(w.rotatedFiles[i], newPath)
			w.rotatedFiles[i] = newPath
		}
		newPath = fmt.Sprintf(baseFileNameFormat, w.basicName, 1, w.extension)
	} else {
		start := n - w.dayRotatedCount
		if start < 0 {
			start = 0
		}
		// The rotating date is the day before the date represented by "dayRotateAt".
		day := time.Unix(w.dayRotateAt, 0).Add(-time.Hour).Format(dayLayout)
		for i := start; i < n; i++ {
			nextCode := w.dayRotatedCount - (i - start)
			newPath = fmt.Sprintf(dailySizeFileNameFormat, w.basicName, day, nextCode, w.extension)
			os.Rename(w.rotatedFiles[i], newPath)
			w.rotatedFiles[i] = newPath
		}
		newPath = fmt.Sprintf(baseFileNameFormat, w.basicName, day, w.extension)
		// update the rotated count of files today.
		w.dayRotatedCount++
	}
	// Close current file before renaming the file.
	w.file.Close()
	os.Rename(w.fileName, newPath)
	// Regardless of the result of rotation for current file, update some information.
	w.rotatedFiles = append(w.rotatedFiles, newPath)
	w.currRotatedSize += w.currSize
	w.resetCurrFile()
}

func (w *FileWriter) getFileSize(name string) int64 {
	fileInfo, err := os.Stat(name)
	if err != nil {
		t := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(w.errW, "FileWriter: unable to get file '%v' information at time: %v, error: %v\n", name, t, err)
		return 0
	}
	return fileInfo.Size()
}

func (w *FileWriter) deleteFiles(names ...string) {
	for _, name := range names {
		err := os.Remove(name)
		if err != nil {
			t := time.Now().Format("2006-01-02 15:04:05")
			fmt.Fprintf(w.errW, "FileWriter: unable to delete file '%v' at time: %v, error: %v\n", name, t, err)
		}
	}
}
