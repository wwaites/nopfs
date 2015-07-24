package dns

import (
	"bytes"
	"fmt"
	"hubs.net.uk/sw/nopfs"
	"net"
	"os"
)

var readme_dns = `
Name and address resolution
===========================

This directory contains name service routines. Reading the files
results in name resolution using the server's local mechanism --
usually /etc/hosts followed by DNS.

  - addr    Look up the IPv4/IPv6 address for a name
  - cname   Look up the canonical version of a name
  - name    Reverse look up of the name that corresponds to an
            address
  - mx      Look up the mail exchanger corresponding for a domain
  - ns      Look up the DNS servers for a domain
  - txt     Look up any text records for a name

`
var Readme nopfs.Dispatcher = nopfs.NewFile([]byte(readme_dns))

func addr(host string) (data []byte, err error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		err = os.ErrNotExist
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(addrs)*16))
	for _, ip := range addrs {
		buf.WriteString(ip)
		buf.WriteByte('\n')
	}
	data = buf.Bytes()
	return
}
var Addr nopfs.Dispatcher = nopfs.NewFun(nopfs.HostF(addr))

func cname(host string) (data []byte, err error) {
	cname, err := net.LookupCNAME(host)
	if err != nil {
		err = os.ErrNotExist
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(cname)+1))
	buf.WriteString(cname)
	buf.WriteByte('\n')
	data = buf.Bytes()
	return
}
var CName nopfs.Dispatcher = nopfs.NewFun(nopfs.HostF(cname))

func name(addr string) (data []byte, err error) {
	names, err := net.LookupAddr(addr)
	if err != nil {
		err = os.ErrNotExist
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(names)*32))
	for _, name := range names {
		buf.WriteString(name)
		buf.WriteByte('\n')
	}
	data = buf.Bytes()
	return
}
var Name nopfs.Dispatcher = nopfs.NewFun(nopfs.HostF(name))

func mx(host string) (data []byte, err error) {
	mxs, err := net.LookupMX(host)
	if err != nil {
		err = os.ErrNotExist
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(mxs)*20))

	for _, mx := range mxs {
		buf.WriteString(fmt.Sprintf("%d %s\n", mx.Pref, mx.Host))
	}
	data = buf.Bytes()
	return
}
var MX nopfs.Dispatcher = nopfs.NewFun(nopfs.HostF(mx))

func ns(domain string) (data []byte, err error) {
	nss, err := net.LookupNS(domain)
	if err != nil {
		err = os.ErrNotExist
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(nss)*16))

	for _, ns := range nss {
		buf.WriteString(ns.Host)
		buf.WriteByte('\n')
	}
	data = buf.Bytes()
	return
}
var NS nopfs.Dispatcher = nopfs.NewFun(nopfs.HostF(ns))

func txt(host string) (data []byte, err error) {
	txts, err := net.LookupTXT(host)
	if err != nil {
		err = os.ErrNotExist
		return
	}

	buf := bytes.NewBuffer(make([]byte, 0, len(txts)*80))

	for _, txt := range txts {
		buf.WriteString(txt)
		buf.WriteByte('\n')
	}
	data = buf.Bytes()
	return
}
var TXT nopfs.Dispatcher = nopfs.NewFun(nopfs.HostF(txt))

var Dir *nopfs.Dir
func init() {
	Dir = nopfs.NewDir()
	Dir.Append("README.txt", Readme)
	Dir.Append("addr", Addr)
	Dir.Append("cname", CName)
	Dir.Append("name", Name)
	Dir.Append("mx", MX)
	Dir.Append("ns", NS)
	Dir.Append("txt", TXT)

}