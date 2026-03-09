package shell

import (
	"adb-backup/internal/log"
	"fmt"
	"io"
	"strings"
	"sync/atomic"
	"time"

	adb "github.com/tablebird/goadb"
	"github.com/tablebird/goadb/wire"
)

type Shell interface {
	RunCommand(cmd string, args ...string) (string, error)
}

type ReaderCloserShell interface {
	Shell
	RunCommandReaderCloser(cmd string, args ...string) (io.ReadCloser, error)
}

type AsyncShell interface {
	Shell
	AsyncRunCommand(cmd string, args ...string) error
}

type AnyShell interface {
	Shell
	ReaderCloserShell
	AsyncShell
}

func NewShell(a *adb.Adb, serial string) AnyShell {
	return &realShell{a: a, serial: serial,
		status: 0, cmd: make(chan string, 100), out: make(chan string, 100)}
}

type realShell struct {
	a      *adb.Adb
	serial string
	cmd    chan string
	out    chan string
	status int32
}

func (s *realShell) RunCommand(cmd string, args ...string) (string, error) {
	cmd, err := prepareCommandLine(cmd, args...)
	if err != nil {
		return "", err
	}
	conn, err := s.dialDevice()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	err = s.sendShell(conn, cmd)
	if err != nil {
		return "", err
	}
	resp, err := conn.ReadUntilEof()
	return string(resp), err
}

func (s *realShell) RunCommandReaderCloser(cmd string, args ...string) (io.ReadCloser, error) {
	cmd, err := prepareCommandLine(cmd, args...)
	if err != nil {
		return nil, err
	}
	conn, err := s.dialDevice()
	if err != nil {
		return nil, err
	}
	err = s.sendShell(conn, cmd)
	if err != nil {
		return nil, err
	}
	return conn, err
}

func (s *realShell) AsyncRunCommand(cmd string, args ...string) error {
	cmd, err := prepareCommandLine(cmd, args...)
	if err != nil {
		return err
	}
	go s.startShell()

	s.cmd <- cmd
	return nil
}
func (s *realShell) sendShell(conn *wire.Conn, cmd string) error {
	req := fmt.Sprintf("shell:%s", cmd)
	// Shell responses are special, they don't include a length header.
	// We read until the stream is closed.
	// So, we can't use conn.RoundTripSingleResponse.
	if err := conn.SendMessage([]byte(req)); err != nil {
		return err
	}
	if _, err := conn.ReadStatus(req); err != nil {
		return err
	}
	return nil
}

func (s *realShell) startShell() error {
	if !atomic.CompareAndSwapInt32(&s.status, 0, 1) {
		return nil
	}
	defer atomic.StoreInt32(&s.status, 0)
	conn, err := s.dialDevice()
	if err != nil {
		log.ErrorF("dial error: %v", err)
		return err
	}
	defer conn.Close()
	s.sendShell(conn, "")
	go s.readOut(conn)
	for cmd := range s.cmd {
		// Shell responses are special, they don't include a length header.
		// We read until the stream is closed.
		// So, we can't use conn.RoundTripSingleResponse.
		if !strings.HasSuffix(cmd, "\n") {
			cmd += "\n"
		}
		if _, err = conn.Write([]byte(cmd)); err != nil {
			log.ErrorF("send message error: %v", err)
			return err
		}
	}
	return nil
}

func (s *realShell) readOut(conn *wire.Conn) error {
	buf := make([]byte, 4096)
	lastLine := ""
	for {
		n, err := conn.Read(buf)
		if err != nil {
			close(s.cmd)
			close(s.out)
			return err
		}

		res := string(buf[:n])
		lines := strings.Split(res, "\n")
		lines[0] = lastLine + lines[0]
		lastLine = ""
		size := len(lines)
		for i, line := range lines {
			if i < size-1 {
				if len(s.out) > 99 {
					<-s.out
				}
				s.out <- line
			} else {
				if line != "" {
					lastLine = line
				}
			}
		}
	}
}

func (s *realShell) dialDevice() (*wire.Conn, error) {
	conn, err := s.a.Dial()
	if err != nil {
		return nil, err
	}
	req := fmt.Sprintf("host:transport:%s", s.serial)
	if err = wire.SendMessageString(conn, req); err != nil {
		conn.Close()
		return nil, err
	}

	if _, err = conn.ReadStatus(req); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

// 执行函数并打印耗时
func executeWithTiming(name string, fn func()) {
	start := time.Now()
	fn()
	elapsed := time.Since(start)
	log.DebugF("%s execution time: %d ms", name, elapsed.Milliseconds())
}
