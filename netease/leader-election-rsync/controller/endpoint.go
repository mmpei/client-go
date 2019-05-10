package controller

import (
	"os"
	"k8s.io/klog"
	"k8s.io/client-go/netease/leader-election-rsync/utils"
)

// FileSwitch will touch a file or rm a file when triggered
// this will be a switch for other process or shell script
type FileSwitch struct {
	filePath string
}

func NewFileController(filePath string) SwitchInterface {
	if ok := utils.FileExists(filePath); ok {
		klog.Warningf("file has exists, removed: %s", filePath)
		os.Remove(filePath)
	}
	return &FileSwitch{
		filePath: filePath,
	}
}

func (fs *FileSwitch) ToMaster() error {
	klog.Infof("to master, touch file: %s", fs.filePath)
	if ok := utils.FileExists(fs.filePath); ok {
		klog.Warningf("to master will touch the file, but has existed: %s", fs.filePath)
		return nil
	}
	os.Create(fs.filePath)
	return nil
}

func (fs *FileSwitch) ToSlave() error {
	klog.Infof("to slave, remove file: %s", fs.filePath)
	os.Remove(fs.filePath)
	return nil
}
