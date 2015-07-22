package main

import (
	"bytes"
	"fmt"
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

func Addr(path []string) (data []byte, err error) {
	host := path2host(path)
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

func CName(path []string) (data []byte, err error) {
	host := path2host(path)
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


func Name(path []string) (data []byte, err error) {
	addr := path2host(path)
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

func MX(path []string) (data []byte, err error) {
	host := path2host(path)
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

func NS(path []string) (data []byte, err error) {
	domain := path2host(path)
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

func TXT(path []string) (data []byte, err error) {
	host := path2host(path)
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
