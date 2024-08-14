package util

import (
	"os"
	"syscall"
)

func CheckProcess(p *os.Process) error {
	return p.Signal(syscall.Signal(0))
}
