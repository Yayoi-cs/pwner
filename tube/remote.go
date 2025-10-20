package tube

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Remoter struct {
	baseTube
	conn net.Conn
	host string
	port int
}

func Remote(host string, port int, opts ...Options) *Remoter {
	options := DefaultOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, options.Timeout)
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", address, err)
	}

	r := &Remoter{
		baseTube: baseTube{
			options: options,
			reader:  conn,
			writer:  conn,
			closed:  false,
		},
		conn: conn,
		host: host,
		port: port,
	}

	return r
}

func (r *Remoter) Send(data []byte) {
	if r.closed {
		log.Fatalf("broken tube")
		return
	}

	if r.options.Timeout > 0 {
		r.conn.SetWriteDeadline(time.Now().Add(r.options.Timeout))
	}

	_, err := r.conn.Write(data)
	if err != nil {
		log.Fatalf("failed to send data: %v", err)
	}
	return
}

func (r *Remoter) SendLine(data []byte) {
	r.Send(append(data, r.options.NewLine...))
	return
}
func (r *Remoter) Recv(n ...int) []byte {
	if r.closed {
		log.Fatalf("broken tube")
	}

	if r.options.Timeout > 0 {
		r.conn.SetReadDeadline(time.Now().Add(r.options.Timeout))
	}

	if len(n) > 0 && n[0] > 0 {
		buf := make([]byte, n[0])
		total := 0
		for total < n[0] {
			read, err := r.conn.Read(buf[total:])
			if err != nil {
				if total > 0 {
					return buf[:total]
				}
				log.Fatalf("recv error: %v", err)
			}
			total += read
		}
		return buf
	}

	buf := make([]byte, 0x1000)
	nRead, err := r.conn.Read(buf)
	if err != nil {
		if nRead > 0 {
			return buf[:nRead]
		}
		log.Fatalf("recv error: %v", err)
	}
	return buf[:nRead]
}

func (r *Remoter) RecvLine() []byte {
	if r.closed {
		log.Fatalf("tube is closed")
	}

	if r.options.Timeout > 0 {
		r.conn.SetReadDeadline(time.Now().Add(r.options.Timeout))
	}

	var buf bytes.Buffer
	b := make([]byte, 1)
	for {
		_, err := r.conn.Read(b)
		if err != nil {
			log.Fatalf("recv error: %v", err)
		}
		buf.Write(b)
		if bytes.HasSuffix(buf.Bytes(), r.options.NewLine) {
			result := buf.Bytes()
			return result[:len(result)-len(r.options.NewLine)]
		}
	}
}

func (r *Remoter) RecvUntil(delim []byte) []byte {
	if r.closed {
		log.Fatalf("tube is closed")
	}

	if r.options.Timeout > 0 {
		r.conn.SetReadDeadline(time.Now().Add(r.options.Timeout))
	}

	var buf bytes.Buffer
	b := make([]byte, 1)
	for {
		_, err := r.conn.Read(b)
		if err != nil {
			log.Fatalf("recv error: %v", err)
		}
		buf.Write(b)
		if bytes.HasSuffix(buf.Bytes(), delim) {
			return buf.Bytes()
		}
	}
}

func (r *Remoter) RecvAll() []byte {
	if r.closed {
		log.Fatalf("tube is closed")
	}

	r.conn.SetReadDeadline(time.Time{})

	ret, err := io.ReadAll(r.conn)
	if err != nil {
		log.Fatalf("recv error: %v", err)
	}
	return ret
}

func (r *Remoter) Interactive() {
	if r.closed {
		log.Fatalf("tube is closed")
	}

	r.conn.SetReadDeadline(time.Time{})

	go io.Copy(os.Stdout, r.conn)

	io.Copy(r.conn, os.Stdin)
}

func (r *Remoter) Close() {
	if r.closed {
		log.Fatalf("broken tube")
	}
	r.closed = true
	err := r.conn.Close()
	if err != nil {
		log.Fatalf("close error: %v", err)
	}
}

func (r *Remoter) GetAddress() string {
	return fmt.Sprintf("%s:%d", r.host, r.port)
}
