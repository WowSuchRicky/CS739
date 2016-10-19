package nfs

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
)

func PathToGen(path string) (int, error) {
	cmd := exec.Command("./generation", path)

	var outb bytes.Buffer
	cmd.Stdout = &outb

	err := cmd.Run()

	if err != nil {
		fmt.Println("Error")
		return 0, err
	}

	out, err := strconv.Atoi(outb.String())
	return out, err
}
