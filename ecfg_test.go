package ecfg

import (
	"fmt"
	"testing"
)

func Test_Some(t *testing.T) {
	var some = struct {
		Some struct {
			SomeString string
		} `ecfg:"some"`
	}{}

	if err := walkThroughStruct(&some); err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%#v\n", some)
}
