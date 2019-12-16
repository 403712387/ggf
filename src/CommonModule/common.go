package common

import (
	"archive/zip"
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
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

// 服务状态
type ServiceStatus string

const (
	Service_Start ServiceStatus = "start" // 在线
	Service_Stop  ServiceStatus = "stop"  // 离线
)

// 服务器的类型
type ServiceType int32

const (
	Service_Invalid_Type   ServiceType = iota //无效的类型
	Service_Collect                           // 采集服务
	Service_Management                        // 管理服务
	Service_Search                            // 搜索服务
	Service_Business                          // 应用服务
	Service_Analysis                          // 分析服务
	Service_Major_Database                    // 主数据库
	Service_Import                            // 接入服务
	Service_Export                            // 接出服务
	Service_Minor_Database                    // 备用数据库
	Service_SubPlatform                       // 下级子平台
	Service_Kafka                             // kafka消息服务器
	Service_Nginx                             // nginx反向代理
	Service_Storage                           // 存储服务
	Service_Host                              // host服务
)

// 服务类型
func (service ServiceType) String() string {
	result := "host service"
	switch service {
	case Service_Invalid_Type:
		result = "invalid service"
	case Service_Collect:
		result = "collect service"
	case Service_Management:
		result = "manager service"
	case Service_Search:
		result = "search service"
	case Service_Business:
		result = "business service"
	case Service_Analysis:
		result = "analysis service"
	case Service_Major_Database:
		result = "major database service"
	case Service_Import:
		result = "import service"
	case Service_Export:
		result = "export service"
	case Service_Minor_Database:
		result = "minor database"
	case Service_SubPlatform:
		result = "subplatform"
	case Service_Kafka:
		result = "kafka service"
	case Service_Nginx:
		result = "nginx service"
	case Service_Storage:
		result = "storage service"
	case Service_Host:
		result = "host service"
	}
	return result
}

// 服务的角色
type ServiceRole int32

const (
	Service_Leader ServiceRole = iota // leader角色
	Service_Worker                    // worker角色
	Service_Parent                    // parent角色
	Service_Child                     // child角色
)

func (role ServiceRole) String() string {
	result := "leader"
	switch role {
	case Service_Leader:
		result = "host leader"
	case Service_Child:
		result = "host child"
	case Service_Worker:
		result = "host worker"
	case Service_Parent:
		result = "host parent"
	}
	return result
}

type OperateType string

const (
	Control_Stop    OperateType = "stop"    // 停止
	Control_Start   OperateType = "start"   // 启动
	Control_Restart OperateType = "restart" // 重启
	Control_Remove  OperateType = "remove"  // 删除（仅对服务组件有效）
	Control_Disable OperateType = "disable" // 禁用(进队服务组件有效)
	Control_Enable  OperateType = "enable"  // 启用
)

// ntp 操作类型
type NtpControlType string

const (
	Ntp_Control_Test NtpControlType = "test"
	Ntp_Control_Set  NtpControlType = "set"
)

// 服务组件信息
type ServiceInfo struct {
	ServiceType  ServiceType // 服务类型
	Ip           string      // 服务器的ip
	Port         int32       // 服务器的普通端口
	HttpPort     int32       // 服务器的http端口
	ProtobufPort int32       // 服务器的protobuf端口
	User         string      // 用户名
	Password     string      // 密码
}

func (s *ServiceInfo) String() string {
	return fmt.Sprintf("ip:%s, http port:%d, protobuf port:%d, user:%s, password:%s, type:%s", s.Ip, s.HttpPort, s.ProtobufPort, s.User, s.Password, s.ServiceType.String())
}

// 主机服务
type HostServiceInfo struct {
	ServiceInfo
	Role       []ServiceRole // 运维服务的角色（leader/worker/child/parent），注意，同一个运维服务可能会有多种角色，比如leader向上级注册的时候，它的另外一个角色就是child
	Hostname   string        // 主机名称
	ServerName string        // 服务器的名称（在界面上展示的)
}

func (o *HostServiceInfo) String() string {
	return fmt.Sprintf("%s, role:%v, hostname:%s, server name:%s", o.ServiceInfo.String(), o.Role, o.Hostname, o.ServerName)
}

type IpInfo struct {
	IP         string `json:"ip"`             // ip地址
	NetMask    string `json:"netmask"`        // 子网掩码
	GateWay    string `json:"gateway"`        // 网关
	AutoConfig bool   `json:"auto_configure"` // 自动获取
	Enable     bool   `json:"enable"`         // 是否启用
}

func (i *IpInfo) String() string {
	return fmt.Sprintf("ip:%s, net mask:%s, gate way:%s, auto config:%t, enable:%t", i.IP, i.NetMask, i.GateWay, i.AutoConfig, i.Enable)
}

// 网口信息
type NetworkInterface struct {
	MAC               string `json:"MAC"`       // mac地址
	Name              string `json:"name"`      //  网口名称
	HostName          string `json:"host_name"` // 主机名字
	IPv4              IpInfo // ipv4
	IPv6              IpInfo // ipv6
	DNS               DNS    // DNS
	Enable            bool   `json:"enable"`             // 是否插了网线
	Bandwidth         uint64 `json:"bandwidth"`          // 网口的传输带宽，单位为byte
	BandwidthReadable string `json:"bandwidth_readable"` // 阅读友好的网口的带宽
}

func (n NetworkInterface) String() string {
	return fmt.Sprintf("name:%s, MAC:%s, host name:%s, ipv4:%s, ipv6:%s, dns:%s", n.Name, n.MAC, n.HostName, n.IPv4.String(), n.IPv6.String(), n.DNS.String())
}

// 网口信息
type NetworkConfigure struct {
	Network []NetworkInterface `json:"network"`
}

func (n *NetworkConfigure) String() (result string) {
	for _, info := range n.Network {
		result += info.String()
	}
	return
}

// DNS信息
type DNS struct {
	MajorDNS string `json:"major_dns"` // 主DNS
	MinorDNS string `json:"minor_dns"` // 备用DNS
}

func (d *DNS) String() string {
	return fmt.Sprintf("major DNS:%s, minor DNS:%s", d.MajorDNS, d.MinorDNS)
}

// 存储信息
type StorageConfigure struct {
	RemoveThreshold         int64  `json:"remove_threshold"`          // 删除存储的阈值
	RemoveThresholdReadable string `json:"remove_threshold_readable"` // 删除存储的阈值
}

func (s *StorageConfigure) String() string {
	return fmt.Sprintf("remove threshold:%d, readable remove threshold:%s", s.RemoveThreshold, ReabableSize(uint64(s.RemoveThreshold)))
}

// 用户信息
type UserInfo struct {
	User     string `json:"user"`     // 用户名
	Password string `json:"password"` // 密码
}

func (u *UserInfo) String() string {
	return fmt.Sprintf("user:%s, password:%s", u.User, u.Password)
}
func (u *UserInfo) Equal(user *UserInfo) bool {
	return u.User == user.User && u.Password == user.Password
}

// 用户修改密码
type ChangePassword struct {
	User        string `json:"user"`         // 用户名
	OldPassword string `json:"old_password"` // 旧的密码
	NewPassword string `json:"new_password"` // 新的密码
}

func (c *ChangePassword) String() string {
	return fmt.Sprintf("user:%s, old password:%s, new password:%s", c.User, c.OldPassword, c.NewPassword)
}

// git信息
type GitInfo struct {
	Branch  string `json:"branch"`  // 分支
	Commit  string `json:"commit"`  // id
	Version string `json:"version"` // 版本信息
}

func (g *GitInfo) String() string {
	return fmt.Sprintf("branch:%s, commit:%s", g.Branch, g.Commit)
}

// 服务的安装信息
type ServiceModuleInfo struct {
	Path      string        `json:"path"`
	Name      string        `json:"name"`       // 名称
	Status    ServiceStatus `json:"status"`     // 状态
	Git       GitInfo       `json:"git"`        // git信息
	HasStatus bool          `json:"has_status"` // 是否有状态信息，像web, flyway是没有状态信息的
	HasLog    bool          `json:"has_log"`    // 是否有日志，像web是没有日志的
	HasUpdate bool          `json:"has_update"` // 是否可以升级，像mysql是不可以升级的
}

func (s ServiceModuleInfo) String() string {
	return fmt.Sprintf("path:%s, name:%s, status:%s, git:%s", s.Path, s.Name, s.Status, s.Git.String())
}

type ServiceModules struct {
	Service []ServiceModuleInfo `json:"service"`
}

func (s *ServiceModules) String() (info string) {

	for _, s := range s.Service {
		info += s.String() + ", "
	}
	return
}

// 实体控制
type EntityControl struct {
	EntryName string      `json:"name"`
	Control   OperateType `json:"operate"` // 控制类型
	Time      string      `json:"time"`    // 操作的事件，定时操作还是立马操作
}

func (e *EntityControl) String() string {
	return fmt.Sprintf("control:%s, time:%s", e.Control, e.Time)
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

// license信息
type LicenseInfo struct {
	DeviceCapacity int32  `json:"device_capacity"` // 支持的license数量
	Expire         string `json:"expire"`          // license过期时间
	MultiInstance  bool   `json:"multi_instance"`  // 是否支持多实例
	Valid          bool   `json:"valid"`           // license是否有效
	Error          string `json:"error"`           //错误原因
}

func (l *LicenseInfo) String() string {
	return fmt.Sprintf("valid:%t, device capacity:%d, expire:%s, multi instalce:%t, error reason:%s", l.Valid, l.DeviceCapacity, l.Expire, l.MultiInstance, l.Error)
}

// 时区信息
type TimeZone struct {
	Offset int32  `json:"offset"` // 相对于格林威治的偏移，单位是分钟
	City   string `json:"city"`   // 城市，中国则有Chongqing和Shanghai两个城市
	Name   string `json:"name"`   // 显示的名称,入GMT+08
	Code   string `json:"code"`
}

func (t *TimeZone) String() string {
	return fmt.Sprintf("city:%s, offset:%d, name:%s, code:%s", t.City, t.Offset, t.Name, t.Code)
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

//磁盘的类型
type DiskType int32

const (
	Disk_Local   DiskType = 1 // 本地磁盘
	Disk_NetWork DiskType = 2 // 网络磁盘
	Disk_Cloud   DiskType = 3 // 云存储
)

// 磁盘信息
type PartitionInfo struct {
	Name        string   `json:"name"`         // 磁盘的名字，比如disk1
	Type        DiskType `json:"type"`         // 磁盘的类型
	MountPoint  string   `json:"mount_point"`  // 挂载点
	Status      bool     `json:"status"`       // 状态
	MountStatus bool     `json:"mount_status"` //磁盘挂载状态，属性
	CapacityInfo
}

func (p *PartitionInfo) String() string {
	return fmt.Sprintf("name:%s, type:%d, mount_point:%s, status: %t, mount_status; %t, total_capacity_readable: %s, used_capacity_readable: %s", p.Name, p.Type, p.MountPoint, p.Status, p.MountStatus, p.CapacityInfo.TotalCapacityReadable, p.CapacityInfo.AvailableCapacityReadable)
}

// 磁盘操作类型
type DiksOperationType string

const (
	FormatDisk DiksOperationType = "format"		// 格式化磁盘
	MountDiks  DiksOperationType = "mount"		// 挂载磁盘
)

//需要挂载的磁盘信息
type OperationDisk struct {
	Name       string 			 `json:"name"`         // 磁盘的名字，比如disk1
	Type	   DiksOperationType `json:"type"`		  // 操作类型
}

func (p *OperationDisk) String() string {
	return  fmt.Sprintf("name: %s","disk_format: %t","disk_mount: %t", p.Name, )
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

//kafka配置信息
type KafkaInfo struct {
		Ip	string `json:"ip"`    // 需要修改的IP
		Port int   `json:"port"`  // 需要修改的端口
}

func (k *KafkaInfo) String() string {
	return fmt.Sprintf("ip:%s, port:%d", k.Ip, k.Port)
}

//kafka配置信息
type MapModeInfo struct {
	Status	int `json:"mode"`    // 地图模式
}

func (k *MapModeInfo) String() string {
	return fmt.Sprintf("ip:%s, port:%d", k.Status)
}


// 获取zone的名称
func GetTimeZoneName(offset int32) (name string) {
	off := float64(offset)
	if offset >= 0 {
		name = fmt.Sprintf("GMT+%.2f", off/3600)
	} else {
		name = fmt.Sprintf("GMT%.2f", off/3600)
	}
	return name
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

// 判断网卡是否插了网线
func IsNetworkInterfaceEnable(name string) (enable bool) {
	info, err := CommondResult(fmt.Sprintf("ifconfig %s", name))
	if err == nil {
		if strings.Index(info, "RUNNING") >= 0 {
			return true
		}
	}
	return false
}

// 获取网络接口的网络带宽
func NetworkInterfaceBandwidth(network string) (bandwidth uint64, err error) {

	// 执行命令ethtool
	output, err := CommondResult(fmt.Sprintf("ethtool %s", network))
	if err != nil {
		logrus.Errorf("get %s bandwidth fail, error %s", err.Error())
		return
	}

	// 从ethtool命令中找到网口速率字段
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.Trim(line, " ")
		line = strings.Trim(line, "\t")
		if strings.Index(line, "Speed") < 0 {
			continue
		}
		prefixIndex := strings.Index(line, ":")
		suffixIndex := strings.Index(line, "/")
		if prefixIndex >= 0 && suffixIndex >= 0 {
			speed := strings.Trim(line[prefixIndex+1:suffixIndex], " ")
			bandwidth, err = convertBandwidth(speed)
		}
		return
	}
	err = fmt.Errorf("fail")
	return
}

// 转换带宽(把字符串,比如100M改为以byte为单位的数字)
func convertBandwidth(bandwidth string) (result uint64, err error) {
	util := [...]string{"B", "K", "M", "G", "T", "P", "E", "Z", "Y"}
	for i := 0; i < len(util); i++ {
		index := strings.Index(bandwidth, util[i])
		if index >= 0 {
			base, e := strconv.ParseInt(bandwidth[:index], 10, 64)
			err = e
			if err != nil {
				return
			}

			result = uint64(base) * Pow(1024, i)
			return
		}
	}
	err = fmt.Errorf("invalid bandwidth")
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

// des进行加密
func DESEncrypt(origData, key []byte) (result string) {
	//将字节秘钥转换成block快
	block,_ := des.NewCipher(key)

	//对明文先进行补码操作
	origData = PKCS5Padding(origData,block.BlockSize())

	//设置加密方式
	blockMode := cipher.NewCBCEncrypter(block,key)

	//创建明文长度的字节数组
	crypted := make([]byte, len(origData))

	//加密明文,加密后的数据放到数组中
	blockMode.CryptBlocks(crypted,origData)

	//将字节数组转换成字符串
	result = base64.StdEncoding.EncodeToString(crypted)
	return
}

//实现明文的补码
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	//计算出需要补多少位
	padding := blockSize - len(ciphertext)%blockSize

	//Repeat()函数的功能是把参数一 切片复制 参数二count个,然后合成一个新的字节切片返回
	// 需要补padding位的padding值
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	//把补充的内容拼接到明文后面
	return append(ciphertext,padtext...)
}

//解密
func DESDecrypt(data string, key []byte) (result string) {

	// 倒叙执行一遍加密方法
	//将字符串转换成字节数组
	crypted,_ := base64.StdEncoding.DecodeString(data)

	//将字节秘钥转换成block快
	block, _ := des.NewCipher(key)

	//设置解密方式
	blockMode := cipher.NewCBCDecrypter(block,key)

	//创建密文大小的数组变量
	origData := make([]byte, len(crypted))

	//解密密文到数组origData中
	blockMode.CryptBlocks(origData,crypted)

	//去补码
	origData = PKCS5UnPadding(origData)

	result = string(origData[:])

	return
}

//去除补码
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)

	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])

	//解密去补码时需取最后一个字节，值为m，则从数据尾部删除m个字节，剩余数据即为加密前的原文
	return origData[:(length - unpadding)]
}