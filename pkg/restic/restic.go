package restic

import (
	"io"

	"github.com/kataras/golog"
	"github.com/rgroemmer/zfs-backupper/pkg/utils"
)

type restic struct {
	log *golog.Logger
}

func NewRestic(logger *golog.Logger) *restic {
	return &restic{
		log: logger,
	}
}

func (c *restic) NewBackup(in io.ReadCloser) error {
	c.log.Info("Recive snapshot stream, backup via stdin")
	err := utils.CMDPipeIn("restic", in, "backup", "--stdin")
	if err != nil {
		return err
	}

	return nil
}
