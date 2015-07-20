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
	root.Append("hello", statfs.NewHello())
	sfs := new(statfs.StatSrv)
	sfs.Debuglevel = *debug
	sfs.Root = root
	sfs.Start(sfs)
	err := sfs.StartNetListener("tcp", *addr)
	if err != nil {
		log.Fatalf("%s", err)
	}
}