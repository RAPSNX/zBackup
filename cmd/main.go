package main

import (
	"log"

	"github.com/kataras/golog"
	"github.com/rgroemmer/zfs-backupper/pkg/zfs"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error running controller: %v", err)
	}
}

func run() error {
	log := golog.New()

	log.Info("Setting up rest-config")
	restConfig, err := config.GetConfig()
	if err != nil {
		return err
	}

	log.Info("Setting up kubernetes clientset")
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	log.Info("Starting BackupManager")
	return (&zfs.BackupManager{
		KubeClient: clientset,
		Log:        log,
	}).Start()
}
