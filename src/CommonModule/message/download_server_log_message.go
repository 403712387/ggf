package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 下载系统日志的消息
type DownloadServerLogMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewDownloadServerLogMessage(pri common.Priority, trs TransType) *DownloadServerLogMessage {
	MessageId++
	return &DownloadServerLogMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}}
}

func (d *DownloadServerLogMessage) String() string {
	return fmt.Sprintf(" %s", d.BaseMessageInfo.String())
}

// 下载系统日志的回应
type DownloadServerLogResponse struct {
	BaseResponseInfo
	LogPath    string // 下载路径
	ZipLogPath string // 日志 压缩的下载路径
}

func NewDownloadServerLogResponse(log, zipLog string, msg BaseMessage) (rsp *DownloadServerLogResponse) {
	rsp = &DownloadServerLogResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, LogPath: log, ZipLogPath: zipLog}
	return rsp
}

func (d *DownloadServerLogResponse) String() string {
	return fmt.Sprintf("log path:%s, zip log path:%s, %s", d.LogPath, d.ZipLogPath, d.BaseResponseInfo.String())
}
