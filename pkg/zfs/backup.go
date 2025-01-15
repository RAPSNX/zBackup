package zfs

import (
	"context"
	"fmt"

	"github.com/kataras/golog"
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

type BackupManager struct {
	KubeClient *kubernetes.Clientset
	Log        *golog.Logger
}

func (b *BackupManager) Start() error {
	ctx := context.Background()
	zfs := client{}

	b.Log.Info("Retriving list of PVs")
	pvList, err := b.KubeClient.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		b.Log.Error("Error while retriving list of PVs")
		return err
	}

	b.Log.Info("Extracting dataset paths out of PVs")
	datasetPathList := getDatasetPathList(pvList)

	pvCount = len(pvList.Items)
	datasetCount = len(datasetPathList)
	b.Log.Infof("Found %d datasets out of %d PVs", datasetCount, pvCount)

	// TODO: remove test
	fmt.Println(zfs.list())

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
