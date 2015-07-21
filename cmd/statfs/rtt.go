package main

import (
	"os/exec"
)

var fping string
var trace string

func init() {
	var err error
	fping, err = exec.LookPath("fping")
	if err != nil {
		panic(err)
	}

	trace, err = exec.LookPath("traceroute")
	if err != nil {
		panic(err)
	}
}

func Ping(path []string) *exec.Cmd {
	host := path[0]
	cmd := exec.Command(fping, "-e", "-r", "0", host)
	return cmd
}

func Trace(path []string) *exec.Cmd {
	host := path[0]
	cmd := exec.Command(trace, host)
	return cmd
}