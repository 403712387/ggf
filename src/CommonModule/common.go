package common

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// 定义优先级
type Priority int32

const (
	Priority_First  Priority = iota // 第一优先级(最高优先级)
	Priority_Second                 // 第二优先级
	Priority_Third                  // 第三优先级
	Priority_Fourth                 // 第四优先级
	Priority_Fifth                  // 第五优先级(最低优先级)
	Priority_Count                  // 消息种类个数
)

// Priority转换为字符串
func (priority Priority) String() string {
	result := "third priority"

	switch priority {
	case Priority_First:
		result = "first priority"
	case Priority_Second:
		result = "second priority"
	case Priority_Third:
		result = "third priority"
	case Priority_Fourth:
		result = "fourth priority"
	case Priority_Fifth:
		result = "fifty priority"
	}

	return result
}

// ntp 操作类型
type NtpControlType string
const (
	Ntp_Control_Test NtpControlType = "test"
	Ntp_Control_Set  NtpControlType = "set"
)

// 服务组件信息
type ServiceInfo struct {
	Ip           string      // 服务器的ip
	Port         int32       // 服务器的普通端口
	HttpPort     int32       // 服务器的http端口
}

func (s *ServiceInfo) String() string {
	return fmt.Sprintf("ip:%s, http port:%d", s.Ip, s.HttpPort)
}

// 主机服务
type HostServiceInfo struct {
	ServiceInfo
}

func (o *HostServiceInfo) String() string {
	return fmt.Sprintf("%s", o.ServiceInfo.String())
}
// git信息
type GitInfo struct {
	Branch  string `json:"branch"`  // 分支
	Commit  string `json:"commit"`  // id
	Version string `json:"version"` // 版本信息
}

// NTP信息
type NTPInfo struct {
	IP                string `json:"ip"`                 // ntp服务器IP
	Port              int32  `json:"port"`               // ntp端口
	ProofreadInterval int64  `json:"proofread_interval"` // 校对间隔
	Enable            bool   `json:"enable"`             // 是否启用ntp校时
}

func (n *NTPInfo) String() string {
	return fmt.Sprintf("ip:%s, port:%d, proofread interval:%d, enable:%t", n.IP, n.Port, n.ProofreadInterval, n.Enable)
}

// 容量信息
type CapacityInfo struct {
	TotalCapacity             uint64 `json:"total_capacity"`              // 总容量
	TotalCapacityReadable     string `json:"total_capacity_readable"`     // 总容量（阅读友好的）
	UsedCapacity              uint64 `json:"used_capacity"`               // 已经使用的容量
	UsedCapacityReadable      string `json:"used_capacity_readable"`      // 已经使用的容量（阅读友好的)
	AvailableCapacity         uint64 `json:"available_capacity"`          // 空闲的容量
	AvailableCapacityReadable string `json:"available_capacity_readable"` // 空闲的容量（阅读友好的)
}

func (c *CapacityInfo) String() string {
	return fmt.Sprintf("total capacity:%d, total capacity readable:%s, used capacity:%d, used capacity readable:%s, available capacity:%d, available capacity readable:%s", c.TotalCapacity, c.TotalCapacityReadable, c.UsedCapacity, c.UsedCapacityReadable, c.AvailableCapacity, c.AvailableCapacityReadable)
}

// 网络统计信息
type NetworkStatisticInfo struct {
	Name                 string `json:"name"`                   //网口的名称
	ReceiveByte          uint64 `json:"receive_byte"`           //  接收的数据
	ReceiveByteReadable  string `json:"receive_byte_readable"`  // 接收的数据（易于读的）
	SendByte             uint64 `json:"send_byte"`              // 发送的数据
	SendByteReadable     string `json:"send_byte_readable"`     // 接收的数据（易于读的）
	CapacityByte         uint64 `json:"capacity_byte"`          // 网口的总容量
	CapacityByteReadable string `json:"capacity_byte_readable"` // 接收的数据（易于读的）
	ReceivePacket        uint64 `json:"receive_packet"`         // 接收的数据包
	SendPacket           uint64 `json:"send_packet"`            // 发送的数据包
	ReceiveError         uint64 `json:"receive_error"`          // 接收到的错误的包
	SendError            uint64 `json:"send_error"`             // 发送的错误的包
	ReceiveDrop          uint64 `json:"receive_drop"`           // 丢弃的接收的包
	SendDrop             uint64 `json:"send_drop"`              // 丢弃的发送的包
}

func (n *NetworkStatisticInfo) String() string {
	return fmt.Sprintf("name:%s, receive %d, receive readable:%s, send:%d, send readable:%s, capacity:%d, capacity readable:%s, receive packet:%d, send packet:%d, receive packet error:%d, send packet error:%d, receive packet drop:%d, send packet drop：%d",
		n.Name, n.ReceiveByte, n.ReceiveByteReadable, n.SendByte, n.SendByteReadable, n.CapacityByte, n.CapacityByteReadable,
		n.ReceivePacket, n.SendPacket, n.ReceiveError, n.SendError, n.ReceiveDrop, n.SendDrop)
}

// 磁盘统计
type DiskStatisticInfo struct {
	Rrqm         	  uint64 `json:"rrqm/s"`
	Wrqm   		 	  uint64 `json:"wrqm/s"`
	Read         	  uint64 `json:"r/s"`
	Write  		 	  uint64 `json:"w/s"`
	ReadByte          uint64 `json:"rkB/s"`
	ReadByteReadable  string `json:"read_byte_readable"`
	WriteByte         uint64 `json:"wkB/s"`
	WriteByteReadable string `json:"write_byte_readable"`
	Avgrq             uint64 `json:"avgrq-sz"`
	Avgqu         	  uint64 `json:"avgqu-sz"`
	Await    		  uint64 `json:"await"`
	Rawait            uint64 `json:"r_await"`
	Wawait       	  uint64 `json:"w_await"`
	Name              string `json:"name"`
	SerialNumber      string `json:"serialNumber"`
	Svctm             uint64 `json:"svctm"`
	Util			  uint64 `json:"util"`
}

func (d *DiskStatisticInfo) String() string {
	return fmt.Sprintf("name:%s, rrqm/s:%d, wrqm/s:%d, r/s:%d, w/s:%d, rkB/s:%d, wkB/s:%d, avgrq-sz:%d, avgqu-sz:%d, await:%d, r_await:%d, w_await:%d, svctm:%d, util:%d",
		d.Name, d.Rrqm, d.Wrqm, d.Read, d.Write, d.ReadByte, d.WriteByte, d.Avgrq, d.Avgqu, d.Await, d.Rawait, d.Wawait, d.Svctm, d.Util)
}

type ByDiskName []DiskStatisticInfo

func (a ByDiskName) Len() int           { return len(a) }
func (a ByDiskName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByDiskName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// 内存使用情况
type MemoryUseInfo struct {
	UsedPrecent  float64 `json:"used_precent"` // 内存占用的百分比
	RSS          uint64  `json:"rss"`          // 提交到物理内存中的内存使用
	RSSReadable  string  `json:"rss_readable"` // 提交到物理内存中的内存使用
	VMS          uint64  `json:"vms"`          // 虚拟内存
	VMSReadable  string  `json:"vms_readable"` // 可读性好的虚拟内存
	Data         uint64  `json:"data"`         // 数据段占用的内存，
	DataReadable string  `json:"data_readable"`
}

func (m *MemoryUseInfo) String() string {
	return fmt.Sprintf("rss:%d, rss readable:%s, vms:%d, vms readable:%s, data:%d, data readable:%s",
		m.RSS, m.RSSReadable, m.VMS, m.VMSReadable, m.Data, m.DataReadable)
}

// io信息
type IOInfo struct {
	ReadCount         uint64 `json:"read_count"`
	WriteCount        uint64 `json:"write_count"`
	ReadByte          uint64 `json:"read_byte"`
	ReadByteReadable  string `json:"read_byte_readable"`
	WriteByte         uint64 `json:"write_byte"`
	WriteByteReadable string `json:"write_byte_readable"`
}

func (i *IOInfo) String() string {
	return fmt.Sprintf("read count:%d, write count:%d, read byte:%d, read byte readable:%s, wirte byte:%d, wirte byte readable:%s",
		i.ReadCount, i.WriteCount, i.ReadByte, i.ReadByteReadable, i.WriteByte, i.WriteByteReadable)
}

// 进程的信息
type ProcessInfo struct {
	Name        string               `json:"service"`      // 服务名称
	Pid         int32                `json:"pid"`          // 进程名称
	ThreadCount int32                `json:"thread_count"` // 线程个数
	CPU         float64              `json:"cpu"`          // cpu的使用率
	Startup     string               `json:"startup"`      // 程序启动时间
	Command     string               `json:"command"`      // 启动参数
	Memory      MemoryUseInfo        `json:"memory"`       // 内存使用情况
	Disk        IOInfo               `json:"disk"`         // io使用
	Network     NetworkStatisticInfo `json:"network"`
	//OpenFile    []process.OpenFilesStat `json:"open_file"` // 打开文件的个数
	//Limit       []process.RlimitStat    `json:"limit"`     // 进程的limit
	//Connection []net.ConnectionStat `json:"connect"` // tcp/ip的信息
}

func (p *ProcessInfo) String() string {
	return fmt.Sprintf("name:%s, pid:%d, thread count:%d, cpu:%f, startup:%s, commond:%s, memory:%s, disk:%s, network:%s",
		p.Name, p.Pid, p.ThreadCount, p.CPU, p.Startup, p.Command, p.Memory.String(), p.Disk.String(), p.Network.String())
}

// 内存使用情况的统计
type MemoryStatisticInfo struct {
	TotalCapacity         uint64 `json:"total_capacity"`          // 总容量
	TotalCapacityReadable string `json:"total_capacity_readable"` // 总容量（阅读友好的）
	UsedSize              uint64 `json:"used_size"`               // 已经使用的容量
	UsedSizeReadable      string `json:"used_size_readable"`      // 已经使用的容量（阅读友好的)
	AvailableSize         uint64 `json:"available_size"`          // 可用的
	AvailableSizeReadable string `json:"available_size_readable"` // 可读性好的可用的内存
	BufferSize            uint64 `json:"buffer_size"`             // buffer的size
	BufferSizeReadable    string `json:"buffer_size_readable"`    // 可读性好的buffer size
	CacheSize             uint64 `json:"cache_size"`              // cache的size
	CacheSizeReadable     string `json:"cache_size_readable"`     // 可读性好的cache的size
}

func (m *MemoryStatisticInfo) String() string {
	return fmt.Sprintf("total:%d, total readable:%s, used:%d, used available:%s, available:%d, available readable:%s, buffer:%d, buffer readable:%s, cache:%d, cache readable:%s",
		m.TotalCapacity, m.TotalCapacityReadable, m.UsedSize, m.UsedSizeReadable, m.AvailableSize, m.AvailableSizeReadable, m.BufferSize, m.BufferSizeReadable, m.CacheSize, m.CacheSizeReadable)
}

// 日志信息
type LogInfo struct {
	LogLevel string `json:"level"` // 日志级别
}

func (l *LogInfo) String() string {
	return l.LogLevel
}

// 调试信息
type DebugInfo struct {
	Log LogInfo `json:"log"` // 日志信息
}

func (d *DebugInfo) String() string {
	return d.Log.String()
}

// 文件/路径是否存在
func IsExist(path string) (result bool) {
	if len(path) <= 0 {
		return false
	}

	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 是否为文件
func IsFile(path string) (result bool) {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// 是否为文件夹
func IsDirectory(path string) (result bool) {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 把byte转换为K,M,G,T等易读的字符串,小数点后保留一位有效数组
func ReabableSize(size uint64) (result string) {
	//  容量的单位
	util := [...]string{"B", "K", "M", "G", "T", "P", "E", "Z", "Y"}

	//  获取容量
	capacity := float64(size)
	for i := 0; i < len(util); i++ {
		if capacity < 1024 {
			result = fmt.Sprintf("%.1f%s", capacity, util[i])
			break
		}

		capacity = capacity / 1024
	}
	return
}

// 压缩文件
func CompressFile(absolutePath string) (data []byte, err error) {
	if !IsExist(absolutePath) {
		err = fmt.Errorf("not find path %s", absolutePath)
		return
	}

	// 开始进行zip压缩
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	// 压缩文件
	if IsFile(absolutePath) {
		err = writeFile(path.Dir(absolutePath), path.Base(absolutePath), writer)
	} else {
		err = writeFile(absolutePath, "", writer)
	}

	writer.Close()
	data = buf.Bytes()
	return
}

// 把文件写入压缩文件中
func writeFile(absoluteDir, relativePath string, writer *zip.Writer) (err error) {

	// 压缩文件的函数
	fileWriter := func(absoluteDir, relativePath string, writer *zip.Writer) {
		w, err := writer.Create(relativePath)
		if err == nil {
			data, _ := ioutil.ReadFile(path.Join(absoluteDir, relativePath))
			w.Write(data)
		} else {
			logrus.Errorf("compress file fail, error reason:%s", err.Error())
		}
	}

	if IsFile(path.Join(absoluteDir, relativePath)) {

		// 压缩文件
		fileWriter(absoluteDir, relativePath, writer)
	} else if IsDirectory(path.Join(absoluteDir, relativePath)) {

		// 读取目录下的子目录和文件信息
		files, e := ioutil.ReadDir(path.Join(absoluteDir, relativePath))
		if e != nil {
			err = e
			logrus.Error("read directory %s file, error reason:%s", relativePath, err.Error())
			return
		}

		// 遍历子目录和文件,并压缩
		for _, file := range files {
			tmpFile := path.Join(relativePath, file.Name())
			if IsFile(path.Join(absoluteDir, tmpFile)) {
				fileWriter(absoluteDir, tmpFile, writer)
			} else if IsDirectory(path.Join(absoluteDir, tmpFile)) {
				writeFile(absoluteDir, tmpFile, writer)
			}
		}
	}
	return
}

// 获取最大值
func MaxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// 获取最小值
func MinInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// 获取文本文件的内容（根据开始行数和结束行数）
func FileText(file string, beginLine, endLine int64) (begin, end int64, info []string, err error) {

	// 判断文件是否存在
	if !IsExist(file) {
		err = fmt.Errorf("no file %s", file)
		return
	}

	// 读取文件内容
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	// 获取文件的所有行数
	lines := strings.Split(string(data[:]), "\n")

	// 如果起始行和结束行都有指定且合法
	if beginLine >= 0 && endLine > 0 && beginLine > endLine {
		begin = MinInt64(beginLine, int64(len(lines)-1))
		end = MinInt64(endLine, int64(len(lines)-1))
		info = lines[begin:end]
		return
	}

	// 指定了起始行，没指定结束行
	if beginLine > 0 && endLine <= 0 {
		begin = MinInt64(beginLine, int64(len(lines)-1))
		endLine = begin + 100
		end = MinInt64(endLine, int64(len(lines)-1))
		info = lines[begin:end]
		return
	}

	// 起始行和结束行都没指定
	if beginLine <= 0 && endLine <= 0 {
		beginLine = int64(len(lines) - 101)
		endLine = beginLine + 100
		begin = MaxInt64(0, beginLine)
		end = MinInt64(endLine, int64(len(lines)-1))
		info = lines[begin:end]
		return
	}

	err = fmt.Errorf("invalid range, begin:%d, end:%d, file lines:%d", beginLine, endLine, len(lines))
	return
}

func Pow(x uint64, y int) (result uint64) {
	result = 1
	for i := 0; i < y; i++ {
		result = result * x
	}
	result = result / 8
	return
}
