package golog

import (
	"encoding/binary"
	"net"
	"time"
)


// SocketHandler writes log to peer.
// Protocol is simple:  length(log) + log |  length(log) + log. 
// Log length is uint32. 

type SocketHandler struct {
	c        net.Conn
	protocol string
	addr     string
	timeout time.Duration
}

func NewSocketHandler(protocol string, addr string, timeout time.Duration) (*SocketHandler, error) {
	s := new(SocketHandler)

	s.protocol = protocol
	s.addr = addr
	s.timeout = timeout

	return s, nil
}

func (h *SocketHandler) Write(p []byte, level int) (n int, err error) {
	if err = h.connect(); err != nil {
		return
	}

	buf := make([]byte, len(p)+4)
	// network is bigendian
	binary.BigEndian.PutUint32(buf, uint32(len(p)))

	copy(buf[4:], p)

	n, err = h.c.Write(buf)
	if err != nil {
		Close()
	}
	return
}

func (h *SocketHandler) Close() error {
	if h.c != nil {
		h.c.Close()
		h.c = nil
	}
	return nil
}

func (h *SocketHandler) connect() error {
	if h.c != nil {
		return nil
	}

	var err error
	h.c, err = net.DialTimeout(h.protocol, h.addr, h.timeout)
	if err != nil {
		return err
	}

	return nil
}



















