package terminal

type Terminal interface {
	Start() error
	Close() error
	Write(p []byte) (n int, err error)
	Read(p []byte) (n int, err error)
	Resize(cols, rows uint32) error
}

type Config struct {
	Cols uint32
	Rows uint32
}
