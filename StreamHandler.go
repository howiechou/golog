package golog

import (
	"errors"
	"io"
)

var (
	ErrNoWriter = errors.New("io.Writer does not exist!")
)

//StreamHandler writes logs to a specified io Writer, maybe stdout, stderr, etc...
type StreamHandler struct {
	ws []io.Writer
}


func NewStreamHandler(w io.Writer) (*StreamHandler, error) {
	h := new(StreamHandler)
	h.ws = make([]io.Writer, getLastLevel()+1)
	if (w != nil) {
		for i := 0; i < len(h.ws); i++ {
			h.ws[i] = w
		}
	}
	return h, nil
}

// SetLevelWriter(level int, w io.Writer) :
// level : set writer from level to last LEVEL 
// if you want to output LevelWarn,LevelError and LevelFatal to std.err
// 	and output LevelInfo to std.out
// 	you need to code SetLevelWriter(LevelInfo, os.Stdout) first, then SetLevelWriter(LevelWarn, os.Stderr).
func (h *StreamHandler) SetLevelWriter(level int, w io.Writer) {
	if w == nil {
		return
	}
	for i := 0; i <= level; i++ {
		h.ws[i] = w
	}
}


func (h *StreamHandler) Write(b []byte, level int) (n int, err error) {
	return h.ws[level].Write(b)
}

func (h *StreamHandler) Close() error {
	return nil
}












