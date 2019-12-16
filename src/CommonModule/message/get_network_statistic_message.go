package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取网络使用
type GetNetworkStatisticMessage struct {
	BaseMessageInfo
	BeginTime string // 开始时间
	EndTime   string // 结束时间
}

// 生成消息
func NewGetNetworkStatisticMessage(begin, end string, pri common.Priority, trs TransType) *GetNetworkStatisticMessage {
	MessageId++
	return &GetNetworkStatisticMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}, BeginTime: begin, EndTime: end}
}

func (d *GetNetworkStatisticMessage) String() string {
	return fmt.Sprintf("begin:%s, end:%s, %s", d.BeginTime, d.EndTime, d.BaseMessageInfo.String())
}

// 获取网络使用的回应
type GetNetworkStatisticResponse struct {
	BaseResponseInfo
	Time      []time.Time
	Statistic [][]common.NetworkStatisticInfo
}

func NewGetNetworkStatisticResponse(statisticTime []time.Time, stat [][]common.NetworkStatisticInfo, msg BaseMessage) (rsp *GetNetworkStatisticResponse) {
	rsp = &GetNetworkStatisticResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Time: statisticTime, Statistic: stat}
	return rsp
}

func (d *GetNetworkStatisticResponse) String() (result string) {
	for index, time := range d.Time {
		result += fmt.Sprintf("time:%s, network:%v", time.Format("2006-01-02 15:04:05"), d.Statistic[index])
	}
	result += d.BaseResponseInfo.String()
	return
}
