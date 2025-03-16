package zfs

import (
	"fmt"
	"io"
	"strings"

	"github.com/rgroemmer/zfs-backupper/pkg/utils"
)

type Dataset string
type Snapshot string

func ListDatasets() ([]Dataset, error) {
	rawOutput, err := utils.CMDWithOuput("zfs", "list", "-H", "-o", "name")
	if err != nil {
		return nil, err
	}
	return utils.SanitizeRawStringToList[Dataset](rawOutput), nil
}

func ListSnaphots() ([]Snapshot, error) {
	rawOutput, err := utils.CMDWithOuput("zfs", "list", "-t", "snap", "-H", "-o", "name")
	if err != nil {
		return nil, err
	}
	return utils.SanitizeRawStringToList[Snapshot](rawOutput), nil
}

func CreateSnapshot(ds Dataset) error {
	_, err := utils.CMDWithOuput("zfs", "snap", fmt.Sprintf("%s@%s", ds, "TEMP-TODO"))
	if err != nil {
		return err
	}
	return nil
}

func SendSnapshot(s Snapshot) (io.ReadCloser, error) {
	out, err := utils.CMDPipeOut("zfs", "send", string(s))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func DestroySnapshot(s Snapshot) error {
	if !strings.Contains(string(s), "@") {
		return fmt.Errorf("Error deleting snapshot doesnt contains '@', %s", s)
	}
	_, err := utils.CMDWithOuput("zfs", "destroy", string(s))
	if err != nil {
		return err
	}
	return nil
}
