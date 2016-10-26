package nfs

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

func InumToPath(inum int) (string, error) {
	cmd := exec.Command("./lib/inum.sh", strconv.Itoa(inum))

	var outb bytes.Buffer
	cmd.Stdout = &outb

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(outb.String()), err
}
