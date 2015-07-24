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

This directory tree exposes live information about the network through
the filesystem. It can be navigated and manipulated using all of the
usual tools for working with files.

  host/    information about specific hosts

`

var readme_host = `
Host operations
===============

Operations that may be done on a hostname or IP address.

  h/icmp/    ping, traceroute, etc.
  h/dns/     gathering information from domain name system.

It suffices to change into the subdirectory named for the host or IP
address. These subdirectories will not appear in a listing but can
nevertheless be descended into, for example,

  % cat 127.0.0.1/icmp/ping
  127.0.0.1 is alive (0.02 ms)

This filesystem caches recently accessed hosts, and they will appear
as directories. There is a control interface for clearing the cache,
which consists of doing a write operation on the 'clear' file, as in,

  % echo > clear

`

func path2host(path []string) string {
	return path[1]
}

func main() {
	flag.Parse()

	root := nopfs.NewDir()
	root.Append("README.txt", nopfs.NewFile([]byte(readme_top)))

	host := nopfs.NewAnyDir()
	root.Append("host", host)
	host.Static("README.txt", nopfs.NewFile([]byte(readme_host)))
	host.Static("clear", &nopfs.Ctl{Writer: nopfs.AnyDirCtlReset})

	icmp := nopfs.NewDir()
	host.Append("icmp", icmp)
	icmp.Append("README.txt", nopfs.NewFile([]byte(readme_icmp)))
	if len(ping) > 0 {
		icmp.Append("ping", nopfs.NewCmd(Ping))
	}
	if len(ping6) > 0 {
		icmp.Append("ping6", nopfs.NewCmd(Ping6))
	}
	if len(trace) > 0 {
		icmp.Append("trace", nopfs.NewCmd(Trace))
	}
	if len(trace6) > 0 {
		icmp.Append("trace6", nopfs.NewCmd(Trace6))
	}
	if len(mtr) > 0 {
		icmp.Append("mtr", nopfs.NewCmd(Mtr))
		icmp.Append("mtrt", nopfs.NewCmd(MtrT))
	}

	dns := nopfs.NewDir()
	host.Append("dns", dns)
	dns.Append("README.txt", nopfs.NewFile([]byte(readme_dns)))
	dns.Append("addr", nopfs.NewFun(Addr))
	dns.Append("cname", nopfs.NewFun(CName))
	dns.Append("name", nopfs.NewFun(Name))
	dns.Append("mx", nopfs.NewFun(MX))
	dns.Append("ns", nopfs.NewFun(NS))
	dns.Append("txt", nopfs.NewFun(TXT))

	sfs := new(nopfs.NopSrv)
	sfs.Debuglevel = *debug
	sfs.Root = root
	sfs.Start(sfs)
	err := sfs.StartNetListener("tcp", *addr)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
