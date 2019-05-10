package controller

import (
	"os"
	"k8s.io/klog"
)

type EnvController struct {
	envName string
}

func NewEnvController(name string) SwitchInterface {
	return &EnvController{
		envName: name,
	}
}

func (ec *EnvController) ToMaster() error {
	os.Setenv(ec.envName, "1")
	return nil
}

func (ec *EnvController) ToSlave() error {
	value := os.Getenv(ec.envName)
	if len(value) == 0 {
		klog.Warning("switch to slave, env not present as master")
	}
	os.Unsetenv(ec.envName)
	return nil
}
