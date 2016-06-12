# golog

------
A new log package for go

## There are several handlers : 
> * StreamHandler : write log to io.Writer interface  
> * SocketHandler : write log to socket 
> * RotatingFileHandler : write log to file which rotating as size
> * TimeRotatingFileHandler : write log to file which rotating as time(second, minute, hour, day)
> * SentryHandler : write log to sentry(Sentry is a modern error logging and aggregation platform.)

## coding
```
    // clear all handlers already exist
	ClearHandlers()
	
	// new a stream handler : write log to stdout
	h,_ := NewStreamHandler(os.Stdout)
	// append the handler to golog
	AppendHandler(h)
	
	// new a socket handler : it is tcp protocol, address is localhost:8989 and timeout is second
	h2,_ := NewSocketHandler("tcp", "localhost:8989", time.Second)
	AppendHandler(h2)
	
	// new a sentry handler : write log to sentry 
	c, err := raven.NewClient("xxx", nil)
	if err != nil {
		t.Error("raven.NewClient fail:" + err.Error())
	}
	h3, _ := NewSentryHandler(c)
	AppendHandler(h3)
	
	// write log to fileName
	// size of rotating is 10 bytes
	// count is 2 : max count of rotating
	h4, _ := NewSizeRotateFileHandler(fileName, 10, 2)
	AppendHandler(h4)
	
	// the log will be written to h1, h2, h3, h4
	
	Warnf("%v\n", "Warnf")
	Warnln("Warnf")
	Warnf("%v\n", "Warnf")
	Debugf("%v\n", "Debugf")
	Tracef("%v\n", "Tracef")
	Fatalf("%v\n", "Fatalf")
	
	Close()
```

### level
```
//value:
LevelTrace 
LevelDebug
LevelInfo
LevelWarn
LevelError
LevelFatal
```

This is severity of log
If you set level is LevelWarn
LevelWarn, LevelError, LevelFatal are available
LevelTrace, LevelDebug, LevelInfo will be ignore

you can set level like this:
```
golog.SetLevel(golog.LevelWarn)

```

### flag
log header format:
```
//value
Ltime  	// 2006/01/02 15:04:05
Lfile	// file.go[123]
Llevel	// [Trace]
```
you can set flag like this:
```
golog.SetFlag(golog.Ltime | golog.Lfile)
```
