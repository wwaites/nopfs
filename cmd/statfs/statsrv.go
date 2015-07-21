package main

import (
	"flag"
	"hubs.net.uk/sw/statfs"
	"log"
	)

var addr = flag.String("addr", ":5640", "network address")
var debug = flag.Int("debug", 0, "print debug messages")

var readme_top = `
This is the top readme
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

	root := statfs.NewDir()
	root.Append("README.txt", statfs.NewFile([]byte(readme_top)))

	host := statfs.NewAnyDir()
	root.Append("host", host)
	host.Static("README.txt", statfs.NewFile([]byte(readme_host)))

	icmp := statfs.NewDir()
	host.Append("icmp", icmp)
	icmp.Append("README.txt", statfs.NewFile([]byte(readme_icmp)))
	icmp.Append("ping", statfs.NewCmd(Ping))
	icmp.Append("ping6", statfs.NewCmd(Ping6))
	icmp.Append("trace", statfs.NewCmd(Trace))
	icmp.Append("trace6", statfs.NewCmd(Trace6))
	icmp.Append("mtr", statfs.NewCmd(Mtr))
	icmp.Append("mtrt", statfs.NewCmd(MtrT))

	sfs := new(statfs.StatSrv)
	sfs.Debuglevel = *debug
	sfs.Root = root
	sfs.Start(sfs)
	err := sfs.StartNetListener("tcp", *addr)
	if err != nil {
		log.Fatalf("%s", err)
	}
}