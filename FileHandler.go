package golog

import (
	"fmt"
	"os"
	"path"
	"time"
)

/*	FileHandler : @TODO : write interface is only cached by OS, add a mechanism timer to flush
 *	SizeRotateFileHandler
 *	TimeRotateFileHandler
 *	LevelRotateFileHandler @TODO : refer : github.com/golang/glog. split file to sevaral as LEVEL
 */

//  write log to a file.
type FileHandler struct {
	fd *os.File
}

func NewFileHandler(fileName string, fileFlag int) (*FileHandler, error) {
	dir := path.Dir(fileName)
	os.Mkdir(dir, 0777) // full permission except suid,guid

	f, err := os.OpenFile(fileName, fileFlag, 0)
	if err != nil {
		return nil, err
	}
	h := new(FileHandler)
	h.fd = f
	return h, nil
}

func (h *FileHandler) Write(b []byte, level int) (n int, err error) {
	return h.fd.Write(b)
}

func (h *FileHandler) Close() error {
	return h.fd.Close()
}

// ####################   SizeRotateFileHandler
// write log a file
type SizeRotateFileHandler struct {
	fd *os.File

	fileName    string
	maxBytes    int64
	backupCount int
}

func NewSizeRotateFileHandler(fileName string, maxBytes int64, backupCount int) (*SizeRotateFileHandler, error) {
	dir := path.Dir(fileName)
	os.Mkdir(dir, 0777)

	h := new(SizeRotateFileHandler)

	if maxBytes <= 0 {
		return nil, fmt.Errorf("invalid max bytes")
	}

	h.fileName = fileName
	h.maxBytes = maxBytes
	h.backupCount = backupCount

	var err error
	h.fd, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *SizeRotateFileHandler) Write(p []byte, level int) (n int, err error) {
	h.doRollover()
	return h.fd.Write(p)
}

func (h *SizeRotateFileHandler) Close() error {
	if h.fd != nil {
		return h.fd.Close()
	}
	return nil
}

func (h *SizeRotateFileHandler) doRollover() {
	f, err := h.fd.Stat()
	if err != nil {
		return
	}

	if h.maxBytes <= 0 {
		return
	} else if f.Size() < int64(h.maxBytes) {
		return
	}

	if h.backupCount > 0 {
		h.fd.Close()

		for i := h.backupCount - 1; i > 0; i-- {
			sfn := fmt.Sprintf("%s.%d", h.fileName, i)
			dfn := fmt.Sprintf("%s.%d", h.fileName, i+1)

			os.Rename(sfn, dfn)
		}

		dfn := fmt.Sprintf("%s.1", h.fileName)
		os.Rename(h.fileName, dfn)

		h.fd, _ = os.OpenFile(h.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}
}

// ###################		TimeRotateFileHandler
// write log to a file,
type TimeRotateFileHandler struct {
	fd *os.File

	baseName   string
	interval   int64
	suffix     string
	rolloverAt int64
}

const (
	WhenSecond = iota
	WhenMinute
	WhenHour
	WhenDay
)

func NewTimeRotateFileHandler(baseName string, when int8, interval int) (*TimeRotateFileHandler, error) {
	dir := path.Dir(baseName)
	os.Mkdir(dir, 0777)

	h := new(TimeRotateFileHandler)

	h.baseName = baseName

	switch when {
	case WhenSecond:
		h.interval = 1
		h.suffix = "2006-01-02_15-04-05"
	case WhenMinute:
		h.interval = 60
		h.suffix = "2006-01-02_15-04"
	case WhenHour:
		h.interval = 3600
		h.suffix = "2006-01-02_15"
	case WhenDay:
		h.interval = 3600 * 24
		h.suffix = "2006-01-02"
	default:
		return nil, fmt.Errorf("invalid when_rotate: %d", when)
	}

	h.interval = h.interval * int64(interval)

	var err error
	h.fd, err = os.OpenFile(h.baseName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	fInfo, _ := h.fd.Stat()
	h.rolloverAt = fInfo.ModTime().Unix() + h.interval

	return h, nil
}

func (h *TimeRotateFileHandler) doRollover() {
	now := time.Now()

	if h.rolloverAt <= now.Unix() {
		fName := h.baseName + now.Format(h.suffix)
		h.fd.Close()
		e := os.Rename(h.baseName, fName)
		if e != nil {
			panic(e)
		}

		h.fd, _ = os.OpenFile(h.baseName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

		h.rolloverAt = time.Now().Unix() + h.interval
	}
}

func (h *TimeRotateFileHandler) Write(b []byte, level int) (n int, err error) {
	h.doRollover()
	return h.fd.Write(b)
}

func (h *TimeRotateFileHandler) Close() error {
	return h.fd.Close()
}

// ######################   LevelFileHandler

type LevelFileHandler struct {
}
