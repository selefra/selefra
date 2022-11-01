package modules

import (
	"context"
	"github.com/selefra/selefra/pkg/internal/getter"
)

func DownloadModule(usePath string) (string, error) {
	ctx := context.Background()
	err := getter.Get(ctx, "./test", usePath)
	return "", err
}
