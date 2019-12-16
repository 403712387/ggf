package ResourceModule

import (
	"CommonModule"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 磁盘读写速度
func statisticDisk() (info []common.DiskStatisticInfo, err error) {
	//firstStat, err := disk.IOCounters()
	time.Sleep(time.Second)
	//secondStat, err := disk.IOCounters()
	//
	//// 统计读写速度
	//for key, first := range firstStat {
	//
	//	// 普通的磁盘以sd或者hd开头
	//	if strings.Index(first.Name, "sd") < 0 && strings.Index(first.Name, "hd") < 0 {
	//		continue
	//	}
	//
	//	// 判断是否两次读取都有改磁盘信息
	//	if _, OK := secondStat[key]; !OK {
	//		continue
	//	}
	//
	//	second := secondStat[key]
	//	var stat common.DiskStatisticInfo
	//	stat.Name = first.Name
	//	stat.IopsInProgress = second.IopsInProgress
	//	stat.IoTime = second.IoTime - first.IoTime
	//	stat.Label = second.Label
	//	stat.MergedReadCount = second.MergedReadCount - first.MergedReadCount
	//	stat.ReadCount = second.ReadCount - first.ReadCount
	//	stat.ReadByte = second.ReadBytes - first.ReadBytes
	//	stat.ReadByteReadable = common.ReabableSize(stat.ReadByte)
	//	stat.ReadTime = second.ReadTime - first.ReadTime
	//	stat.MergedWriteCount = second.MergedWriteCount - first.MergedWriteCount
	//	stat.WriteCount = second.WriteCount - first.WriteCount
	//	stat.WriteByte = second.WriteBytes - first.WriteBytes
	//	stat.WriteByteReadable = common.ReabableSize(stat.WriteByte)
	//	stat.WriteTime = second.WriteTime - first.WriteTime
	//	stat.WeightedIO = second.WeightedIO
	//	stat.SerialNumber = first.SerialNumber
	//
	//	info = append(info, stat)
	//}
	multimap := make(map[string][]uint64)
	diskNum, _ := common.CommondResult("fdisk -l |grep /dev/sd |grep Disk |wc -l")
	diskIoStat, _ := common.CommondResult(fmt.Sprintf("iostat -x -d -k 1 2 |grep sd |awk '{NF -=0}1'|tail -n %s", diskNum))
	diskIoStatSlice := strings.Split(strings.Trim(diskIoStat, "\n"), "\n")
	for _, disks := range diskIoStatSlice {
		diskIoStatSliceSec := strings.Split(disks, " ")
		for index, disk := range diskIoStatSliceSec {
			if index == 0 {
				continue
			}
			disk, _ := strconv.ParseFloat(disk, 64)
			multimap[diskIoStatSliceSec[0]] = append(multimap[diskIoStatSliceSec[0]], uint64(disk))
		}
	}
	//logrus.Infof("getDiskIoInfo:%f",multimap)

	for k, v := range multimap {
		var stat common.DiskStatisticInfo
		stat.Name = k
		stat.Rrqm = v[0]
		stat.Wrqm = v[1]
		stat.Read = v[2]
		stat.Write = v[3]
		stat.ReadByte = v[4]
		stat.WriteByte = v[5]
		//stat.ReadByteReadable = common.ReabableSize(stat.ReadByte)
		stat.Avgrq = v[6]
		stat.Avgqu = v[7]
		//stat.WriteByteReadable = common.ReabableSize(stat.WriteByte)
		stat.Await = v[8]
		stat.Rawait = v[9]
		stat.Wawait = v[10]
		stat.Svctm = v[11]
		//stat.SerialNumber = ""
		stat.Util = v[12]

		info = append(info, stat)
	}

	sort.Sort(common.ByDiskName(info))
	return
}
