package controller

import "sync"

type SwitchInterface interface {
	ToMaster() error
	ToSlave() error
}

type MasterController struct {
	sync.RWMutex
	master bool
	delegates []SwitchInterface
}

func NewMasterController(d []SwitchInterface) *MasterController {
	return &MasterController{
		delegates: d,
	}
}

func (mc *MasterController) ToMaster() error {
	mc.Lock()
	defer mc.Unlock()

	mc.master = true
	for _, s := range mc.delegates {
		s.ToMaster()
	}
	return nil
}

func (mc *MasterController) ToSlave() error {
	mc.Lock()
	defer mc.Unlock()

	mc.master = false
	for _, s := range mc.delegates {
		s.ToSlave()
	}
	return nil
}

func (mc *MasterController) IsMaster() bool {
	mc.RLock()
	defer mc.RUnlock()

	return mc.master
}
