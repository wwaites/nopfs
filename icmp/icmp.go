package icmp

import (
	"os/exec"
	"hubs.net.uk/sw/nopfs"
)

var readme_icmp = `
ICMP and similar network probes
===============================

This directory contains tests that use ICMP to measure characteristics
of the host. Chiefly this means ping(1) and traceroute(1) as well as
the more advanced mtr(1).

`
var Readme nopfs.Dispatcher = nopfs.NewFile([]byte(readme_icmp))

var ping_prog string
var ping6_prog string
var trace_prog string
var trace6_prog string
var mtr_prog string

var fping = true
var fping6 = true

var Ping nopfs.Dispatcher
var Ping6 nopfs.Dispatcher
var Trace nopfs.Dispatcher
var Trace6 nopfs.Dispatcher
var Mtr nopfs.Dispatcher
var MtrT nopfs.Dispatcher
var Dir *nopfs.Dir

func init() {
	var err error
	ping_prog, err = exec.LookPath("fping")
	if err != nil {
		ping_prog, err = exec.LookPath("ping")
		fping = false
	}
	if err == nil {
		Ping = nopfs.NewCmd(nopfs.HostC(ping))
	}

	ping6_prog, err = exec.LookPath("fping6")
	if err != nil {
		ping6_prog, err = exec.LookPath("ping6")
		fping6 = false
	}
	if err == nil {
		Ping6 = nopfs.NewCmd(nopfs.HostC(ping6))
	}

	trace_prog, err = exec.LookPath("traceroute")
	if err == nil {
		Trace = nopfs.NewCmd(nopfs.HostC(trace))
	}

	trace6_prog, err = exec.LookPath("traceroute6")
	if err == nil {
		Trace6 = nopfs.NewCmd(nopfs.HostC(trace6))
	}

	mtr_prog, err = exec.LookPath("mtr")
	if err == nil {
		Mtr = nopfs.NewCmd(nopfs.HostC(mtr))
		MtrT = nopfs.NewCmd(nopfs.HostC(mtrt))
	}

	Dir = nopfs.NewDir()
	Dir.Append("README.txt", Readme)

	if Ping != nil {
		Dir.Append("ping", Ping)
	}
	if Ping6 != nil {
		Dir.Append("ping6", Ping6)
	}
	if Trace != nil {
		Dir.Append("trace", Trace)
	}
	if Trace6 != nil {
		Dir.Append("trace6", Trace6)
	}
	if Mtr != nil {
		Dir.Append("mtr", Mtr)
		Dir.Append("mtrt", MtrT)
	}
}

func ping(host string) *exec.Cmd {
	var cmd *exec.Cmd
	if fping {
		cmd = exec.Command(ping_prog, "-e", "-r", "0", host)
	} else {
		cmd = exec.Command(ping_prog, "-c", "1", host)
	}
	return cmd
}

func ping6(host string) *exec.Cmd {
	var cmd *exec.Cmd
	if fping6 {
		cmd = exec.Command(ping6_prog, "-e", "-r", "0", host)
	} else {
		cmd = exec.Command(ping6_prog, "-c", "1", host)
	}
	return cmd
}

func trace6(host string) *exec.Cmd {
	return exec.Command(trace6_prog, "-I", host)
}

func trace(host string) *exec.Cmd {
	return exec.Command(trace_prog, "-I", host)
}

func mtr(host string) *exec.Cmd {
	return exec.Command(mtr_prog, "-w", "-e", "-b", "-r", host)
}

func mtrt(host string) *exec.Cmd {
	return exec.Command(mtr_prog, "-w", "-e", "-b", "-r", "-T", host)
}
