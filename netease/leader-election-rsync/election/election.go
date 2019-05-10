package election

import (
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/rest"
	"context"
	clientset "k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/netease/leader-election-rsync/config"
	"time"
	"k8s.io/klog"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/netease/leader-election-rsync/controller"
)

type Interface interface {
	Run() error
}

type Election struct {
	ctx context.Context
	id string
	lock *resourcelock.ConfigMapLock
	config *config.ElectionConfig
	ctl controller.SwitchInterface
}

func NewElection(ctx context.Context, cfg *rest.Config, id, ns, name string, ctl controller.SwitchInterface) *Election {
	client := clientset.NewForConfigOrDie(cfg)

	// use ConfigMapLock for version equal or less than 1.11
	lock := &resourcelock.ConfigMapLock{
		ConfigMapMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Client: client.CoreV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}
	return &Election{
		ctx: ctx,
		id: id,
		lock: lock,
		config: config.NewElectionConfig(),
		ctl: ctl,
	}
}

func (el *Election) Run() error {
	// start the leader election code loop
	leaderelection.RunOrDie(el.ctx, leaderelection.LeaderElectionConfig{
		Lock: el.lock,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: el.config.ReleaseOnCancel,
		LeaseDuration:   el.config.LeaseDuration * time.Second,
		RenewDeadline:   el.config.RenewDeadline * time.Second,
		RetryPeriod:     el.config.RetryPeriod * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				// we're notified when we start - this is where you would
				// usually put your code
				klog.Infof("%s: leading", el.id)
				el.ctl.ToMaster()
			},
			OnStoppedLeading: func() {
				// we can do cleanup here, or after the RunOrDie method
				// returns
				klog.Infof("%s: lost", el.id)
				el.ctl.ToSlave()
			},
			OnNewLeader: func(identity string) {
				// we're notified when new leader elected
				if identity == el.id {
					// I just got the lock
					return
				}
				klog.Infof("new leader elected: %v", identity)
			},
		},
	})
	return nil
}
