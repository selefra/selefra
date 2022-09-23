package oci

import (
	"github.com/selefra/selefra/pkg/utils"
	"testing"
)

func TestDownloadDB(t *testing.T) {
	RunDB()
}

func TestChangePort(t *testing.T) {
	tempDir, _ := utils.GetTempPath()
	confPath := tempDir + "/pgsql/data/postgresql.conf"
	ChangePort(confPath, "15432")
}
