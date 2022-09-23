package modules

import (
	"context"
	"github.com/selefra/selefra/pkg/internal/getter"
)

func DownloadModule(usePath string, modulesPath string) error {
	ctx := context.Background()
	err := getter.ModuleGet(ctx, modulesPath, usePath)
	return err
}
