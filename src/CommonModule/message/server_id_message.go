package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取服务器标识符的消息
type ServerIdMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewServerIdMessage(pri common.Priority, trs TransType) *ServerIdMessage {
	MessageId++
	return &ServerIdMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}}
}

func (s *ServerIdMessage) String() string {
	return s.BaseMessageInfo.String()
}

// 获取服务器标识的回应
type ServerIdResponse struct {
	BaseResponseInfo
	ID   		string // 服务器标识
	Path 		string // 下载的路径
	ZipPath 	string // 压缩文件
}

func NewServerIdResponse(id, path, zip string, msg BaseMessage) (rsp *ServerIdResponse) {
	rsp = &ServerIdResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, ID: id, Path: path, ZipPath:zip}
	return rsp
}

func (s *ServerIdResponse) String() string {
	return fmt.Sprintf("id:%s, path:%s, %s", s.ID, s.Path, s.BaseResponseInfo.String())
}
