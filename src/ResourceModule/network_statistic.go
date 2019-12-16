package ResourceModule

import (
	"CommonModule"
	"github.com/shirou/gopsutil/net"
	"time"
)

// 统计网卡的使用情况
func statisticNetwork() (info []common.NetworkStatisticInfo, err error) {
	firstStat, err := net.IOCounters(true)
	time.Sleep(time.Second)
	secondStat, err := net.IOCounters(true)

	//  根据两次统计间隔，获取网卡的速率
	for _, first := range firstStat {
		if first.Name == "lo" {
			continue
		}
		for _, second := range secondStat {
			if first.Name != second.Name {
				continue
			}
			capacity, e := common.NetworkInterfaceBandwidth(first.Name)
			if e != nil {
				continue
			}

			// 计算网口速率信息
			var network common.NetworkStatisticInfo
			network.Name = first.Name
			network.CapacityByte = capacity
			network.CapacityByteReadable = common.ReabableSize(network.CapacityByte)
			network.ReceiveByte = second.BytesRecv - first.BytesRecv
			network.ReceiveByteReadable = common.ReabableSize(network.ReceiveByte)
			network.SendByte = second.BytesSent - first.BytesSent
			network.SendByteReadable = common.ReabableSize(network.SendByte)
			network.SendPacket = second.PacketsSent - first.PacketsSent
			network.ReceivePacket = second.PacketsRecv - first.PacketsRecv
			network.SendError = second.Errout - first.Errout
			network.ReceiveError = second.Errin - first.Errin
			network.SendDrop = second.Dropout - first.Dropout
			network.ReceiveDrop = second.Dropin - first.Dropin

			info = append(info, network)
			break
		}
	}
	return
}

