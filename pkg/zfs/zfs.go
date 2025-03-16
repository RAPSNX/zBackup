package zfs

import (
	"fmt"

	"github.com/rgroemmer/zfs-backupper/pkg/utils"
)

type Dataset string
type Snapshot string

func ListDatasets() ([]Dataset, error) {
	rawOutput, err := utils.RunCMD("zfs", "list", "-H", "-o", "name")
	if err != nil {
		return nil, err
	}
	return utils.SanitizeRawStringToList[Dataset](rawOutput), nil
}

func ListSnaphots() ([]Dataset, error) {
	rawOutput, err := utils.RunCMD("zfs", "list", "-t", "snap", "-H", "-o", "name")
	if err != nil {
		return nil, err
	}
	return utils.SanitizeRawStringToList[Dataset](rawOutput), nil
}

func CreateSnapshot(ds Dataset) error {
	_, err := utils.RunCMD("zfs", "snap", fmt.Sprintf("%s@%s", ds, "TEMP-TODO"))
	if err != nil {
		return err
	}
	return nil
}
