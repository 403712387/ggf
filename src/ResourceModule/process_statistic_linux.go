package ResourceModule

import (
	"CommonModule"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"strconv"
	"strings"
	"time"
)

func processString(process *process.Process) (result string) {
	status, _ := process.Status()
	threadCount, _ := process.NumThreads()
	name, _ := process.Name()
	cpu, _ := process.CPUPercent()
	memory, _ := process.MemoryInfoEx()
	command, _ := process.Cmdline()
	second, _ := process.CreateTime()
	startup := time.Unix(second/1000, 0)
	return fmt.Sprintf("pid:%d, name:%s, status:%s, threads:%d, memory:%s, cpu:%f, command:%s, start:%s", process.Pid,
		name, status, threadCount, memory.String(), cpu, command, startup.Format("2006-01-02 15:04:05"))
}

// 统计进程使用CPU的情况
func statisticProcess(services map[string]bool) (result []common.ProcessInfo, err error) {

	// 获取服务组件的进程的信息(获取两次，中间相差一秒，是因为要获取磁盘IO，网络IO这些信息)
	firstProcesses, err := process.Processes()

	//因为nginx有父子进程，父进程不知道为啥没有使用cpu资源，所以删掉父进程pid
	res ,_ := common.CommondResult("ps -ef |grep nginx |grep -v grep |grep root |awk '{print $2}'")
	res1 ,_ := strconv.ParseInt(strings.Trim(res,"\n"),10,32)
	for id, info := range firstProcesses {
		if info.Pid == int32(res1) {
			firstProcesses = append(firstProcesses[:id], firstProcesses[id+1:]...)
			break
		}
	}

	firstProcesses, serviceName := getServiceInfo(services, firstProcesses)

	//第一次获取每个进程的cpu使用时间
	firstProcessCpuTime,err := getCpuPercent(serviceName)

	time.Sleep(time.Second)
	secondProcesses, err := process.Processes()

	//因为nginx有父子进程，父进程不知道为啥没有使用cpu资源，所以删掉父进程pid
	for id, info := range secondProcesses {
		if info.Pid == int32(res1) {
			secondProcesses = append(secondProcesses[:id], secondProcesses[id+1:]...)
			break
		}
	}
	secondProcesses, serviceName = getServiceInfo(services, secondProcesses)
	//第二次获取每个进程的cpu使用时间
	secondProcessCpuTime,_ := getCpuPercent(serviceName)

	//获取cpu的核数和新建个存储进程使用率的字典
	cpuCore,_ := common.CommondResult("cat /proc/cpuinfo |grep processor |wc -l")
	cpuCoreToInt,_ := strconv.ParseInt(strings.Trim(cpuCore,"\n"),10,64)
	processCpuResult := make(map[string]float64)
	//统计总的cpu使用率时长
	cpuTotalTime := secondProcessCpuTime["cpu"] - firstProcessCpuTime["cpu"]
	//获取每个进程的cpu占用率
	for firk,firv :=range firstProcessCpuTime {
		for seck,secv := range secondProcessCpuTime {
			if firk != seck || firk == "cpu" || seck == "cpu" {
				continue
			}
			resultCpu := float64(cpuCoreToInt * 100) * float64((secv - firv)) / float64(cpuTotalTime)
			processCpuResult[seck],_ = strconv.ParseFloat(fmt.Sprintf("%.2f",resultCpu),64)
		}
	}

	// 定义个float64来相加nginx线程的cpu占用率和内存占用率
	var nginxCpuResult float64
	var nginxMemResult float64
	Num := 1

	// 获取nginx的pid的个数
	pidNum, _ := common.CommondResult("ps -ef |grep nginx |grep -v grep |grep -v root |awk '{print $2}'")
	pidSlice := strings.Split(strings.Trim(pidNum,"\n"),"\n")
	// 获取服务组件的详细信息
	for _, first := range firstProcesses {
		for _, second := range secondProcesses {

			// 如果进程ID不相等，则不处理
			if first.Pid != second.Pid {
				continue
			}

			var info common.ProcessInfo
			pidName, _:= first.Name()
			if pidName == "nginx" {
				cpu, _ := first.CPUPercent()
				CPUresult := fmt.Sprintf("%.2f", cpu)
				CPUresultFloat, _ := strconv.ParseFloat(CPUresult, 64)
				nginxCpuResult += CPUresultFloat

				memoryUsed, _ := first.MemoryPercent()
				memoryUsedToFloat, _ := strconv.ParseFloat(fmt.Sprintf("%.2f",memoryUsed), 64)
				nginxMemResult += memoryUsedToFloat

				if Num == len(pidSlice) {
					info.Pid = first.Pid
					name, _ := serviceName[info.Pid]
					info.Name = name
					thread, _ := first.NumThreads()
					info.ThreadCount = thread
					info.CPU = nginxCpuResult
					millsecond, _ := first.CreateTime()
					info.Startup = time.Unix(millsecond/1000, 0).Format("2006-01-02 15:04:05")
					command, _ := first.Cmdline()
					info.Command = command

					// io使用情况
					info.Memory.UsedPrecent = nginxMemResult
					memory, _ := second.MemoryInfo()
					info.Memory.RSS = memory.RSS
					info.Memory.RSSReadable = common.ReabableSize(info.Memory.RSS)
					info.Memory.VMS = memory.VMS
					info.Memory.VMSReadable = common.ReabableSize(info.Memory.VMS)
					info.Memory.Data = memory.Data
					info.Memory.DataReadable = common.ReabableSize(info.Memory.Data)

					// io信息
					firstIO, _ := first.IOCounters()
					secondIO, _ := second.IOCounters()
					info.Disk.ReadCount = secondIO.ReadCount - firstIO.ReadCount
					info.Disk.WriteCount = secondIO.WriteCount - firstIO.WriteCount
					info.Disk.ReadByte = secondIO.ReadBytes - firstIO.ReadBytes
					info.Disk.ReadByteReadable = common.ReabableSize(info.Disk.ReadByte)
					info.Disk.WriteByte = secondIO.WriteBytes - firstIO.WriteBytes
					info.Disk.WriteByteReadable = common.ReabableSize(info.Disk.WriteByte)

					// 网络信息
					firstNetwork, _ := first.NetIOCounters(false)
					secondNetwork, _ := second.NetIOCounters(false)
					info.Network.Name = firstNetwork[0].Name
					info.Network.ReceiveByte = secondNetwork[0].BytesRecv - firstNetwork[0].BytesRecv
					info.Network.ReceiveByteReadable = common.ReabableSize(info.Network.ReceiveByte)
					info.Network.SendByte = secondNetwork[0].BytesSent - firstNetwork[0].BytesSent
					info.Network.SendByteReadable = common.ReabableSize(info.Network.SendByte)
					info.Network.ReceivePacket = secondNetwork[0].PacketsRecv - firstNetwork[0].PacketsRecv
					info.Network.SendPacket = secondNetwork[0].PacketsSent - firstNetwork[0].PacketsSent
					info.Network.ReceiveError = secondNetwork[0].Errin - firstNetwork[0].Errin
					info.Network.SendError = secondNetwork[0].Errout - firstNetwork[0].Errout
					info.Network.ReceiveDrop = secondNetwork[0].Dropin - firstNetwork[0].Dropin
					info.Network.SendDrop = secondNetwork[0].Dropout - firstNetwork[0].Dropout

					result = append(result, info)
				}
				Num ++

				break
			}
			// cpu使用情况
			info.Pid = first.Pid
			name, _ := serviceName[info.Pid]
			info.Name = name
			thread, _ := first.NumThreads()
			info.ThreadCount = thread
			//cpu, _ := getCpuPercent(info.Pid)
			info.CPU = processCpuResult[info.Name]
			millsecond, _ := first.CreateTime()
			info.Startup = time.Unix(millsecond/1000, 0).Format("2006-01-02 15:04:05")
			command, _ := first.Cmdline()
			info.Command = command

			// 内存使用情况
			memoryUsed, _ := second.MemoryPercent()
			memoryUsedToFloat, _ := strconv.ParseFloat(fmt.Sprintf("%.2f",memoryUsed), 64)
			info.Memory.UsedPrecent = memoryUsedToFloat
			memory, _ := second.MemoryInfo()
			info.Memory.RSS = memory.RSS
			info.Memory.RSSReadable = common.ReabableSize(info.Memory.RSS)
			info.Memory.VMS = memory.VMS
			info.Memory.VMSReadable = common.ReabableSize(info.Memory.VMS)
			info.Memory.Data = memory.Data
			info.Memory.DataReadable = common.ReabableSize(info.Memory.Data)

			// io信息
			firstIO, _ := first.IOCounters()
			secondIO, _ := second.IOCounters()
			info.Disk.ReadCount = secondIO.ReadCount - firstIO.ReadCount
			info.Disk.WriteCount = secondIO.WriteCount - firstIO.WriteCount
			info.Disk.ReadByte = secondIO.ReadBytes - firstIO.ReadBytes
			info.Disk.ReadByteReadable = common.ReabableSize(info.Disk.ReadByte)
			info.Disk.WriteByte = secondIO.WriteBytes - firstIO.WriteBytes
			info.Disk.WriteByteReadable = common.ReabableSize(info.Disk.WriteByte)

			// 网络信息
			firstNetwork, _ := first.NetIOCounters(false)
			secondNetwork, _ := second.NetIOCounters(false)
			info.Network.Name = firstNetwork[0].Name
			info.Network.ReceiveByte = secondNetwork[0].BytesRecv - firstNetwork[0].BytesRecv
			info.Network.ReceiveByteReadable = common.ReabableSize(info.Network.ReceiveByte)
			info.Network.SendByte = secondNetwork[0].BytesSent - firstNetwork[0].BytesSent
			info.Network.SendByteReadable = common.ReabableSize(info.Network.SendByte)
			info.Network.ReceivePacket = secondNetwork[0].PacketsRecv - firstNetwork[0].PacketsRecv
			info.Network.SendPacket = secondNetwork[0].PacketsSent - firstNetwork[0].PacketsSent
			info.Network.ReceiveError = secondNetwork[0].Errin - firstNetwork[0].Errin
			info.Network.SendError = secondNetwork[0].Errout - firstNetwork[0].Errout
			info.Network.ReceiveDrop = secondNetwork[0].Dropin - firstNetwork[0].Dropin
			info.Network.SendDrop = secondNetwork[0].Dropout - firstNetwork[0].Dropout

			result = append(result, info)

			// 打开文件的信息
			//info.OpenFile, _ = second.OpenFiles()

			// 进程的limit
			//info.Limit, _ = second.Rlimit()
			//logrus.Infof("second info :%v", info)
			// 进程的tcp/ip信息
			//info.Connection, _ = second.Connections()
		}
	}
	return
}

// 根据服务名称，从所有进程中获取服务组件的进程信息
func getServiceInfo(services map[string]bool, processes []*process.Process) (result []*process.Process, serviceName map[int32]string) {
	serviceName = make(map[int32]string)
	for _, proc := range processes {
		name, _ := proc.Name()
		commandLine, _ := proc.Cmdline()

		// 如果该进程属于监控进程，则忽略
		if strings.Index(commandLine, "Monitor_") > 0 {
			continue
		}

		// 如果该进程属于控制进程，则忽略
		if strings.Index(commandLine, "controller") > 0 {
			continue
		}

		//  判断进程是否属于服务
		isService := false
		service := ""
		for k, _ := range services {
			if strings.Index(commandLine,"sersync.jar")  >= 0 {
				isService = true
				service = "datasync"
				break
			}

			if strings.Index(commandLine, k+".jar") >= 0 || strings.Index(commandLine, "/mars/"+k) >= 0 || strings.Index(name, k) >= 0 {
				isService = true
				service = k
				break
			}
		}

		// 如果进程不属于服务，则忽略
		if !isService {
			continue
		}

		serviceName[proc.Pid] = service
		result = append(result, proc)

	}
	return
}

func getCpuPercent(services map[int32]string)(result map[string]int64, err error) {

	result = make(map[string]int64)
	for k,v := range services {
		processResult, _ := common.CommondResult(fmt.Sprintf("cat /proc/%s/stat |awk '{print $14,$15,$16,$17}' |awk '{b[NR]=$0; for(i=1;i<=NF;i++)a[NR]+=$i;}END{for(i=1;i<=NR;i++) print a[i]}'", fmt.Sprintf("%d",k)))
		//组成字典返回
		result[v],_ = strconv.ParseInt(strings.Trim(processResult,"\n"),10,64)
	}
	cpuResult, _ := common.CommondResult("cat /proc/stat |head -n 1 |awk '{b[NR]=$0; for(i=2;i<=NF;i++)a[NR]+=$i;}END{for(i=1;i<=NR;i++) print a[i]}'")
	result["cpu"],_ = strconv.ParseInt(strings.Trim(cpuResult,"\n"),10,64)
	return
}
