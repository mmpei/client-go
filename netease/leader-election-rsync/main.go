/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/transport"
	"k8s.io/klog"
	"k8s.io/client-go/netease/leader-election-rsync/election"
	"k8s.io/client-go/netease/leader-election-rsync/controller"
)

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	klog.InitFlags(nil)

	var kubeconfig string
	var cmLockName string
	var cmLockNamespace string
	var id string
	var rsyncFile string
	var httpServer string
	var watchDir string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&id, "id", "", "the holder identity name")
	flag.StringVar(&rsyncFile, "rsync-file", "/var/leader-election-rsyncfile", "the rsync file to control the rsync action")
	flag.StringVar(&cmLockName, "config-map-lock-name", "leader-election-rsync", "the configmap lock resource name")
	flag.StringVar(&cmLockNamespace, "config-map-lock-namespace", "default", "the configmap lock resource namespace")
	flag.StringVar(&httpServer,"http-server", "0.0.0.0:29991", "http server for election")
	flag.StringVar(&watchDir,"watch-dir", "/data", "the dir should be watched")
	flag.Parse()

	if id == "" {
		klog.Fatal("unable to get id (missing id flag).")
	}

	// generate master controller
	var switchList []controller.SwitchInterface
	rsyncCtl := controller.NewFileController(rsyncFile)
	switchList = append(switchList, rsyncCtl)
	masterCtl := controller.NewMasterController(switchList)

	// start metrics server
	stat, closer, err := controller.NewPrometheusScope(httpServer)
	if err != nil {
		klog.Error("failed to initialize metrics: %s", err)
	}
	defer closer.Close()
	dirWatcher := controller.NewDirWatcher(watchDir, stat)
	dirWatcher.Start()

	//http.HandleFunc("/data/du", controller.NewDirWatcher(watchDir).Handle)
	//go http.ListenAndServe(httpServer, nil)

	// leader election uses the Kubernetes API by writing to a
	// lock object, which can be a LeaseLock object (preferred),
	// a ConfigMap, or an Endpoints (deprecated) object.
	// Conflicting writes are detected and each client handles those actions
	// independently.
	config, err := buildConfig(kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	// use a Go context so we can tell the leaderelection code when we
	// want to step down
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ele := election.NewElection(ctx, config, id, cmLockNamespace, cmLockName, masterCtl)

	// use a client that will stop allowing new requests once the context ends
	config.Wrap(transport.ContextCanceller(ctx, fmt.Errorf("the leader is shutting down")))

	// listen for interrupts or the Linux SIGTERM signal and cancel
	// our context, which the leader election code will observe and
	// step down
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		log.Printf("Received termination, signaling shutdown")
		cancel()
	}()

	// start the leader election code loop
	ele.Run()
}
