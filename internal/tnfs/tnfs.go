package tnfs

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/fujiNetWIFI/tnfs-gui/internal/config"
	"github.com/nwah/gotnfsd"
)

type EventType string

const (
	StatusChange EventType = "status"
	Log          EventType = "log"
	Err          EventType = "err"

	DEFAULT_PORT = 16384
)

type Event struct {
	Type EventType
	Data string
	Time time.Time
}

type Status int

const (
	STOPPED Status = iota
	STOPPING
	STARTING
	STARTED
	FAILED
)

type Server struct {
	Status  Status
	Log     io.Reader
	Process *os.Process
	EventCh chan Event
	Err     error

	cfg *config.Config
}

func NewServer(cfg *config.Config, ch chan Event) *Server {
	s := &Server{
		Status:  STOPPED,
		EventCh: ch,
		cfg:     cfg,
	}

	r, w, _ := os.Pipe()

	go func() {
		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanBytes)
		for scanner.Scan() {
			s.sendLogEvent(scanner.Text())
		}
	}()

	gotnfsd.Init(w)

	return s
}

func (s *Server) fail(err error) error {
	defer s.setStatus(FAILED)
	s.Err = err
	return err
}

func (s *Server) setStatus(status Status) {
	s.Status = status
	s.sendEvent(StatusChange)
}

func (s *Server) sendEvent(t EventType) {
	e := Event{Type: t, Time: time.Now()}
	if s.EventCh != nil {
		s.EventCh <- e
	}
}

func (s *Server) sendLogEvent(msg string) {
	e := Event{Type: Log, Data: msg, Time: time.Now()}
	s.EventCh <- e
}

func (s *Server) captureStderr() {
	r, w, _ := os.Pipe()
	os.Stderr = w

	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			s.sendLogEvent(scanner.Text())
		}
	}()
}

func (s *Server) Start() {
	s.setStatus(STARTED)

	// s.captureStderr()

	go func() {
		err := gotnfsd.Start(s.cfg.TnfsRootPath, DEFAULT_PORT, false)
		if err != nil {
			s.fail(err)
		}
	}()
}

func (s *Server) Stop() {
	if s.Status != STARTED {
		return
	}
	gotnfsd.Stop()
	s.setStatus(STOPPED)
}
