package tnfs

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"github.com/fujiNetWIFI/tnfs-gui/internal/config"
	"github.com/mitchellh/go-ps"
)

type EventType string

const (
	StatusChange EventType = "status"
	Log          EventType = "log"
	Err          EventType = "err"
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
	a := fyne.CurrentApp()
	s := &Server{
		Status:  STOPPED,
		EventCh: ch,
		cfg:     cfg,
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigch
		s.killSubprocess()
	}()

	a.Lifecycle().SetOnStopped(func() {
		s.killSubprocess()
	})

	go s.findExistingProcess()

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

func (s *Server) launchSubprocess() error {
	cmd := exec.Command(s.cfg.ExePath, s.cfg.TnfsRootPath)

	errR, _ := cmd.StderrPipe()
	outR, _ := cmd.StdoutPipe()
	s.Log = io.MultiReader(errR, outR)

	err := cmd.Start()
	if err != nil {
		return s.fail(err)
	}
	s.Process = cmd.Process

	go func() {
		err = cmd.Wait()
		if err != nil {
			fmt.Println(err.Error())
			if cmd.ProcessState.ExitCode() == 255 {
				s.fail(errors.New("TNFS port (16384) may be in use"))
			}
		}
	}()

	return nil
}

func (s *Server) killSubprocess() error {
	if s.Process == nil {
		return errors.New("Not started")
	}
	// TODO: timeout when killing
	s.Process.Signal(syscall.SIGTERM)
	_, err := s.Process.Wait()

	if err != nil && !errors.Is(err, syscall.ECHILD) {
		return s.fail(err)
	}
	return nil
}

func (s *Server) findExistingProcess() *os.Process {
	all, err := ps.Processes()
	if err != nil {
		return nil
	}
	for _, pp := range all {
		name := pp.Executable()
		if strings.HasSuffix(name, "tnfsd") || strings.HasSuffix(name, "tnfsd.exe") {
			p, _ := os.FindProcess(pp.Pid())
			if p == nil {
				continue
			}
			err := p.Signal(syscall.Signal(0))
			if err == nil {
				s.Process = p
				s.setStatus(STARTED)
				return p
			}
		}
	}
	return nil
}

func (s *Server) Start() error {
	existing := s.findExistingProcess()
	if existing != nil {
		return nil
	}

	s.setStatus(STARTING)

	err := s.launchSubprocess()
	if err != nil {
		return s.fail(err)
	}

	s.setStatus(STARTED)

	go func() {
		scanner := bufio.NewScanner(s.Log)
		for scanner.Scan() {
			s.sendLogEvent(scanner.Text())
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	if s.Status != STARTED {
		return errors.New("Not running")
	}

	s.setStatus(STOPPING)
	err := s.killSubprocess()

	if err != nil {
		fmt.Println(err.Error())
	}

	s.setStatus(STOPPED)
	return err
}
