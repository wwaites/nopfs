# Network Operations File System

This package implements a file system server that 
exposes information about the network. It allows,
for example, reading a certain file to obtain real
time statistics.

This filesystem is made available via the 9p/Styx
protocol and can be mounted over the network from
all modern operating systems. Under Linux, for example,
after running the daemon,

    % nopfs &
    
it may be installed into the filesystem with either
of the short or the long form commands,

    % 9mount -u tcp!localhost!5640 /mnt

    - or -

    % mount -t 9p -o tcp,trans=tcp,nodev,port=5640 127.0.0.1 /mnt

The `Hello World' of this arrangement is,

    % cat /mnt/host/news.bbc.co.uk/icmp/ping
    news.bbc.co.uk is alive (84.1 ms)

## Prerequisites

The Go language compiler version 1.4 or later is
required to build this package. Furthermore the 
following executables are runtime dependencies and
must be present in the search path,

* fping
* traceroute
* mtr

## Installation

Should be as simple as

    go get hubs.net.uk/sw/nopfs

which will result in the ${GOPATH}/bin/nopfs program being
built and installed
