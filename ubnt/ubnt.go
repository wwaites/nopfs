package ubnt

import (
	"bytes"
	"hubs.net.uk/sw/nopfs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"
)

func fixture() (data []byte, err error) {
	file, err := os.Open("aflist.txt") // For read access.
	if err != nil {
		return
	}
	data, err = ioutil.ReadAll(file)
	file.Close()
	return
}

var aflist_re *regexp.Regexp
var aflist_pat string =
	`^  (?P<k>[^ ][^.]+)\.+(?P<v>[^.].*)[ \t]*`
var aflist_prog string

func aflist_update(data []byte) {
	AfList.Lock()
	defer AfList.Unlock()
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		m := aflist_re.FindSubmatch(line)
		if m != nil {
			k := m[1]
			v := m[2]
//			AfList.AppendUnsafe(string(k), nopfs.NewFile(v))
			if bytes.HasPrefix(k, []byte("rxpower")) {
				AfList.AppendUnsafe(string(k), nopfs.NewFile(v))
			}
		}
		
	}
	AfList.AppendUnsafe("all", nopfs.NewFile(data))
}


func aflist(tick *time.Ticker) {
	for _ = range tick.C {
		cmd := exec.Command(aflist_prog)
		data, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("aflist: %s", err)
		} else {
			aflist_update(data)
		}
	}
}

var AfList *nopfs.Dir
var Dir *nopfs.Dir

func init() {
	aflist_re = regexp.MustCompile(aflist_pat)

	Dir = nopfs.NewDir()

	var err error
	aflist_prog, err = exec.LookPath("aflist")
	if err == nil {
		AfList = nopfs.NewDir()
		Dir.Append("aflist", AfList)

		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			aflist(ticker)
		}()
	}
}

