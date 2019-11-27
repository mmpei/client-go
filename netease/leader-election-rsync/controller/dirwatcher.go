package controller

import (
	"net/http"
	"os/exec"
	"strings"
)

type DirWatcher struct {
	Dir string
}

func NewDirWatcher(dir string) *DirWatcher {
	return &DirWatcher{
		Dir: dir,
	}
}

func (dw *DirWatcher) Handle(w http.ResponseWriter, r *http.Request) {
	cmdLine := "du -s " + dw.Dir
	cmd := exec.Command("/bin/sh", "-c", cmdLine)
	if res, err := cmd.Output(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		strArray := strings.Fields(strings.TrimSpace(string(res)))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strArray[0]))
		return
	}
}
