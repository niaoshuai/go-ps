// +build linux solaris

package ps

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

// UnixProcess is an implementation of Process that contains Unix-specific
// fields and information.
type UnixProcess struct {
	pid   int
	ppid  int
	state rune
	pgrp  int
	sid   int

	binary string
	cmd    string
	env    []string
}

func (p *UnixProcess) Pid() int {
	return p.pid
}

func (p *UnixProcess) PPid() int {
	return p.ppid
}

func (p *UnixProcess) Executable() string {
	return p.binary
}

func (p *UnixProcess) Env() []string {
	return p.env
}

func (p *UnixProcess) Cmd() string {
	return p.cmd
}

func findProcess(pid int) (Process, error, error) {
	dir := fmt.Sprintf("/proc/%d", pid)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}

		return nil, err, nil
	}

	return newUnixProcess(pid)
}

func processes() ([]Process, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	results := make([]Process, 0, 50)
	for {
		fis, err := d.Readdir(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, fi := range fis {
			// We only care about directories, since all pids are dirs
			if !fi.IsDir() {
				continue
			}

			// We only care if the name starts with a numeric
			name := fi.Name()
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			// From this point forward, any errors we just ignore, because
			// it might simply be that the process doesn't exist anymore.
			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}
			p, err, _ := newUnixProcess(int(pid))
			if err != nil {
				continue
			}
			results = append(results, p)
		}
	}

	return results, nil
}

func newUnixProcess(pid int) (*UnixProcess, error, error) {
	p := &UnixProcess{pid: pid}
	return p, p.Refresh(), p.CmdLine()
}
