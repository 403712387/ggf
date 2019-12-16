package message

import (
	"CommonModule"
	"fmt"
	"time"
)

// 获取存储配置
type GetStorageConfigureMessage struct {
	BaseMessageInfo
}

// 生成新的消息
func NewGetStorageConfigureMessage(pri common.Priority, trs TransType) *GetStorageConfigureMessage {
	MessageId++
	return &GetStorageConfigureMessage{BaseMessageInfo: BaseMessageInfo{id: MessageId, priority: pri, trans: trs, birthday: time.Now()}}
}

func (g *GetStorageConfigureMessage) String() string {
	return fmt.Sprintf("%s", g.BaseMessageInfo.String())
}

// 获取更新配置的回应
type GetStorageConfigureResponse struct {
	BaseResponseInfo
	StorageConfigure common.StorageConfigure // 存储的配置
	LastRemove       []string                // 最后删除的存储目录
}

func NewGetStorageConfigureResponse(conf common.StorageConfigure, remove []string, msg BaseMessage) (rsp *GetStorageConfigureResponse) {
	rsp = &GetStorageConfigureResponse{BaseResponseInfo: BaseResponseInfo{message: msg}, StorageConfigure: conf, LastRemove: remove}
	return rsp
}

func (g *GetStorageConfigureResponse) String() string {
	return fmt.Sprintf("storage:%s, %s", g.StorageConfigure.String(), g.BaseResponseInfo.String())
}
