package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取内存使用率
type GetMemoryStatisticMessage struct {
	BaseMessageInfo
	BeginTime string // 开始时间
	EndTime   string // 结束时间
}

// 生成消息
func NewGetMemoryStatisticMessage(begin, end string, pri common.Priority, trs TransType) *GetMemoryStatisticMessage {
	MessageId++
	return &GetMemoryStatisticMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}, BeginTime: begin, EndTime: end}
}

func (d *GetMemoryStatisticMessage) String() string {
	return fmt.Sprintf("begin:%s, end:%s, %s", d.BeginTime, d.EndTime, d.BaseMessageInfo.String())
}

// 获取内存使用率的回应
type GetMemoryStatisticResponse struct {
	BaseResponseInfo
	Time      []time.Time
	Statistic []common.MemoryStatisticInfo
}

func NewGetMemoryStatisticResponse(statisticTime []time.Time, stat []common.MemoryStatisticInfo, msg BaseMessage) (rsp *GetMemoryStatisticResponse) {
	rsp = &GetMemoryStatisticResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Time: statisticTime, Statistic: stat}
	return rsp
}

func (d *GetMemoryStatisticResponse) String() (result string) {
	for index, time := range d.Time {
		result += fmt.Sprintf("time:%s, memory:%s", time.Format("2006-01-02 15:04:05"), d.Statistic[index].String())
	}
	result += d.BaseResponseInfo.String()
	return
}
