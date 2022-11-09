package logout

import (
	"fmt"
	"testing"
)

func TestRunFunc(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}
	err := RunFunc(nil, nil)
	fmt.Println(err)
}
