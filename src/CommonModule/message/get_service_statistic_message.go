package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取网络使用
type GetServiceStatisticMessage struct {
	BaseMessageInfo
	BeginTime string // 开始时间
	EndTime   string // 结束时间
}

// 生成消息
func NewGetServiceStatisticMessage(begin, end string, pri common.Priority, trs TransType) *GetServiceStatisticMessage {
	MessageId++
	return &GetServiceStatisticMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}, BeginTime: begin, EndTime: end}
}

func (d *GetServiceStatisticMessage) String() string {
	return fmt.Sprintf("begin:%s, end:%s, %s", d.BeginTime, d.EndTime, d.BaseMessageInfo.String())
}

// 获取网络使用的回应
type GetServiceStatisticResponse struct {
	BaseResponseInfo
	Time      []time.Time
	Statistic [][]common.ProcessInfo
}

func NewGetServiceStatisticResponse(statisticTime []time.Time, stat [][]common.ProcessInfo, msg BaseMessage) (rsp *GetServiceStatisticResponse) {
	rsp = &GetServiceStatisticResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Time: statisticTime, Statistic: stat}
	return rsp
}

func (d *GetServiceStatisticResponse) String() (result string) {
	for index, time := range d.Time {
		result += fmt.Sprintf("time:%s, memory:%v", time.Format("2006-01-02 15:04:05"), d.Statistic[index])
	}
	result += d.BaseResponseInfo.String()
	return
}
