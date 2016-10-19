package nfs

import (
	"bytes"
	"strconv"
	"os/exec"
)

func InumToPath(inum int) (string, error) {
	cmd := exec.Command("./lib/inum.sh", strconv.Itoa(inum))

	var outb bytes.Buffer
	cmd.Stdout = &outb;

	err := cmd.Run()

	if err != nil{
	    return "", err
	}
	
	return outb.String(), err
}
