package nopfs

import (
	"os/exec"
)

func HostF(f func(string) ([]byte, error)) func ([]string) ([]byte, error) {
	return func(path []string) ([]byte, error) {
		return f(path[1])
	}
}

func HostC(f func(string) *exec.Cmd) func ([]string) *exec.Cmd {
	return func(path []string) *exec.Cmd {
		return f(path[1])
	}
}
