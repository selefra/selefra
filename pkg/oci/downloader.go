package oci

import (
	"fmt"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/selefra/selefra/pkg/utils"
	"path/filepath"
)

func DownloadDB() {
	ociPath, err := utils.GetOCIPath()
	if err != nil {
		return
	}
	image, err := crane.Pull("docker.io/library/postgres:latest")
	fmt.Println(image.Size())
	if err != nil {
		return
	}
	err = crane.SaveOCI(image, filepath.Join(ociPath, "postgres"))
	if err != nil {
		return
	}
	err = crane.SaveLegacy(image, "docker.io/library/postgres:latest", filepath.Join(ociPath, "postgres.tar"))
	if err != nil {
		fmt.Println(err)
		return
	}
	img, err := crane.Load(filepath.Join(ociPath, "postgres.tar"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(img.Size())
}
