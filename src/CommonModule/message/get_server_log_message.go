package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取系统日志的消息
type GetServerLogMessage struct {
	BaseMessageInfo
	Begin int64 // 开始行数
	End   int64 // 结束行数
}

// 生成消息
func NewGetServerLogMessage(begin, end int64, pri common.Priority, trs TransType) *GetServerLogMessage {
	MessageId++
	return &GetServerLogMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}, Begin: begin, End: end}
}

func (d *GetServerLogMessage) String() string {
	return fmt.Sprintf("begin:%d, end:%d, %s", d.Begin, d.End, d.BaseMessageInfo.String())
}

// 获取系统日志的回应
type GetServerLogResponse struct {
	BaseResponseInfo
	Begin int64    // 开始行数
	End   int64    // 结束行数
	Log   []string // 日志内容
}

func NewGetServerLogResponse(begin, end int64, log []string, msg BaseMessage) (rsp *GetServerLogResponse) {
	rsp = &GetServerLogResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Begin: begin, End: end, Log: log}
	return rsp
}

func (d *GetServerLogResponse) String() string {
	return fmt.Sprintf("begin:%d, end:%d, %s", d.Begin, d.End, d.BaseResponseInfo.String())
}
