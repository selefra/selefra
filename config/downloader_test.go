package config

import (
	"fmt"
	"testing"
)

func TestDownloader_Get(t *testing.T) {
	d := Downloader{Url: "s3://static-picture/job3.yaml"}
	b, err := d.Get()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(b)
	d1 := Downloader{
		Url: "https://static-picture.s3.amazonaws.com/job3.yaml",
	}
	b1, err1 := d1.Get()
	if err1 != nil {
		t.Error(err1)
	}
	fmt.Println(b1)
}
