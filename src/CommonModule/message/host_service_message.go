package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取host服务的信息，包括开机时间，git信息
type GgfServiceMessage struct {
	BaseMessageInfo
}

// 生成消息
func NewGgfServiceMessage(pri common.Priority, tra TransType) (msg *GgfServiceMessage) {
	MessageId++
	return &GgfServiceMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: tra, birthday: time.Now()}}
}

func (h *GgfServiceMessage) String() string {
	return h.BaseMessageInfo.String()
}

// 获取服务器标识的回应
type HostServiceResponse struct {
	BaseResponseInfo
	HostStartup   time.Time // 服务启动时间
	SystemStartup time.Time // 系统开机时间
	GitBranch     string    // git分支信息
	GitCommitID   string    // git commit id信息
}

func NewGgfServiceResponse(host, system time.Time, gitBranch, gitCommit string, msg BaseMessage) (rsp *HostServiceResponse) {
	rsp = &HostServiceResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, HostStartup: host, SystemStartup: system, GitBranch: gitBranch, GitCommitID: gitCommit}
	return rsp
}

func (h *HostServiceResponse) String() string {
	return fmt.Sprintf("host start up time:%s, system start up time:%s, git branch:%s, git commit id:%s, %s", h.HostStartup.Format("2006-01-02 15:04:05"), h.SystemStartup.Format("2006-01-02 15:04:05"), h.GitBranch, h.GitCommitID, h.BaseResponseInfo.String())
}
