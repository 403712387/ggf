package message

import (
	"CommonModule"
	"time"
)

// 统计进程的资源使用
type ServiceStatisticMessage struct {
	BaseMessageInfo
	time.Time
	Statistic []common.ProcessInfo
}

// 生成消息
func NewServiceStatisticMessage(samplintTime time.Time, stat []common.ProcessInfo, pri common.Priority, tra TransType) (msg *ServiceStatisticMessage) {
	MessageId++
	return &ServiceStatisticMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, Statistic: stat, Time: samplintTime}
}

func (d *ServiceStatisticMessage) String() (result string) {
	for _, stat := range d.Statistic {
		result += stat.String() + ", "
	}

	result += d.BaseMessageInfo.String()
	return
}
