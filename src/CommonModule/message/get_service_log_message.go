package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 下载系统日志的消息
type GetServiceLogMessage struct {
	BaseMessageInfo
	Service string // 服务名称
	Begin   int64  // 开始行数
	End     int64  // 结束行数
}

// 生成消息
func NewGetServiceLogMessage(service string, begin, end int64, pri common.Priority, trs TransType) *GetServiceLogMessage {
	MessageId++
	return &GetServiceLogMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}, Service: service, Begin: begin, End: end}
}

func (d *GetServiceLogMessage) String() string {
	return fmt.Sprintf("begin:%d, end:%d, %s", d.Begin, d.End, d.BaseMessageInfo.String())
}

// 下载系统日志的回应
type GetServiceLogResponse struct {
	BaseResponseInfo
	Begin int64    // 开始行数
	End   int64    // 结束行数
	Log   []string // 日志内容
}

func NewGetServiceLogResponse(begin, end int64, log []string, msg BaseMessage) (rsp *GetServiceLogResponse) {
	rsp = &GetServiceLogResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Begin: begin, End: end, Log: log}
	return rsp
}

func (d *GetServiceLogResponse) String() string {
	return fmt.Sprintf("%s", d.BaseResponseInfo.String())
}
