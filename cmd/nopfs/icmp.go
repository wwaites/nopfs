package main

import (
	"os/exec"
)

var readme_icmp = `
This directory contains tests that use ICMP to measure
characteristics of the host. Chiefly this means ping(1)
and traceroute(1) as well as the more advanced mtr(1).

`

var fping string
var fping6 string
var trace string
var trace6 string
var mtr string

func init() {
	var err error
	fping, err = exec.LookPath("fping")
	if err != nil {
		panic(err)
	}

	fping6, err = exec.LookPath("fping6")
	if err != nil {
		panic(err)
	}

	trace, err = exec.LookPath("traceroute")
	if err != nil {
		panic(err)
	}

	trace6, err = exec.LookPath("traceroute6")
	if err != nil {
		panic(err)
	}

	mtr, err = exec.LookPath("mtr")
	if err != nil {
		panic(err)
	}
}

func Ping(path []string) *exec.Cmd {
	host := path[1]
	cmd := exec.Command(fping, "-e", "-r", "0", host)
	return cmd
}

func Trace(path []string) *exec.Cmd {
	host := path[1]
	return exec.Command(trace, "-I", host)
}

func Ping6(path []string) *exec.Cmd {
	host := path[1]
	return exec.Command(fping6, "-e", "-r", "0", host)
}

func Trace6(path []string) *exec.Cmd {
	host := path[1]
	return exec.Command(trace6, "-I", host)
}

func Mtr(path []string) *exec.Cmd {
	host := path[1]
	return exec.Command(mtr, "-w", "-e", "-b", "-r", host)
}

func MtrT(path []string) *exec.Cmd {
	host := path[1]
	return exec.Command(mtr, "-w", "-e", "-b", "-r", "-T", host)
}