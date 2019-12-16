package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// CPU使用情况
type CPUStatisticMessage struct {
	BaseMessageInfo
	time.Time
	common.CapacityInfo
}

// 生成消息
func NewCPUStatisticMessage(samplingTime time.Time, capacity common.CapacityInfo, pri common.Priority, tra TransType) (msg *CPUStatisticMessage) {
	MessageId++
	return &CPUStatisticMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, CapacityInfo: capacity, Time: samplingTime}
}

func (c *CPUStatisticMessage) String() (result string) {
	return fmt.Sprintf("cpu:%s, %s", c.CapacityInfo.String(), c.BaseMessageInfo.String())
}
