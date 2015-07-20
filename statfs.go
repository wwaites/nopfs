// Copyright 2009 The go9p Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statfs

import (
	"github.com/rminnich/go9p"
	"errors"
	"log"
	"syscall"
)

type StatSrv struct {
	go9p.Srv
	DebugLevel int
	Root       Dispatcher
}

func (sfs *StatSrv) Attach(req *go9p.SrvReq) {
	if req.Afid != nil {
		req.RespondError(go9p.Enoauth)
		return
	}

	req.Fid.Aux = sfs.Root
	if sfs.Debuglevel > 0 {
		log.Printf("attach")
	}

	req.RespondRattach(Qid(sfs.Root))
}

func (sfs *StatSrv) Stat(req *go9p.SrvReq) {
	fid := req.Fid.Aux.(Dispatcher)
	if sfs.Debuglevel > 0 {
		log.Printf("stat %s", fid)
	}
	req.RespondRstat(Fstat(fid))
}


func (sfs *StatSrv) Walk(req *go9p.SrvReq) {
	fid := req.Fid.Aux.(Dispatcher)
	tc := req.Tc

	if sfs.Debuglevel > 0 {
		log.Printf("walk %s %s", fid, tc.Wname)
	}

	if req.Newfid.Aux == nil {
		req.Newfid.Aux = fid.Clone()
	}

	if len(tc.Wname) == 0 {
		w := make([]go9p.Qid, 0)
		req.RespondRwalk(w)
	} else {
		nfid, err := fid.Walk(tc.Wname[0])
		if err != nil {
			req.RespondError(toError(err))
			return
		}
		req.Newfid.Aux = nfid
		w := []go9p.Qid{*Qid(nfid)}
		req.RespondRwalk(w)
	}
}

func (sfs *StatSrv) Open(req *go9p.SrvReq) {
	fid := req.Fid.Aux.(Dispatcher)
	if sfs.Debuglevel > 0 {
		log.Printf("open %s", fid)
	}
	req.RespondRopen(Qid(fid), 0)
}

func (sfs *StatSrv) Read(req *go9p.SrvReq) {
	fid := req.Fid.Aux.(Dispatcher)
	tc := req.Tc
	rc := req.Rc
	
	if sfs.Debuglevel > 0 {
		log.Printf("read %T %s %d:%d", fid, fid, tc.Offset, tc.Count)
	}

	buf, err := fid.Read()
	if err != nil {
		req.RespondError(toError(err))
		return
	}
	
	go9p.InitRread(rc, tc.Count)
	count := 0
	switch {
	case tc.Offset > uint64(len(buf)):
		count = 0
	case len(buf[tc.Offset:]) > int(tc.Count):
		count = int(tc.Count)
	default:
		count = len(buf[tc.Offset:])
	}

	copy(rc.Data, buf[tc.Offset:int(tc.Offset)+count])
	go9p.SetRreadCount(rc, uint32(count))
	req.Respond()
}

func toError(err error) *go9p.Error {
	var ecode uint32

	ename := err.Error()
	if e, ok := err.(syscall.Errno); ok {
		ecode = uint32(e)
	} else {
		ecode = go9p.EIO
	}

	return &go9p.Error{ename, ecode}
}

func (s *StatSrv) ConnOpened(conn *go9p.Conn) {
	if conn.Srv.Debuglevel > 0 {
		log.Println("connected")
	}
	s.Debuglevel = conn.Srv.Debuglevel
}

func (*StatSrv) ConnClosed(conn *go9p.Conn) {
	if conn.Srv.Debuglevel > 0 {
		log.Println("disconnected")
	}
}

func (sfs *StatSrv) FidDestroy(sfid *go9p.SrvFid) {
	if sfid.Aux == nil {
		return
	}
	fid := sfid.Aux.(Dispatcher)
	if sfs.Debuglevel > 0 {
		log.Printf("destroy %s", fid)
	}
	fid.Close()
}

func (*StatSrv) Flush(req *go9p.SrvReq) {
	log.Printf("flush: %p", req)
}

func (*StatSrv) Create(req *go9p.SrvReq) {
	log.Printf("create: %p", req)
	req.RespondError(errors.New("create: ..."))
}


func (*StatSrv) Write(req *go9p.SrvReq) {
	log.Printf("write: %p", req)
	req.RespondError(errors.New("write: ..."))
}

func (*StatSrv) Clunk(req *go9p.SrvReq) {
	req.RespondRclunk() 
}

func (*StatSrv) Remove(req *go9p.SrvReq) {
	log.Printf("remove: %p", req)
	req.RespondError(errors.New("remove: ..."))
}

func (*StatSrv) Wstat(req *go9p.SrvReq) {
	log.Printf("wstat: %p", req)
	req.RespondError(errors.New("wstat: ..."))
}
