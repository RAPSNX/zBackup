package main

import (
	"log"

	"github.com/kataras/golog"
	"github.com/rgroemmer/zfs-backupper/pkg/zfs"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error executing BackupRunner: %v", err)
	}
}

func run() error {
	log := golog.New()

	// Get dataset list
	log.Info("Listing datasets")
	datasets, err := zfs.ListDatasets()
	if err != nil {
		return err
	}

	// Create snapshots forEach
	log.Info("Createing snapshot for each dataset")
	for _, ds := range datasets {
		err := zfs.CreateSnapshot(ds)
		if err != nil {
			return err
		}
	}

	// Pipe to restic forEach snap

	return nil
}
