package main

import (
	"log"

	"github.com/rgroemmer/zfs-backupper/pkg/zfs"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error running controller: %v", err)
	}
}

func run() error {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info("Setting up rest-config")
	restConfig, err := config.GetConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	return (&zfs.BackupRunner{Client: clientset}).Exec()
}
