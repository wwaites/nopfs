package statfs

import (
	"strings"
	"crypto/sha256"
	"encoding/binary"
	"github.com/rminnich/go9p"
	"time"
	"os"
	"os/exec"
	"sync"
)

type Dispatcher interface {
	Clone()      Dispatcher
	GetPath()    []string
	SetPath([]string)
	IsDir()      bool
	Read()       ([]byte, error)
	Inode()      uint64
	Flush()
	Walk(string) (Dispatcher, error)
	Close()
}

func Qid(d Dispatcher) (q *go9p.Qid) {
	q = new(go9p.Qid)
	if d.IsDir() {
		q.Type = go9p.QTDIR
	} else {
		q.Type = go9p.QTFILE
	}
	q.Path = d.Inode()
	return
}

func Fstat(d Dispatcher) (p *go9p.Dir) {
	p = new(go9p.Dir)

	p.Qid = *Qid(d)
	if d.IsDir() {
		p.Mode |= go9p.DMDIR
		p.Mode |= 0111
	}
	p.Mode |= 0444

	p.Uid = "none"
	p.Uidnum = uint32(0)
	p.Gid = "none"
	p.Gidnum = uint32(0)
	p.Muid = "none"
	p.Muidnum = go9p.NOUID
	p.Ext = ""

	now := time.Now().Unix()
	p.Atime = uint32(now)
	p.Mtime = uint32(now)

	path := d.GetPath()
	if len(path) > 0 {
		p.Name = path[len(path)-1]
	}
	return
}

type Path struct {
	path []string
}

func (p *Path) GetPath() []string {
	return p.path
}

func (p *Path) SetPath(path []string) {
	p.path = path
}

func (p *Path) Inode() uint64 {
	buf  := strings.Join(p.path, "/")
	hash := sha256.Sum256([]byte(buf))
	i, _ := binary.Uvarint(hash[:])
	return i
}

func (p *Path) String() string {
	return strings.Join(p.path, "/")
}

type Dir struct {
	sync.Mutex
	Path
	entries map[string]Dispatcher
	listing []byte
}

func NewDir() (d *Dir) {
	d = &Dir{}
	d.SetPath(make([]string, 0))
	d.entries = make(map[string]Dispatcher)
	return
}


func (d *Dir) IsDir() bool {
	return true
}

func (d *Dir) Clone() Dispatcher {
	d.Lock()
	defer d.Unlock()
	n := new(Dir)
	n.SetPath(d.GetPath())
	n.entries = d.entries
	n.listing = d.listing
	return n
}

func (d *Dir) Walk(name string) (Dispatcher, error) {
	subDisp, ok := d.entries[name]
	if !ok {
		return nil, os.ErrNotExist
	} else {
		newDisp := subDisp.Clone()
		path := append(d.path, name)
		newDisp.SetPath(path)
		return newDisp, nil
	}
}

func (d *Dir) Read() ([]byte, error) {
	d.Lock()
	defer d.Unlock()
	if d.listing == nil {
		d.listing = make([]byte, 0)
		for name, subDisp := range d.entries {
			path := append(d.GetPath(), name)
			newDisp := subDisp.Clone()
			newDisp.SetPath(path)
			b := go9p.PackDir(Fstat(newDisp), true)
			d.listing = append(d.listing, b...)
		}
	}
	return d.listing, nil
}

func (d *Dir) Append(name string, disp Dispatcher) {
	d.entries[name] = disp
}

func (d *Dir) Close() {
	d.Lock()
	defer d.Unlock()
	d.listing = nil
}

func (d *Dir) Flush() {}

type AnyDir struct {
	Path
	entries map[string]Dispatcher
}

func NewAnyDir() (a *AnyDir) {
	a = &AnyDir{}
	a.entries = make(map[string]Dispatcher)
	return
}

func (a *AnyDir) IsDir() bool {
	return true
}

func (a *AnyDir) Walk(name string) (Dispatcher, error) {
	subDisp := NewDir()
	path := append(a.path, name)
	subDisp.SetPath(path)
	subDisp.entries = a.entries
	return subDisp, nil
}

func (a *AnyDir) Read() ([]byte, error) {
	return make([]byte, 0), nil
}

func (a *AnyDir) Append(name string, disp Dispatcher) {
	a.entries[name] = disp
}

func (a *AnyDir) Clone() Dispatcher {
	n := NewAnyDir()
	n.SetPath(a.GetPath())
	return n
}

func (a *AnyDir) Close() {}
func (a *AnyDir) Flush() {}

type PseudoFile struct {
	Path
}

func (f *PseudoFile) IsDir() bool {
	return false
}

func (f *PseudoFile) Walk(name string) (d Dispatcher, err error) {
	err = os.ErrInvalid
	return
}

type Cmd struct {
	PseudoFile

	cfun  func ([]string)*exec.Cmd
	clock sync.Mutex
	cmd   *exec.Cmd

	dlock sync.Mutex
	data  []byte

	err   error
}

func NewCmd(cmd func ([]string)*exec.Cmd) (c *Cmd) {
	c = &Cmd{}
	c.cfun = cmd
	c.SetPath(make([]string, 0))
	return
}

func (c *Cmd) Clone() Dispatcher {
	n := NewCmd(c.cfun)
	n.SetPath(c.GetPath())
	return n
}

func (c *Cmd) Close() {
	c.dlock.Lock()
	defer c.dlock.Unlock()
	c.data, c.err = nil, nil
}

func (c *Cmd) Read() ([]byte, error) {
	c.dlock.Lock()
	defer c.dlock.Unlock()
	if c.data == nil {
		c.clock.Lock()
		c.cmd = c.cfun(c.path)
		c.clock.Unlock()

		c.data, c.err = c.cmd.CombinedOutput()

		c.clock.Lock()
		c.cmd = nil
		c.clock.Unlock()
	}
	return c.data, c.err
}

func (c *Cmd) Flush() {
	c.clock.Lock()
	if c.cmd != nil {
		c.cmd.Process.Kill()
	}
	c.clock.Unlock()
}