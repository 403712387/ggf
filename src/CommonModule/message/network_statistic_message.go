package message

import (
	"CommonModule"
	"time"
)

// 网络使用情况
type NetworkStatisticMessage struct {
	BaseMessageInfo
	time.Time
	Statistic []common.NetworkStatisticInfo
}

// 生成消息
func NewNetworkStatisticMessage(samplingTime time.Time, stat []common.NetworkStatisticInfo, pri common.Priority, tra TransType) (msg *NetworkStatisticMessage) {
	MessageId++
	return &NetworkStatisticMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, Statistic: stat, Time: samplingTime}
}

func (n *NetworkStatisticMessage) String() (result string) {
	for _, network := range n.Statistic {
		result += network.String() + ", "
	}
	result += n.BaseMessageInfo.String()
	return
}
