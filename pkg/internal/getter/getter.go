package getter

import (
	"context"
	"github.com/selefra/selefra/ui/progress"
	"os"
	"time"

	getter "github.com/hashicorp/go-getter"
)

type Detector struct {
	Name     string
	Detector getter.Detector
}

var (
	detectors = []getter.Detector{
		new(getter.GitHubDetector),
		new(getter.GitDetector),
		new(getter.S3Detector),
		new(getter.GCSDetector),
		new(getter.FileDetector),
	}

	decompressors = map[string]getter.Decompressor{
		"bz2": new(getter.Bzip2Decompressor),
		"gz":  new(getter.GzipDecompressor),
		"xz":  new(getter.XzDecompressor),
		"zip": new(getter.ZipDecompressor),

		"tar.bz2":  new(getter.TarBzip2Decompressor),
		"tar.tbz2": new(getter.TarBzip2Decompressor),

		"tar.gz": new(getter.TarGzipDecompressor),
		"tgz":    new(getter.TarGzipDecompressor),

		"tar.xz": new(getter.TarXzDecompressor),
		"txz":    new(getter.TarXzDecompressor),
	}

	getters = map[string]getter.Getter{
		"file":   new(getter.FileGetter),
		"gcs":    new(getter.GCSGetter),
		"github": new(getter.GitGetter),
		"git":    new(getter.GitGetter),
		"hg":     new(getter.HgGetter),
		"s3":     new(getter.S3Getter),
		"http":   httpGetter,
		"https":  httpGetter,
	}
)

var httpGetter = &getter.HttpGetter{
	ReadTimeout:           30 * time.Second,
	MaxBytes:              500000000,
	XTerraformGetDisabled: true,
}

func Get(ctx context.Context, installPath, url string, options ...getter.ClientOption) error {
	pwd, _ := os.Getwd()
	client := getter.Client{
		Src:           url,
		Dst:           installPath,
		Pwd:           pwd,
		Mode:          getter.ClientModeDir,
		Detectors:     detectors,
		Decompressors: decompressors,
		Getters:       getters,
		Ctx:           ctx,
		// Extra options provided by caller to overwrite default behavior
		Options:          options,
		ProgressListener: progress.CreateProgress(),
	}

	if err := client.Get(); err != nil {
		return err
	}
	return nil
}

func ModuleGet(ctx context.Context, installPath, url string, options ...getter.ClientOption) error {
	pwd, _ := os.Getwd()
	client := getter.Client{
		Src:           url,
		Dst:           installPath,
		Pwd:           pwd,
		Mode:          getter.ClientModeDir,
		Detectors:     detectors,
		Decompressors: decompressors,
		Getters:       getters,
		Ctx:           ctx,
		// Extra options provided by caller to overwrite default behavior
		Options: options,
	}

	if err := client.Get(); err != nil {
		return err
	}
	return nil
}
