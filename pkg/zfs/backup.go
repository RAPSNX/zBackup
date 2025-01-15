package zfs

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	OpenEBSPoolNameKey = "openebs.io/poolname"
)

var (
	pvCount      int
	datasetCount int
)

type BackupRunner struct {
	Client *kubernetes.Clientset
}

func (b *BackupRunner) Exec() error {
	ctx := context.Background()

	pvList, err := b.Client.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	datasetPathList := getDatasetPathList(pvList)

	pvCount = len(pvList.Items)
	datasetCount = len(datasetPathList)
	log.Printf("Found %d datasets out of %d PVs", datasetCount, pvCount)

	return nil
}

func getDatasetPathList(pvList *v1.PersistentVolumeList) (list []string) {
	for _, pv := range pvList.Items {
		var poolName string
		var datasetName string

		if pv.Spec.CSI != nil {
			poolName = pv.Spec.CSI.VolumeAttributes[OpenEBSPoolNameKey]
			datasetName = pv.Spec.CSI.VolumeHandle
		}
		if poolName != "" || datasetName != "" {
			list = append(list, fmt.Sprintf("%s/%s", poolName, datasetName))
		}
	}
	return list
}
