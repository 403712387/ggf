package common

import "os/exec"

func GetPartitionInfo(path string) (total, used, free uint64, err error) {
	return
}

func CommondResult(commond string) (result string, err error) {
	cmd := exec.Command("cmd", "/C", commond)
	out, err := cmd.Output()
	if err == nil {
		result = string(out[:])
	}
	return
}

func CommondShell(commond string, argv []string) (err error) {
	exec.Command("cmd", "/C", commond)
	return
}
