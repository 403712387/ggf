package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取CPU使用率
type GetCpuStatisticMessage struct {
	BaseMessageInfo
	BeginTime string // 开始时间
	EndTime   string // 结束时间
}

// 生成消息
func NewGetCpuStatisticMessage(begin, end string, pri common.Priority, trs TransType) *GetCpuStatisticMessage {
	MessageId++
	return &GetCpuStatisticMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}, BeginTime: begin, EndTime: end}
}

func (d *GetCpuStatisticMessage) String() string {
	return fmt.Sprintf("begin:%s, end:%s, %s", d.BeginTime, d.EndTime, d.BaseMessageInfo.String())
}

// 获取CPU使用率的回应
type GetCpuStatisticResponse struct {
	BaseResponseInfo
	Time []time.Time
	Cpu  []common.CapacityInfo
}

func NewGetCpuStatisticResponse(statisticTime []time.Time, cpu []common.CapacityInfo, msg BaseMessage) (rsp *GetCpuStatisticResponse) {
	rsp = &GetCpuStatisticResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Time: statisticTime, Cpu: cpu}
	return rsp
}

func (d *GetCpuStatisticResponse) String() (result string) {
	for index, time := range d.Time {
		result += fmt.Sprintf("time:%s, cpu:%s", time.Format("2006-01-02 15:04:05"), d.Cpu[index].String())
	}
	result += d.BaseResponseInfo.String()
	return
}
