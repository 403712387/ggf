package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取硬盘使用率
type GetDiskStatisticMessage struct {
	BaseMessageInfo
	BeginTime string // 开始时间
	EndTime   string // 结束时间
}

// 生成消息
func NewGetDiskStatisticMessage(begin, end string, pri common.Priority, trs TransType) *GetDiskStatisticMessage {
	MessageId++
	return &GetDiskStatisticMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}, BeginTime: begin, EndTime: end}
}

func (d *GetDiskStatisticMessage) String() string {
	return fmt.Sprintf("begin:%s, end:%s, %s", d.BeginTime, d.EndTime, d.BaseMessageInfo.String())
}

// 获取硬盘使用率的回应
type GetDiskStatisticResponse struct {
	BaseResponseInfo
	Time      []time.Time
	Statistic [][]common.DiskStatisticInfo
}

func NewGetDiskStatisticResponse(statisticTime []time.Time, stat [][]common.DiskStatisticInfo, msg BaseMessage) (rsp *GetDiskStatisticResponse) {
	rsp = &GetDiskStatisticResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Time: statisticTime, Statistic: stat}
	return rsp
}

func (d *GetDiskStatisticResponse) String() (result string) {
	for index, time := range d.Time {
		result += fmt.Sprintf("time:%s, disk:%v", time.Format("2006-01-02 15:04:05"), d.Statistic[index])
	}
	result += d.BaseResponseInfo.String()
	return
}
