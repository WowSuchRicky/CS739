package nfs

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
)

func PathToGen(path string) (uint64, error) {
	cmd := exec.Command("./lib/generation", path)

	var outb bytes.Buffer
	cmd.Stdout = &outb

	err := cmd.Run()

	if err != nil {
		fmt.Println("Error")
		return 0, err
	}

	out, err := strconv.Atoi(outb.String())
	return uint64(out), err
}
