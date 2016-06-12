package golog

import (
	"errors"
	"github.com/getsentry/raven-go"
)

type SentryHandler struct {
	sentry *raven.Client
}


func NewSentryHandler(dsn string) (*SentryHandler, error) {
	c, err := raven.NewClient(dsn, nil)
	if err != nil {
		return nil, err
	}
	h := new(SentryHandler)
	h.sentry = c
	return h, nil
}

func (h *SentryHandler) Write(b []byte, level int) (n int, err error) {
	
	if level >= LevelWarn {
		str := string(b)  // confirm : no copy : change slice to string 
		packet := raven.NewPacket(str, nil, raven.NewException(errors.New(str), raven.NewStacktrace(3, 3, []string{})))

		switch level {
		case LevelWarn:
			packet.Level = raven.WARNING
		case LevelError:
			packet.Level = raven.ERROR
		case LevelFatal:
			packet.Level = raven.FATAL
		}

		_, ch := h.sentry.Capture(packet, nil)
		if level == LevelFatal {
			<-ch
		}
	}
	// @TODO WARN : return n, error : async send , 
	return len(b), nil
}


func (h *SentryHandler) Close() error {
	h.sentry.Close()
	return nil
}


