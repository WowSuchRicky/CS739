package nfs

import (
	"fmt"
	"testing"
)

func TestGen(t *testing.T) {
	st, _ := PathToGen("./gen.go")
	fmt.Printf("Hi %v\n", st)
}
