//go:build windows

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
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

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

	if err := s.Process.Kill(); err != nil {
		return s.fail(err)
	}

	return nil
}
