package logout

import (
	"fmt"
	"testing"
)

func TestRunFunc(t *testing.T) {
	err := RunFunc(nil, nil)
	fmt.Println(err)
}
