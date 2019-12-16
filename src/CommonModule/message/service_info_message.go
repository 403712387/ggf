package message

import (
	"CommonModule"
	"time"
)

// 获取服务组件信息的消息
type ServiceInfoMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewServiceInfoMessage(pri common.Priority, tra TransType) (msg *ServiceInfoMessage) {
	MessageId++
	return &ServiceInfoMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}}
}

func (s *ServiceInfoMessage) String() string {
	return s.BaseMessageInfo.String()
}

type ServiceInfoResponse struct {
	BaseResponseInfo
	Info common.ServiceModules // 服务组件的信息
}

// 生成回应
func NewServiceInfoResponse(info common.ServiceModules, msg BaseMessage) (response *ServiceInfoResponse) {
	resp := ServiceInfoResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, Info: info}
	return &resp
}

func (s *ServiceInfoResponse) String() (result string) {

	result = s.Info.String() + " ," + s.BaseResponseInfo.String()
	return
}
