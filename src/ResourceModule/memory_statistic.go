package ResourceModule

import (
	"CommonModule"
	"github.com/shirou/gopsutil/mem"
)

// 统计内存的使用情况
func statisticMemory() (info common.MemoryStatisticInfo, err error) {
	stat, err := mem.VirtualMemory()
	if err == nil {
		info.TotalCapacity = stat.Total
		info.UsedSize = stat.Used
		info.CacheSize = stat.Cached
		info.BufferSize = stat.Buffers
		info.AvailableSize = stat.Available

		info.TotalCapacityReadable = common.ReabableSize(info.TotalCapacity)
		info.UsedSizeReadable = common.ReabableSize(info.UsedSize)
		info.CacheSizeReadable = common.ReabableSize(info.CacheSize)
		info.BufferSizeReadable = common.ReabableSize(info.BufferSize)
		info.AvailableSizeReadable = common.ReabableSize(info.AvailableSize)
	}
	return
}
