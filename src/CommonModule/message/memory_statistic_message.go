package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 内存使用情况
type MemoryStatisticMessage struct {
	BaseMessageInfo
	time.Time
	common.MemoryStatisticInfo
}

// 生成消息
func NewMemoryInfoMessage(samplingTime time.Time, capacity common.MemoryStatisticInfo, pri common.Priority, tra TransType) (msg *MemoryStatisticMessage) {
	MessageId++
	return &MemoryStatisticMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, MemoryStatisticInfo: capacity, Time: samplingTime}
}

func (m *MemoryStatisticMessage) String() (result string) {
	return fmt.Sprintf("memory:%s, %s", m.MemoryStatisticInfo.String(), m.BaseMessageInfo.String())
}
