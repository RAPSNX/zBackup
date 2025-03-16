package main

import (
	"os"

	"github.com/kataras/golog"
	"github.com/rgroemmer/zfs-backupper/pkg/utils"
	"github.com/rgroemmer/zfs-backupper/pkg/zfs"
)

func main() {
	log := golog.New()
	if err := run(log); err != nil {
		log.Errorf("Error executing BackupRunner: %v", err)
	}
}

func run(log *golog.Logger) error {
	pipeout, err := zfs.SendSnapshot("kubex-main/pvc-eddec1f4-4792-4609-a596-8a197d66baa2@test")
	if err != nil {
		return err
	}

	out, err := utils.CMDPipeIn("tee", pipeout, "temp")
	if err != nil {
		return err
	}

	log.Info(out)

	os.Exit(0)
	// Get dataset list
	log.Info("Listing datasets")
	datasets, err := zfs.ListDatasets()
	if err != nil {
		return err
	}

	// Create snapshots forEach
	for _, ds := range datasets {
		log.Info("Creating snap for ", ds)
		err := zfs.CreateSnapshot(ds)
		if err != nil {
			return err
		}
	}

	// List all snaps
	snaps, err := zfs.ListSnaphots()
	if err != nil {
		return err
	}

	// Backup all snaps
	// for _, snap := range snaps {
	// 	log.Info("Backup snaps now")
	// 	err := zfs.SendSnapshot(snap)
	// 	if err != nil {
	// 		panic("lul")
	// 	}
	// }

	// Destroy all snaps
	for _, snap := range snaps {
		log.Info("Destroy snap ", snap)
		err := zfs.DestroySnapshot(snap)
		if err != nil {
			return err
		}
	}

	// Pipe to restic forEach snap

	return nil
}
