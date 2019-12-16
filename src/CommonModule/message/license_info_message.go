package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取license信息的消息
type LicenseInfoMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewLicenseInfoMessage(pri common.Priority, trs TransType) *LicenseInfoMessage {
	MessageId++
	return &LicenseInfoMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}}
}

func (s *LicenseInfoMessage) String() string {
	return s.BaseMessageInfo.String()
}

// 获取服务器标识的回应
type LicenseInfoResponse struct {
	BaseResponseInfo
	License common.LicenseInfo `json:"license"`
}

func NewLicenseInfoResponse(license common.LicenseInfo, msg BaseMessage) (rsp *LicenseInfoResponse) {
	rsp = &LicenseInfoResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, License: license}
	return rsp
}

func (s *LicenseInfoResponse) String() string {
	return fmt.Sprintf("license:%s, %s", s.License.String(), s.BaseResponseInfo.String())
}
