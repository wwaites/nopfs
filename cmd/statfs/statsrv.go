package main

import (
	"flag"
	"hubs.net.uk/sw/statfs"
	"log"
	)

var addr = flag.String("addr", ":5640", "network address")
var debug = flag.Int("debug", 0, "print debug messages")

func main() {
	flag.Parse()

	root := statfs.NewAnyDir()
	root.Append("rtt", statfs.NewCmd(Ping))
	root.Append("rtt6", statfs.NewCmd(Ping6))
	root.Append("trace", statfs.NewCmd(Trace))
	root.Append("trace6", statfs.NewCmd(Trace6))
	root.Append("mtr", statfs.NewCmd(Mtr))
	root.Append("mtrt", statfs.NewCmd(MtrT))

	sfs := new(statfs.StatSrv)
	sfs.Debuglevel = *debug
	sfs.Root = root
	sfs.Start(sfs)
	err := sfs.StartNetListener("tcp", *addr)
	if err != nil {
		log.Fatalf("%s", err)
	}
}