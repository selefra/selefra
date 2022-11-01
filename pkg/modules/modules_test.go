package modules

import (
	"fmt"
	"testing"
)

func TestDownloadModule(t *testing.T) {
	s, e := DownloadModule("s3://static-picture/job3.yaml")
	fmt.Println(s, e)
}
