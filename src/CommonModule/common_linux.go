package common

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func GetPartitionInfo(path string) (total, used, free uint64, err error) {
	// 获取分区信息
	fs := syscall.Statfs_t{}
	err = syscall.Statfs(path, &fs)
	if err != nil {
		return
	}

	//  获取分区使用情况
	total = fs.Blocks * uint64(fs.Bsize)
	free = fs.Bfree * uint64(fs.Bsize)
	used = total - free
	return
}

func CommondResult(commond string) (result string, err error) {
	out := exec.Command("/bin/bash", "-c", commond)
	data, err := out.Output()
	result = string(data)
	return
}

func CommondShell(commond string, argv []string) (err error) {
	attr := new(os.ProcAttr)
	newProcess, err := os.StartProcess(commond, argv, attr) //运行脚本
	if err != nil {
		return
	}

	status, err := newProcess.Wait() //等待命令执行完
	if !status.Success() {
		err = fmt.Errorf("exec %s fail", commond)
	}
	return
}
