package golog

import (
	"testing"
	"time"
	"bytes"
	"strings"
	"os"
	"github.com/getsentry/raven-go"
)

type TestHandler struct {
	bytes.Buffer
}

func (h *TestHandler) Write(b []byte, level int) (n int, err error) {
	//fmt.Printf("testhandler : %s\n", string(b))
	return h.Buffer.Write(b)
}

func (h *TestHandler) Close() error {
	return nil
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}

func TestTestHandler(t *testing.T) {
	h := new(TestHandler)
	h.Write([]byte("hello"), 1)
	if h.String() != "hello" {
		t.Error("hello does not exist")
	}
	t.Log(h.String())
}

func TestLevel(t *testing.T) {
	h := new(TestHandler)
	AppendHandler(h)
	logger.level = LevelTrace

	for i := LevelTrace; i <= LevelFatal; i++ {
		level := LevelName[i]
		logData := "test interface"
		h.Reset()
		
		switch i {
			case LevelTrace:
				Tracef("%s", logData)
			case LevelDebug:
				Debugf("%s", logData)
			case LevelInfo:
				Infof("%s", logData)
			case LevelWarn:
				Warnf("%s", logData)
			case LevelError:
				Errorf("%s", logData)
			case LevelFatal:
				Fatalf("%s", logData)
		}
		time.Sleep(100 * time.Millisecond)

		if !contains(h.Buffer.String(), level) {
			t.Log(h.String())
			t.Errorf("has not %s", level)
		}
		if !contains(h.String(), "test") {
			t.Error("has not test string")
		}
	}
}

func TestHandlers(t *testing.T) {
	ClearHandlers()
	
	h,_ := NewStreamHandler(os.Stdout)
	AppendHandler(h)
	h2,_ := NewSocketHandler("tcp", "localhost:8989", time.Second)
	AppendHandler(h2)
	
	c, err := raven.NewClient("xxx", nil)
	if err != nil {
		t.Error("raven.NewClient fail:" + err.Error())
	}
	h3, _ := NewSentryHandler(c)
	AppendHandler(h3)
	
	
	Warnf("%v\n", "Warnf")
	Infof("%v\n", "Infof")
	Debugf("%v\n", "Debugf")
	Tracef("%v\n", "Tracef")
	Fatalf("%v\n", "Fatalf")
	
	Close()
}

func TestSizeRotatingFileLog(t *testing.T) {
	
	path := "/tmp/testlogrotating/test_log"
	os.RemoveAll(path)

	os.Mkdir(path, 0777)
	fileName := path + "/test"

	h, err := NewSizeRotateFileHandler(fileName, 10, 2)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 10)

	h.Write(buf, LevelDebug)

	h.Write(buf, LevelDebug)

	if _, err := os.Stat(fileName + ".1"); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(fileName + ".2"); err == nil {
		t.Fatal(err)
	}

	h.Write(buf, LevelDebug)
	if _, err := os.Stat(fileName + ".2"); err != nil {
		t.Fatal(err)
	}

	h.Close()

	os.RemoveAll(path)
}

func TestTimeRotatingFileLog(t *testing.T) {
	
}