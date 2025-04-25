package zfs

import (
	"fmt"
	"io"
	"strings"

	"github.com/kataras/golog"
	"github.com/rgroemmer/zfs-backupper/pkg/utils"
)

type Dataset string
type Snapshot string

type zfsClient struct {
	log *golog.Logger
}

func NewZfsClient(logger *golog.Logger) *zfsClient {
	return &zfsClient{
		log: logger,
	}
}

func (c *zfsClient) ListDatasets() ([]Dataset, error) {
	c.log.Info("List datasets")
	rawOutput, err := utils.CMDWithOuput("zfs", "list", "-H", "-o", "name")
	if err != nil {
		return nil, err
	}
	return utils.SanitizeRawStringToList[Dataset](rawOutput), nil
}

func (c *zfsClient) ListSnaphots() ([]Snapshot, error) {
	rawOutput, err := utils.CMDWithOuput("zfs", "list", "-t", "snap", "-H", "-o", "name")
	if err != nil {
		return nil, err
	}
	return utils.SanitizeRawStringToList[Snapshot](rawOutput), nil
}

func (c *zfsClient) CreateSnapshot(ds Dataset) error {
	_, err := utils.CMDWithOuput("zfs", "snap", fmt.Sprintf("%s@%s", ds, "TEMP-TODO"))
	if err != nil {
		return err
	}
	return nil
}

func (c *zfsClient) SendSnapshot(s Snapshot) (io.ReadCloser, error) {
	c.log.Info("Sending snapshot stream")
	o, err := utils.CMDPipeOut("zfs", "send", string(s))
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (c *zfsClient) DestroySnapshot(s Snapshot) error {
	if !strings.Contains(string(s), "@") {
		return fmt.Errorf("error deleting snapshot doesnt contains '@', %s", s)
	}
	_, err := utils.CMDWithOuput("zfs", "destroy", string(s))
	if err != nil {
		return err
	}
	return nil
}
