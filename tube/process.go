package tube

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type Proc struct {
	baseTube
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	scanner *bufio.Scanner
}

func Process(args []string, opts ...Options) (*Proc, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	options := DefaultOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	cmd := exec.Command(args[0], args[1:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start process: %v", err)
	}

	p := &Proc{
		baseTube: baseTube{
			options: options,
			reader:  stdout,
			writer:  stdin,
			closed:  false,
		},
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		scanner: bufio.NewScanner(stdout),
	}

	return p, nil
}

func (p *Proc) Send(data []byte) {
	if p.closed {
		log.Fatalf("broken tube")
		return
	}
	_, err := p.stdin.Write(data)
	if err != nil {
		log.Fatalf("failed to send data: %v", err)
	}
	return
}

func (p *Proc) SendLine(data []byte) {
	p.Send(append(data, p.options.NewLine...))
	return
}

func (p *Proc) Recv(n ...int) []byte {
	if p.closed {
		log.Fatalf("broken tube")
	}

	if len(n) > 0 && n[0] > 0 {
		buf := make([]byte, n[0])
		total := 0
		for total < n[0] {
			read, err := p.stdout.Read(buf[total:])
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

	buf := make([]byte, 4096)
	nRead, err := p.stdout.Read(buf)
	if err != nil {
		if nRead > 0 {
			return buf[:nRead]
		}
		log.Fatalf("recv error: %v", err)
	}
	return buf[:nRead]
}

func (p *Proc) RecvLine() []byte {
	if p.closed {
		log.Fatalf("tube is closed")
	}
	var buf bytes.Buffer
	b := make([]byte, 1)
	for {
		_, err := p.stdout.Read(b)
		if err != nil {
			log.Fatalf("recv error: %v", err)
		}
		buf.Write(b)
		if bytes.HasSuffix(buf.Bytes(), p.options.NewLine) {
			result := buf.Bytes()
			return result[:len(result)-len(p.options.NewLine)]
		}
	}
}

func (p *Proc) RecvUntil(delim []byte) []byte {
	if p.closed {
		log.Fatalf("tube is closed")
	}
	var buf bytes.Buffer
	b := make([]byte, 1)
	for {
		_, err := p.stdout.Read(b)
		if err != nil {
			log.Fatalf("recv error: %v", err)
		}
		buf.Write(b)
		if bytes.HasSuffix(buf.Bytes(), delim) {
			return buf.Bytes()
		}
	}
}

func (p *Proc) RecvAll() []byte {
	if p.closed {
		log.Fatalf("tube is closed")
	}
	ret, err := io.ReadAll(p.stdout)
	if err != nil {
		log.Fatalf("recv error: %v", err)
	}
	return ret
}

func (p *Proc) Interactive() {
	if p.closed {
		log.Fatalf("tube is closed")
	}

	go io.Copy(os.Stdout, p.stdout)

	go io.Copy(os.Stderr, p.stderr)

	io.Copy(p.stdin, os.Stdin)
}

func (p *Proc) Close() {
	if p.closed {
		log.Fatalf("broken tube")
	}
	p.closed = true

	p.stdin.Close()
	p.stdout.Close()
	p.stderr.Close()

	if p.cmd.Process != nil {
		p.cmd.Process.Kill()
	}

	err := p.cmd.Wait()
	if err != nil {
		log.Fatalf("close error: %v", err)
	}
}
