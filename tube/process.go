package tube

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"pwner/utils"
	"time"
)

type Proc struct {
	baseTube
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	scanner *bufio.Scanner
}

func Process(params ...interface{}) *Proc {
	var args []string
	var options Options = DefaultOptions

	for _, param := range params {
		switch v := param.(type) {
		case string:
			args = append(args, v)
		case Options:
			options = v
		case []string:
			args = append(args, v...)
		default:
			utils.Fatal("invalid parameter type: %T", v)
		}
	}

	if len(args) == 0 {
		utils.Fatal("no command specified")
	}

	cmd := exec.Command(args[0], args[1:]...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		utils.Fatal("failed to create stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		utils.Fatal("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		utils.Fatal("failed to create stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		utils.Fatal("failed to start process: %v", err)
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

	return p
}

func (p *Proc) Send(data []byte) {
	if p.closed {
		utils.Fatal("broken tube")
		return
	}
	_, err := p.stdin.Write(data)
	if err != nil {
		utils.Fatal("failed to send data: %v", err)
	}
	return
}

func (p *Proc) SendLine(data []byte) {
	p.Send(append(data, p.options.NewLine...))
	return
}

func (p *Proc) Recv(n ...int) []byte {
	if p.closed {
		utils.Fatal("broken tube")
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
				utils.Fatal("recv error: %v", err)
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
		utils.Fatal("recv error: %v", err)
	}
	return buf[:nRead]
}

func (p *Proc) RecvLine() []byte {
	if p.closed {
		utils.Fatal("tube is closed")
	}
	var buf bytes.Buffer
	b := make([]byte, 1)
	for {
		_, err := p.stdout.Read(b)
		if err != nil {
			utils.Fatal("recv error: %v", err)
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
		utils.Fatal("tube is closed")
	}
	var buf bytes.Buffer
	b := make([]byte, 1)
	for {
		_, err := p.stdout.Read(b)
		if err != nil {
			utils.Fatal("recv error: %v", err)
		}
		buf.Write(b)
		if bytes.HasSuffix(buf.Bytes(), delim) {
			return buf.Bytes()
		}
	}
}

func (p *Proc) RecvAll() []byte {
	if p.closed {
		utils.Fatal("tube is closed")
	}
	ret, err := io.ReadAll(p.stdout)
	if err != nil {
		utils.Fatal("recv error: %v", err)
	}
	return ret
}

func (p *Proc) Interactive() {
	if p.closed {
		utils.Fatal("tube is closed")
	}

	fmt.Printf("%s[*] copying tube for interactive shell....\n%s", utils.ColorYellow, utils.ColorReset)

	go io.Copy(os.Stdout, p.stdout)
	go io.Copy(os.Stderr, p.stderr)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s[pwner]$%s ", utils.ColorRed, utils.ColorReset)
	for scanner.Scan() {
		line := scanner.Text()
		p.stdin.Write([]byte(line + "\n"))
		time.Sleep(100 * time.Millisecond)
		if p.closed {
			break
		}
		fmt.Printf("%s[pwner]$%s ", utils.ColorRed, utils.ColorReset)
	}
	os.Exit(0)
}

func (p *Proc) Close() {
	if p.closed {
		utils.Fatal("broken tube")
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
		utils.Fatal("close error: %v", err)
	}
}
