package ResourceModule

import (
	"CommonModule"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

// 统计CPU的使用情况
func statisticCPU() (info common.CapacityInfo, err error) {
	used, err := cpu.Percent(time.Second, false)
	if err == nil {
		count, _ := cpu.Counts(true)
		info.TotalCapacity = uint64(count)               // CPU的总核数
		info.UsedCapacity = uint64(used[0])              //  已经使用的CPU
		info.AvailableCapacity = 100 - info.UsedCapacity // 空闲的CPU
	}
	return
}
