package config

import (
	"time"
)

type ElectionConfig struct {
	// ReleaseOnCancel
	ReleaseOnCancel bool `yaml:"release_on_cancel"`
	// LeaseDuration
	LeaseDuration time.Duration `yaml:"lease_duration"`
	// RenewDeadline
	RenewDeadline time.Duration `yaml:"renew_deadline"`
	// RetryPeriod
	RetryPeriod time.Duration `yaml:"retry_period"`
}

func NewElectionConfig() *ElectionConfig {
	return &ElectionConfig{
		true,
		15,
		10,
		3,
	}
}
