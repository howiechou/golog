package golog

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"strings"
)

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func getLastLevel() int {
	return LevelFatal
}


const (
	Ltime  = 1 << iota	// 2006/01/02 15:04:05
	Lfile				// file.go[123]
	Llevel			// [Trace]
)

const TimeFormat = "2006/01/02 15:04:05 "

var LevelName [6]string = [6]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}



type Logger struct {
	level int
	flag int
	quit chan struct{}
	msg chan *buffer
	bufs buffers
	wg sync.WaitGroup
	closed bool
	handlers []Handler
}

var logger *Logger

func init() {
	logger = new(Logger)
	
	logger.level = LevelInfo
	logger.flag = Ltime | Lfile | Llevel
	
	logger.quit = make(chan struct{})
	logger.closed = false
	
	logger.msg = make(chan *buffer, 8)
	
	logger.wg.Add(1)
	go logger.start()
}

func ClearHandlers() {
	for _, h := range logger.handlers {
		if h != nil {
			h.Close()
		}
	}
	logger.handlers = logger.handlers[:0]
}

func AppendHandler(h Handler) {
	if h != nil {
		logger.handlers = append(logger.handlers, h)
	}
}


func (l *Logger)start() {
	defer l.wg.Done()
	for {
		select {
			case msg := <- l.msg:
				for _, h := range l.handlers {
					if h != nil {
						h.Write(msg.Bytes(), msg.level)
					}
				}
				l.bufs.putBuffer(msg)
			case <- l.quit:
				if len(l.msg) == 0 {
					// maybe msg chan is not empty, continue until msg chan is empty
					return
				}
		}
	}
}

// close and flush 
func Close() {
	if logger.closed {
		return
	}
	logger.closed = true

	close(logger.quit)
	logger.wg.Wait()
	logger.quit = nil

	ClearHandlers()
}




func (l *Logger) headers(level int, depth int) (*buffer) {
	_, file, line, ok := runtime.Caller(3 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}

	buf := l.bufs.getBuffer()

	if l.flag & Ltime > 0 {
		now := time.Now().Format(TimeFormat)
		buf.WriteString(now)
	}
	if l.flag & Llevel > 0 {
		buf.WriteString(LevelName[level] + " ")
	}
	if l.flag * Lfile > 0 {
		fmt.Fprintf(buf, "%s[%d] : ", file, line)
	}

	return buf
}


func (l *Logger) isValidLevel(level int) bool {
	return (level >= l.level && level <= LevelFatal)
}

func (l *Logger)printf(level int, format string, args ...interface{}) {
	if !l.isValidLevel(level) || len(l.handlers) == 0 {
		return
	}
	buf := l.headers(level, 0)
	fmt.Fprintf(buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.msg <- buf
}

func (l *Logger)println(level int, args ...interface{}) {
	if !l.isValidLevel(level)  || len(l.handlers) == 0{
		return
	}
	buf := l.headers(level, 0)
	fmt.Fprintln(buf, args...)
	l.msg <- buf
}


func Traceln(v ...interface{}) {
	logger.println(LevelTrace, v...)
}

func Debugln(v ...interface{}) {
	logger.println(LevelDebug, v...)
}

func Infoln(v ...interface{}) {
	logger.println(LevelInfo, v...)
}

func Warnln(v ...interface{}) {
	logger.println(LevelWarn, v...)
}

func Errorln(v ...interface{}) {
	logger.println(LevelError, v...)
}

func Fatalln(v ...interface{}) {
	logger.println(LevelFatal, v...)
}

//func Trace(arg0 interface{}, args ...interface{}) {
//	switch arg0.(type) {
//	case string:
//		logger.printf(LevelTrace, arg0.(string), args...)
//	default:
//		logger.println(LevelTrace, []interface{}{arg0, args}...)
//	}
//}

//func Debug(arg0 interface{}, args ...interface{}) {
//	switch arg0.(type) {
//	case string:
//		logger.printf(LevelDebug, arg0.(string), args...)
//	default:
//		logger.println(LevelDebug, []interface{}{arg0, args}...)
//	}
//}
//func Info(arg0 interface{}, args ...interface{}) {
//	switch arg0.(type) {
//	case string:
//		logger.printf(LevelInfo, arg0.(string), args...)
//	default:
//		logger.println(LevelInfo, []interface{}{arg0, args}...)
//	}
//}
//func Warn(arg0 interface{}, args ...interface{}) {
//	switch arg0.(type) {
//	case string:
//		logger.printf(LevelWarn, arg0.(string), args...)
//	default:
//		logger.println(LevelWarn, []interface{}{arg0, args}...)
//	}
//}
//func Error(arg0 interface{}, args ...interface{}) {
//	switch arg0.(type) {
//	case string:
//		logger.printf(LevelError, arg0.(string), args...)
//	default:
//		logger.println(LevelError, []interface{}{arg0, args}...)
//	}
//}
//func Fatal(arg0 interface{}, args ...interface{}) {
//	switch arg0.(type) {
//	case string:
//		logger.printf(LevelFatal, arg0.(string), args...)
//	default:
//		logger.println(LevelFatal, []interface{}{arg0, args}...)
//	}
//}


func Tracef(format string, v ...interface{}) {
	logger.printf(LevelTrace, format, v...)
}

func Debugf(format string, v ...interface{}) {
	logger.printf(LevelDebug, format, v...)
}

func Infof(format string, v ...interface{}) {
	logger.printf(LevelInfo, format, v...)
}

func Warnf(format string, v ...interface{}) {
	logger.printf(LevelWarn, format, v...)
}

func Errorf(format string, v ...interface{}) {
	logger.printf(LevelError, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	logger.printf(LevelFatal, format, v...)
}




func SetLevel(level int) {
	if (level >= LevelTrace && level <= LevelFatal) {
		logger.level = level
	}
}

func GetLevel() int {
	return logger.level
}

func SetFlag(flag int) {
	logger.flag = flag
}
func GetFlag() int {
	return logger.flag
}

