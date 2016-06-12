package golog

import (
	"sync"
	"bytes"
)

const (
	maxBufPoolSize = 32
)

type buffer struct {
	bytes.Buffer
	level int
}


type buffers struct {
	sync.Mutex
	bufs []*buffer
}

func (l *buffers) getBuffer() *buffer {
	l.Lock()
	var buf *buffer
	if len(l.bufs) == 0 {
		buf = new(buffer)
	} else {
		buf = l.bufs[len(l.bufs) - 1]
		l.bufs = l.bufs[0 : len(l.bufs) - 1]
	}
	l.Unlock()

	return buf
}

func (l *buffers) putBuffer(b *buffer) {
	if b.Len() >= 256 {
		// Let big buffer die a natrual death.
		return
	}
	l.Lock()
	if len(l.bufs) < maxBufPoolSize {
		b.Reset()
		l.bufs = append(l.bufs, b)
	}
	l.Unlock()
}





