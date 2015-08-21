package main

import (
	"flag"
	"hubs.net.uk/sw/nopfs"
	"hubs.net.uk/sw/nopfs/ubnt"
	"log"
)

var addr = flag.String("addr", ":5641", "network address")
var debug = flag.Int("debug", 0, "print debug messages")

func main() {
	flag.Parse()

	sfs := new(nopfs.NopSrv)
	sfs.Debuglevel = *debug
	sfs.Root = ubnt.Dir
	sfs.Start(sfs)
	err := sfs.StartNetListener("tcp", *addr)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
