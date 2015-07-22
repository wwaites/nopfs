package nopfs

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
	Clone()               Dispatcher
	GetPath()             []string
	SetPath([]string)
	GetParent()           Dispatcher
	SetParent(Dispatcher)

	IsDir()           bool
	Size()            uint64
	Read()            ([]byte, error)
	Write([]byte)     error
	Inode()           uint64
	Perms()           uint32
	Flush()
	Walk(string)      (Dispatcher, error)
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
	}
	p.Mode |= d.Perms()

	p.Uid = "none"
	p.Uidnum = uint32(0)
	p.Gid = "none"
	p.Gidnum = uint32(0)
	p.Muid = "none"
	p.Muidnum = go9p.NOUID
	p.Ext = ""

	p.Length = d.Size()

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
	parent Dispatcher
	path   []string
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

func (p *Path) GetParent() Dispatcher {
	return p.parent
}

func (p *Path) SetParent(parent Dispatcher) {
	p.parent = parent
}

type Dir struct {
	sync.RWMutex
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

func (d *Dir) Perms() uint32 {
	return 0555
}

func (d *Dir) IsDir() bool {
	return true
}

func (d *Dir) Size() uint64 {
	return uint64(0)
}

func (d *Dir) Clone() Dispatcher {
	d.RLock()
	defer d.RUnlock()
	n := new(Dir)
	n.SetPath(d.GetPath())
	n.SetParent(d.GetParent())
	n.entries = d.entries
	n.listing = d.listing
	return n
}

func (d *Dir) Walk(name string) (Dispatcher, error) {
	d.RLock()
	defer d.RUnlock()
	subDisp, ok := d.entries[name]
	if !ok {
		return nil, os.ErrNotExist
	} else {
		newDisp := subDisp.Clone()
		path := append(d.path, name)
		newDisp.SetPath(path)
		newDisp.SetParent(d)
		return newDisp, nil
	}
}

func (d *Dir) Read() ([]byte, error) {
	d.RLock()
	defer d.RUnlock()
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

func (d *Dir) Append(name string, disp Dispatcher) *Dir {
	d.Lock()
	defer d.Unlock()
	d.entries[name] = disp
	d.listing = nil
	return d
}

func (d *Dir) Write([]byte) error {
	return os.ErrInvalid
}
func (d *Dir) Close() {}
func (d *Dir) Flush() {}

type AnyDir struct {
	Path
	lock    *sync.RWMutex
	entries map[string]Dispatcher
	history map[string]bool
	static  map[string]Dispatcher
	listing []byte
}

func NewAnyDir() (a *AnyDir) {
	a = &AnyDir{}
	a.lock    = &sync.RWMutex{}
	a.entries = make(map[string]Dispatcher)
	a.static  = make(map[string]Dispatcher)
	a.history = make(map[string]bool)
	return
}

func (a *AnyDir) IsDir() bool {
	return true
}

func (a *AnyDir) Perms() uint32 {
	return 0555
}

func (a *AnyDir) Size() uint64 {
	return uint64(0)
}

func (a *AnyDir) Walk(name string) (Dispatcher, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	subDisp, ok := a.static[name]
	if ok {
		newDisp := subDisp.Clone()
		path := append(a.path, name)
		newDisp.SetPath(path)
		newDisp.SetParent(a)
		return newDisp, nil
	} else {
		a.history[name] = true
		subDir := NewDir()
		path := append(a.path, name)
		subDir.SetPath(path)
		subDir.SetParent(a)
		subDir.entries = a.entries
		return subDir, nil
	}
}

func (a *AnyDir) Read() ([]byte, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()
	if a.listing == nil {
		a.listing = make([]byte, 0)
		for name, subDisp := range a.static {
			path := append(a.GetPath(), name)
			newDisp := subDisp.Clone()
			newDisp.SetPath(path)
			b := go9p.PackDir(Fstat(newDisp), true)
			a.listing = append(a.listing, b...)
		}
		for name, _ := range a.history {
			subDir := NewDir()
			path := append(a.path, name)
			subDir.SetPath(path)
			subDir.entries = a.entries
			b := go9p.PackDir(Fstat(subDir), true)
			a.listing = append(a.listing, b...)
		}
	}
	return a.listing, nil
}


func (a *AnyDir) Clone() Dispatcher {
	a.lock.RLock()
	defer a.lock.RUnlock()
	n := NewAnyDir()
	n.SetPath(a.GetPath())
	n.SetParent(a.GetParent())
	n.lock    = a.lock
	n.entries = a.entries
	n.static  = a.static
	n.listing = a.listing
	n.history = a.history
	return n
}

func (a *AnyDir) Append(name string, disp Dispatcher) *AnyDir {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.entries[name] = disp
	return a
}

func (a *AnyDir) Static(name string, disp Dispatcher) *AnyDir {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.static[name] = disp
	a.listing = nil
	return a
}

func (a *AnyDir) Write([]byte) error {
	return os.ErrInvalid
}
func (a *AnyDir) Close() {}
func (a *AnyDir) Flush() {}

func AnyDirCtlReset(c *Ctl, data []byte) (resp []byte, err error) {
	dir, ok := c.GetParent().(*AnyDir)
	if !ok {
		err = os.ErrInvalid
		return
	}
	dir.lock.Lock()
	defer dir.lock.Unlock()
	for k, _ := range dir.history {
		delete(dir.history, k)
	}
	resp = []byte("ok")
	return
}

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

func (f *PseudoFile) Write([]byte) error {
	return os.ErrInvalid
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

func (c *Cmd) Perms() uint32 {
	return 0444
}


func (c *Cmd) Clone() Dispatcher {
	n := NewCmd(c.cfun)
	n.SetPath(c.GetPath())
	n.SetParent(c.GetParent())
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
		c.cmd = c.cfun(c.GetPath())
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

func (c *Cmd) Size() uint64 {
	return uint64(0)
}


type Fun struct {
	PseudoFile
	sync.Mutex
	fun  func([]string) ([]byte, error)
	data []byte
}

func NewFun(fun func([]string) ([]byte, error)) *Fun {
	f := &Fun{}
	f.fun = fun
	f.SetPath(make([]string, 0))
	return f
}

func (f *Fun) Perms() uint32 {
	return 0444
}

func (f *Fun) Clone() Dispatcher {
	n := NewFun(f.fun)
	n.SetPath(f.GetPath())
	n.SetParent(f.GetParent())
	return n
}

func (f *Fun) Read() (data []byte, err error) {
	f.Lock()
	defer f.Unlock()
	if f.data == nil {
		data, err = f.fun(f.GetPath())
		if err != nil {
			f.data = data
		}
	} else {
		data = f.data
	}
	return
}

func (f *Fun) Size() uint64 { return uint64(0) }
func (f *Fun) Flush() {}
func (f *Fun) Close() {}

type File struct {
	PseudoFile
	data []byte
}

func NewFile(data []byte) *File {
	f := &File{}
	f.data = data
	return f
}

func (f *File) Perms() uint32 {
	return 0444
}

func (f *File) Clone() Dispatcher {
	n := NewFile(f.data)
	n.SetPath(f.GetPath())
	n.SetParent(f.GetParent())
	return n
}

func (f *File) Read() ([]byte, error) {
	return f.data, nil
}

func (f *File) Size() uint64 {
	return uint64(len(f.data))
}

func (f *File) Close() {}
func (f *File) Flush() {}

type Ctl struct {
	Path
	sync.RWMutex
	Writer   func(*Ctl, []byte) ([]byte, error)
	buf      []byte
}

func (c *Ctl) Clone() Dispatcher {
	n := &Ctl{Writer: c.Writer}
	n.SetPath(c.GetPath())
	n.SetParent(c.GetParent())
	return n
}

func (c *Ctl) Perms() uint32 { return 0666 }
func (c *Ctl) IsDir() bool { return false }
func (c *Ctl) Close() {}
func (c *Ctl) Flush() {}

func (c *Ctl) Read() (data []byte, err error) {
	c.RLock()
	defer c.RUnlock()
	if c.buf == nil {
		err = os.ErrNotExist
		return
	}
	data = c.buf
	return
}

func (c *Ctl) Write(data []byte) (err error) {
	c.Lock()
	defer c.Unlock()
	c.buf, err = c.Writer(c, data)
	return
}

func (c *Ctl) Size() uint64 {
	c.RLock()
	if c.buf == nil {
		return uint64(0)
	} else {
		return uint64(len(c.buf))
	}
}

func (c *Ctl) Walk(name string) (d Dispatcher, err error) {
	err = os.ErrInvalid
	return
}
