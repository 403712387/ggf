package ResourceModule

import (
	"CommonModule"
	"errors"
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
	diskNum,_ := common.CommondResult("fdisk -l |grep /dev/sd |grep Disk |wc -l")
	diskIoStat,_ := common.CommondResult(fmt.Sprintf("iostat -x -d -k 1 2 |grep sd |awk '{NF -=0}1'|tail -n %s", diskNum))
	diskIoStatSlice := strings.Split(strings.Trim(diskIoStat, "\n"), "\n")
	for _,disks := range diskIoStatSlice {
		diskIoStatSliceSec := strings.Split(disks, " ")
		for index,disk := range diskIoStatSliceSec {
			if index == 0 {
				continue
			}
			disk,_ := strconv.ParseFloat(disk,64)
			multimap[diskIoStatSliceSec[0]] = append(multimap[diskIoStatSliceSec[0]], uint64(disk))
		}
	}
	//logrus.Infof("getDiskIoInfo:%f",multimap)

	for k,v := range multimap {
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

//获取磁盘的容量，挂载状态，挂载点，磁盘类型，属性等
func statusDisk() (info []common.PartitionInfo, err error) {
	type DiskInfoMap map[string][]string //用来存放磁盘的信息

	mountDiskMap := make(DiskInfoMap) // 存放挂载的磁盘信息
	allDiskMap := make(DiskInfoMap)   // 存放所有的磁盘信息（包括挂载的和未挂载的）

	var diskstat common.PartitionInfo

	//获取全部磁盘信息，并转换成列表
	diskAll, _ := common.CommondResult("fdisk -l |grep 'Disk /dev' |grep -v Linux |grep -v swap |awk '{print $2$3$4}'")
	diskAllSlice := strings.Split(strings.Trim(diskAll, ",\n"), ",\n")

	//获取已挂载磁盘信息，并转换成列表
	diskMount, _ := common.CommondResult("df -h |grep -v devtmp |grep -v tmpfs |awk -v OFS=':' '{print $1,$NF,$2,$4}' |grep -v Filesystem")
	diskMountSlice := strings.Split(strings.Trim(diskMount, "\n"), "\n")

	//把全部磁盘信息塞入到字典
	for _, disks := range diskAllSlice {
		allDiskMap[strings.Split(disks, ":")[0]] = []string{strings.Split(disks, ":")[1]}
	}

	//把已挂载磁盘信息塞入到字典
	for _, mounts := range diskMountSlice {
		for i := 1; i < len(strings.Split(mounts, ":")); i++ {
			mountDiskMap[strings.Split(mounts, ":")[0]] = append(mountDiskMap[strings.Split(mounts, ":")[0]], strings.Split(mounts, ":")[i])
		}
	}

	//校验已挂载磁盘是否在全部磁盘字典中
	for key2, _ := range allDiskMap {
		isMount := false
		for key1, _ := range mountDiskMap {

			// 处理已经挂载的磁盘
			if strings.Contains(key1, key2) {
				isMount = true
				diskstat.Name = key1
				diskstat.MountPoint = mountDiskMap[key1][0]
				diskstat.Status = true
				diskstat.MountStatus = true
				diskstat.TotalCapacityReadable = mountDiskMap[key1][1]
				diskstat.AvailableCapacityReadable = mountDiskMap[key1][2]
				diskstat.Type = common.Disk_Local
				info = append(info, diskstat)
				break
			}
		}

		//  处理没有挂载的磁盘
		if !isMount {
			diskstat.Name = key2
			diskstat.MountPoint = " "
			diskstat.Status = true
			diskstat.MountStatus = false
			diskstat.TotalCapacityReadable = allDiskMap[key2][0]
			diskstat.AvailableCapacityReadable = allDiskMap[key2][0]
			diskstat.Type = common.Disk_Local
			info = append(info, diskstat)
		}
	}
	return
}

// 对磁盘进行一些操作（挂载，格式化等）
func operateDiskMount(conf common.OperationDisk) (err error){
	if conf.Type == "format" {
		if strings.Contains(conf.Name,"/dev/mapper/") {
			err = errors.New("formated")
			return
		}
		resultFir,_ := common.CommondResult(fmt.Sprintf("blkid %s", conf.Name))
		resultSec,_ := common.CommondResult(fmt.Sprintf("ls %s*", conf.Name))
		resultSecToSlice := strings.Split(strings.Trim(resultSec,"\n"),"\n")
		if strings.Contains(resultFir,"UUID") || len(resultSecToSlice) > 1  {
			err = errors.New("formated")
			return
		}
		if len(resultSecToSlice) == 1 {
			_, err1 := common.CommondResult(fmt.Sprintf("parted %s mklabel gpt -s", conf.Name))
			if err1 != nil {
				return
			}
			_, err1 = common.CommondResult(fmt.Sprintf("parted %s mkpart primary 1 100%%", conf.Name))
			if err1 != nil {
				return
			}
			_, err1 = common.CommondResult(fmt.Sprintf("mkfs.xfs %s1",conf.Name))
			if err1 != nil {
				return
			}
		}
	}
	if conf.Type == "mount" {
		if strings.Contains(conf.Name,"/dev/mapper/") {
			err = errors.New("mounted")
			return
		}
		resultFir,_ := common.CommondResult(fmt.Sprintf("blkid %s", conf.Name))
		resultSec,_ := common.CommondResult(fmt.Sprintf("ls %s*", conf.Name))
		resultSecToSlice := strings.Split(strings.Trim(resultSec,"\n"),"\n")
		if strings.Contains(resultFir,"UUID") || len(resultSecToSlice) > 1  {
			diskMountInfo ,_ := common.CommondResult("df -h |grep -v devtmp |grep -v tmpfs |awk -v OFS=':' '{print $1,$NF,$2,$4}' |grep -v Filesystem")
			diskMountInfoSlice := strings.Split(strings.Trim(diskMountInfo,"\n"),"\n")
			for _,info := range diskMountInfoSlice {
				if conf.Name == info || strings.Contains(info,conf.Name){
					err = errors.New("mounted")
					return
				}
			}
			diskMountDir ,_  := common.CommondResult("ls -d /data*")
			diskMountDirSlice := strings.Split(strings.Trim(diskMountDir,"\n"),"\n")
			if len(diskMountDirSlice) == 0 {
				_ ,err  := common.CommondResult("mkdir /data")
				if err != nil {
					return err
				}
				if strings.Contains(resultFir,"UUID") {
					_, err = common.CommondResult(fmt.Sprintf("mount %s /data", conf.Name))
					if err != nil {
						return err
					}
					UUID,_ := common.CommondResult(fmt.Sprintf("blkid %s |awk '{print $2}' |sed 's/\"//g'", conf.Name))
					_, err = common.CommondResult(fmt.Sprintf("echo '%s /data xfs defaults 0 0' >> /etc/fstab", strings.Trim(UUID,"\n")))
					if err != nil {
						return err
					}
				}
				if len(resultSecToSlice) > 1 {
					mountName := resultSecToSlice[len(resultSecToSlice)-1:][0]
					_, err = common.CommondResult(fmt.Sprintf("mount %s /data", mountName))
					if err != nil {
						return err
					}
					UUID,_ := common.CommondResult(fmt.Sprintf("blkid %s |awk '{print $2}' |sed 's/\"//g'", mountName))
					_, err = common.CommondResult(fmt.Sprintf("echo '%s /data xfs defaults 0 0' >> /etc/fstab", strings.Trim(UUID,"\n")))
					if err != nil {
						return err
					}
				}
			} else {
				lastMkDir := diskMountDirSlice[len(diskMountDirSlice)-1:]
				num,_ := strconv.Atoi(string(lastMkDir[0][len(lastMkDir[0])-1:]))
				_ ,err  := common.CommondResult(fmt.Sprintf("mkdir /data%d", num+1))
				if err != nil {
					return err
				}
				if strings.Contains(resultFir,"UUID") {
					_, err = common.CommondResult(fmt.Sprintf("mount %s /data%d", conf.Name,num+1))
					if err != nil {
						return err
					}
					UUID,_ := common.CommondResult(fmt.Sprintf("blkid %s |awk '{print $2}' |sed 's/\"//g'", conf.Name))
					_, err = common.CommondResult(fmt.Sprintf("echo '%s /data%d xfs defaults 0 0' >> /etc/fstab", strings.Trim(UUID,"\n"), num+1))
					if err != nil {
						return err
					}
				}
				if len(resultSecToSlice) > 1 {
					mountName := resultSecToSlice[len(resultSecToSlice)-1:][0]
					_, err = common.CommondResult(fmt.Sprintf("mount %s /data%d", mountName,num+1))
					if err != nil {
						return err
					}
					UUID,_ := common.CommondResult(fmt.Sprintf("blkid %s |awk '{print $2}' |sed 's/\"//g'", mountName))
					_, err = common.CommondResult(fmt.Sprintf("echo '%s /data%d xfs defaults 0 0' >> /etc/fstab", strings.Trim(UUID,"\n"), num+1))
					if err != nil {
						return err
					}
				}
			}
		}else{
			err = errors.New("not format,please format")
			return err
		}
	}
	return
}