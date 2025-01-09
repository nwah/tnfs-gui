//go:build darwin || bsd || linux

package tnfs

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

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
