package main

import (
	"flag"
	"hubs.net.uk/sw/nopfs"
	"log"
	)

var addr = flag.String("addr", ":5640", "network address")
var debug = flag.Int("debug", 0, "print debug messages")

var readme_top = `

Network Operations File System
==============================

This directory tree exposes live information about the
network through the filesystem. It can be navigated and
manipulated using all of the usual tools for working 
with files.

  host/    information about specific hosts

`

var readme_host = `
Operations that may be done on a hostname or IP address.

  h/icmp/    ping, traceroute, etc.
  h/dns/     gathering information from domain name system.

It suffices to change into the subdirectory named for the
host or IP address. These subdirectories will not appear
in a listing but can nevertheless be descended into, for
example,

  % cat 127.0.0.1/icmp/ping
  127.0.0.1 is alive (0.02 ms)

`

func main() {
	flag.Parse()

	root := nopfs.NewDir()
	root.Append("README.txt", nopfs.NewFile([]byte(readme_top)))

	host := nopfs.NewAnyDir()
	root.Append("host", host)
	host.Static("README.txt", nopfs.NewFile([]byte(readme_host)))

	icmp := nopfs.NewDir()
	host.Append("icmp", icmp)
	icmp.Append("README.txt", nopfs.NewFile([]byte(readme_icmp)))
	icmp.Append("ping", nopfs.NewCmd(Ping))
	icmp.Append("ping6", nopfs.NewCmd(Ping6))
	icmp.Append("trace", nopfs.NewCmd(Trace))
	icmp.Append("trace6", nopfs.NewCmd(Trace6))
	icmp.Append("mtr", nopfs.NewCmd(Mtr))
	icmp.Append("mtrt", nopfs.NewCmd(MtrT))

	sfs := new(nopfs.NopSrv)
	sfs.Debuglevel = *debug
	sfs.Root = root
	sfs.Start(sfs)
	err := sfs.StartNetListener("tcp", *addr)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
