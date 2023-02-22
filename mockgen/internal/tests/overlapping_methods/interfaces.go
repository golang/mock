package overlap

type ReadCloser interface {
	Read([]byte) (int, error)
	Close() error
}

type WriteCloser interface {
	Write([]byte) (int, error)
	Close() error
}
