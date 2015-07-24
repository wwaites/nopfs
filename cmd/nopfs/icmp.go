package main

import (
	"log"
	"os/exec"
)

var readme_icmp = `
ICMP and similar network probes
===============================

This directory contains tests that use ICMP to measure characteristics
of the host. Chiefly this means ping(1) and traceroute(1) as well as
the more advanced mtr(1).

`

var ping string
var ping6 string
var trace string
var trace6 string
var mtr string
var fping = true
var fping6 = true

func init() {
	var err error
	ping, err = exec.LookPath("fping")
	if err != nil {
		ping, _ = exec.LookPath("ping")
		fping = false
	}
	log.Printf("[icmp] echo program: %s", ping)
	ping6, err = exec.LookPath("fping6")
	if err != nil {
		ping6, _ = exec.LookPath("ping6")
		fping6 = false
	}
	log.Printf("[icmp6] echo program: %s", ping6)
	trace, _ = exec.LookPath("traceroute")
	log.Printf("[icmp] traceroute program: %s", trace)
	trace6, _ = exec.LookPath("traceroute6")
	log.Printf("[icmp6] traceroute program: %s", trace6)
	mtr, _ = exec.LookPath("mtr")
	log.Printf("[icmp] pretty traceroute: %s", mtr)
}

func Ping(path []string) *exec.Cmd {
	host := path2host(path)
	var cmd *exec.Cmd
	if fping {
		cmd = exec.Command(ping, "-e", "-r", "0", host)
	} else {
		cmd = exec.Command(ping, "-c", "1", host)
	}
	return cmd
}

func Trace(path []string) *exec.Cmd {
	host := path2host(path)
	return exec.Command(trace, "-I", host)
}

func Ping6(path []string) *exec.Cmd {
	host := path2host(path)
	var cmd *exec.Cmd
	if fping6 {
		cmd = exec.Command(ping6, "-e", "-r", "0", host)
	} else {
		cmd = exec.Command(ping6, "-c", "1", host)
	}
	return cmd
}

func Trace6(path []string) *exec.Cmd {
	host := path2host(path)
	return exec.Command(trace6, "-I", host)
}

func Mtr(path []string) *exec.Cmd {
	host := path2host(path)
	return exec.Command(mtr, "-w", "-e", "-b", "-r", host)
}

func MtrT(path []string) *exec.Cmd {
	host := path2host(path)
	return exec.Command(mtr, "-w", "-e", "-b", "-r", "-T", host)
}
