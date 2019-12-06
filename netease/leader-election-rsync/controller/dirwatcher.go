package controller

import (
	"net/http"
	"os/exec"
	"strings"
	"github.com/uber-go/tally"
	"strconv"
	"k8s.io/klog"
	"time"
)

type DirWatcher struct {
	Scope tally.Scope
	Dir string
}

func NewDirWatcher(dir string, scope tally.Scope) *DirWatcher {
	return &DirWatcher{
		Scope: scope,
		Dir: dir,
	}
}

func (dw *DirWatcher) Handle(w http.ResponseWriter, r *http.Request) {
	size, err := dw.du()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatInt(size,10)))
		return
	}
}

func (dw *DirWatcher) Start() {
	go func() {
		for {
			dw.Watch()
			time.Sleep(5 * time.Second)
		}
	}()
}

func (dw *DirWatcher) Watch() {
	size, err := dw.du()
	if err != nil {
		klog.Errorf("du error: v%", err)
	} else {
		dw.Scope.Tagged(map[string]string{
			"service":   "harbor",
			"module": "rsync",
		}).Gauge("dir_du").Update(float64(size))
	}
}

func (dw *DirWatcher) du() (int64, error ) {
	cmdLine := "du -s " + dw.Dir
	cmd := exec.Command("/bin/sh", "-c", cmdLine)
	if res, err := cmd.Output(); err != nil {
		return 0, err
	} else {
		strArray := strings.Fields(strings.TrimSpace(string(res)))
		return strconv.ParseInt(strArray[0], 10, 64)
	}
}
