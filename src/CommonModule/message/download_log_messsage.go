package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取服务组件信息的消息
type DownloadServiceLogMessage struct {
	BaseMessageInfo
	Service string //   服务名称
}

// 生成消息
func NewDownloadServiceLogMessage(name string, pri common.Priority, tra TransType) (msg *DownloadServiceLogMessage) {
	MessageId++
	return &DownloadServiceLogMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}, Service: name}
}

func (s *DownloadServiceLogMessage) String() string {
	return fmt.Sprintf("service:%s, %s", s.Service, s.BaseMessageInfo.String())
}

type DownloadServiceLogResponse struct {
	BaseResponseInfo
	Log string // 日志的下载路径
}

// 生成回应
func NewDownloadServiceLogResponse(log string, msg BaseMessage) (response *DownloadServiceLogResponse) {
	resp := DownloadServiceLogResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Log: log}
	return &resp
}

func (s *DownloadServiceLogResponse) String() (result string) {
	result = fmt.Sprintf("log:%s, %s", s.Log, s.BaseResponseInfo.String())
	return
}
