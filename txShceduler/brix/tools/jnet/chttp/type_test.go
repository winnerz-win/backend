package chttp

import (
	"fmt"
	"net/url"
	"testing"
)

func TestTTT(t *testing.T) {
	kv := MakeHTTPBody()

	fmt.Println(kv.String())

	type Dm struct {
		a string
		b int
	}
	dm := Dm{
		"ssss",
		1,
	}

	u := url.Values{}

	u.Set("c", "1")
	u.Set("c", fmt.Sprintf("%v", dm))
	fmt.Println(u.Get("c"))
	u.Set("ass", "1")
	u.Set("z", "1")
	fmt.Println(u.Encode())

}
