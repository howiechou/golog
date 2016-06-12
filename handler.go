package golog



// Handler writes logs to somewhere you want
type Handler interface {
	Write(p []byte, level int) (n int, err error)
	Close() error
}






















