package tube

import (
	"io"
	"time"
)

type Tube interface {
	Send(data []byte) error
	Recv(n int) ([]byte, error)
	SendLine(data []byte) error
	RecvLine() ([]byte, error)
	RecvUntil(delim []byte) ([]byte, error)
	RecvAll() ([]byte, error)
	Interactive() error
	Close() error
	SetTimeout(timeout time.Duration)
	IsOpen() bool
}

type Options struct {
	Timeout  time.Duration
	NewLine  []byte
	LogLevel int
}

var DefaultOptions = Options{
	Timeout:  30 * time.Second,
	NewLine:  []byte("\n"),
	LogLevel: 0,
}

type baseTube struct {
	options Options
	reader  io.Reader
	writer  io.Writer
	closed  bool
}

func (t *baseTube) SetTimeout(timeout time.Duration) {
	t.options.Timeout = timeout
}

func (t *baseTube) IsOpen() bool {
	return !t.closed
}
