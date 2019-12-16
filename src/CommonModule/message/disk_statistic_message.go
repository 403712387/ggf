package message

import (
	"CommonModule"
	"time"
)

// 磁盘的读写速率
type DiskStatisticMessage struct {
	BaseMessageInfo
	time.Time
	Statistic []common.DiskStatisticInfo
}

// 生成消息
func NewDiskStatisticMessage(samplingTime time.Time, stat []common.DiskStatisticInfo, pri common.Priority, tra TransType) (msg *DiskStatisticMessage) {
	MessageId++
	return &DiskStatisticMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, Statistic: stat, Time: samplingTime}
}

func (d *DiskStatisticMessage) String() (result string) {
	for _, stat := range d.Statistic {
		result += stat.String() + ", "
	}

	result += d.BaseMessageInfo.String()
	return
}
