package layer

import (
	"fmt"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	timeout := time.Second * 30
	layer := NewLayer("", "", "1.0", timeout)

	p := Parameters{
		dedupe: nil,
		path:   "users/1234-5678-abcd/conversations",
	}

	res, err := layer.request("GET", p)
	fmt.Println(err)
	fmt.Println(res)
}
