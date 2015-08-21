package nopfs

import (
	"errors"
	"github.com/rminnich/go9p"
	"log"
	"syscall"
)

type NopSrv struct {
	go9p.Srv
	DebugLevel int
	Root       Dispatcher
}

func (sfs *NopSrv) Attach(req *go9p.SrvReq) {
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

func (sfs *NopSrv) Stat(req *go9p.SrvReq) {
	fid := req.Fid.Aux.(Dispatcher)
	if sfs.Debuglevel > 0 {
		log.Printf("stat %s", fid)
	}
	req.RespondRstat(Fstat(fid))
}

func (sfs *NopSrv) Walk(req *go9p.SrvReq) {
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
		nfid, err := fid.Walk(req, tc.Wname[0])
		if err != nil {
			req.RespondError(toError(err))
			return
		}
		req.Newfid.Aux = nfid
		w := []go9p.Qid{*Qid(nfid)}
		req.RespondRwalk(w)
	}
}

func (sfs *NopSrv) Open(req *go9p.SrvReq) {
	fid := req.Fid.Aux.(Dispatcher)
	if sfs.Debuglevel > 0 {
		log.Printf("open %s", fid)
	}
	req.RespondRopen(Qid(fid), 0)
}

func (sfs *NopSrv) Read(req *go9p.SrvReq) {
	fid := req.Fid.Aux.(Dispatcher)
	tc := req.Tc
	rc := req.Rc

	if sfs.Debuglevel > 0 {
		log.Printf("read %T %s %d:%d", fid, fid, tc.Offset, tc.Count)
	}

	buf, err := fid.Read(req)
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

func (s *NopSrv) ConnOpened(conn *go9p.Conn) {
	if conn.Srv.Debuglevel > 0 {
		log.Println("connected")
	}
	s.Debuglevel = conn.Srv.Debuglevel
}

func (*NopSrv) ConnClosed(conn *go9p.Conn) {
	if conn.Srv.Debuglevel > 0 {
		log.Println("disconnected")
	}
}

func (sfs *NopSrv) FidDestroy(sfid *go9p.SrvFid) {
	if sfid.Aux == nil {
		return
	}
	fid := sfid.Aux.(Dispatcher)
	if sfs.Debuglevel > 0 {
		log.Printf("destroy %s", fid)
	}
	fid.Close()
}

func (sfs *NopSrv) Flush(req *go9p.SrvReq) {
	fid := req.Fid.Aux.(Dispatcher)
	if sfs.Debuglevel > 0 {
		log.Printf("flush %s", fid)
	}
	fid.Flush(req)
}

func (*NopSrv) Create(req *go9p.SrvReq) {
	log.Printf("create: %p", req)
	req.RespondError(errors.New("create: ..."))
}

func (sfs *NopSrv) Write(req *go9p.SrvReq) {
	fid := req.Fid.Aux.(Dispatcher)
	tc := req.Tc
	if sfs.Debuglevel > 0 {
		log.Printf("write: %f", fid)
	}

	e := fid.Write(req, tc.Data)
	if e != nil {
		req.RespondError(toError(e))
		return
	}

	req.RespondRwrite(uint32(len(tc.Data)))
}

func (*NopSrv) Clunk(req *go9p.SrvReq) {
	req.RespondRclunk()
}

func (*NopSrv) Remove(req *go9p.SrvReq) {
	log.Printf("remove: %p", req)
	req.RespondError(errors.New("remove: ..."))
}

func (*NopSrv) Wstat(req *go9p.SrvReq) {
	//	log.Printf("wstat: %p", req)
	//	req.RespondError(errors.New("wstat: ..."))
	req.RespondRwstat()
}
