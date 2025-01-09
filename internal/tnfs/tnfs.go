package server

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

	"github.com/mitchellh/go-ps"
)

type TnfsEventType string

const (
	statusChange TnfsEventType = "status"
	log          TnfsEventType = "log"
	err          TnfsEventType = "err"
)

type TnfsEvent struct {
	Type TnfsEventType
	Data string
	Time time.Time
}

type TnfsdStatus int

const (
	STOPPED TnfsdStatus = iota
	STOPPING
	STARTING
	STARTED
	FAILED
)

type TnfsServer struct {
	Status  TnfsdStatus
	Log     io.Reader
	Process *os.Process
	EventCh chan TnfsEvent
	Err     error
}

func NewTnfsServer(ch chan TnfsEvent) *TnfsServer {
	s := &TnfsServer{
		Status:  STOPPED,
		EventCh: ch,
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigch
		s.killSubprocess()
	}()

	return s
}

func (s *TnfsServer) fail(err error) error {
	defer s.setStatus(FAILED)
	s.Err = err
	return err
}

func (s *TnfsServer) setStatus(status TnfsdStatus) {
	s.Status = status
	s.sendEvent(statusChange)
}

func (s *TnfsServer) sendEvent(t TnfsEventType) {
	e := TnfsEvent{Type: t, Time: time.Now()}
	if s.EventCh != nil {
		s.EventCh <- e
	}
}

func (s *TnfsServer) sendLogEvent(msg string) {
	e := TnfsEvent{Type: log, Data: msg, Time: time.Now()}
	s.EventCh <- e
}

func (s *TnfsServer) launchSubprocess(exePath, tnfsRootPath string) error {
	cmd := exec.Command(exePath, tnfsRootPath)

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

func (s *TnfsServer) killSubprocess() error {
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

func (s *TnfsServer) findExistingProcess() *os.Process {
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

func (s *TnfsServer) Start() error {
	existing := s.findExistingProcess()
	if existing != nil {
		return nil
	}

	s.setStatus(STARTING)

	err := s.launchSubprocess(exePath, tnfsRootPath)
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

func (s *TnfsServer) Stop() error {
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
