package golog

// NullHandler does nothing, it discards anything.
type NullHandler struct {
}

func NewNullHandler() (*NullHandler, error) {
	return new(NullHandler), nil
}

func (h *NullHandler) Write(b []byte, level int) (n int, err error) {
	return len(b), nil
}

func (h *NullHandler) Close() error {
	return nil
}




