package core

import (
	"CommonModule"
	"CommonModule/message"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

type MessageQueue []message.BaseMessage
type ResponseQueue []message.BaseResponse

// 消息队列
type MessageList struct {
	message          [common.Priority_Count]MessageQueue // 存放消息的数组
	messageCondition sync.Cond                           // 消息队列的通知

	response          [common.Priority_Count]ResponseQueue // 存放回应的的数组
	responseCondition sync.Cond                            // 回应队列的通知

	exit         bool             // 是否退出
	MessageRoute *MessageRoute    // 消息路由模块
	Process      ProcessInterface // 当前模块
	ModuleName   string           // 模块名称

	emptyErr error // 当队列为空的时候的错误
}

// 初始化
func (m *MessageList) Init() {
	m.exit = false
	m.emptyErr = fmt.Errorf("empty message list")
}

// 反初始化
func (m *MessageList) Uninit() {

}

// 开始工作
func (m *MessageList) BeginWork() {

}

// 停止工作
func (m *MessageList) StopWork() {

}

// 获取模块名称
func (m *MessageList) GetModuleName() string {
	return m.ModuleName
}

// 开启消息队列
func (m *MessageList) BeginProcessLoop() {
	// 初始化变量
	m.messageCondition = sync.Cond{L: &sync.RWMutex{}}
	m.responseCondition = sync.Cond{L: &sync.RWMutex{}}

	// 处理消息的循环
	go m.processMessageLoop()

	// 处理回应的循环
	go m.processResponseLoop()
}

// 退出消息队列
func (m *MessageList) StopProcessLoop() {
	m.exit = true

	// 唤醒
	m.messageCondition.Broadcast()
	m.responseCondition.Broadcast()
}

// 把消息放入消息队列
func (m *MessageList) PushMessage(msg message.BaseMessage) {
	m.messageCondition.L.Lock()
	defer m.messageCondition.L.Unlock()

	priority := msg.MessagePriority()
	m.message[priority] = append(m.message[priority], msg)
	m.messageCondition.Signal()
}

// 把消息的回应放入回应队列
func (m *MessageList) PushResponse(rsp message.BaseResponse) {
	if rsp == nil {
		return
	}

	m.responseCondition.L.Lock()
	defer m.responseCondition.L.Unlock()

	priority := rsp.Message().MessagePriority()
	m.response[priority] = append(m.response[priority], rsp)
	m.responseCondition.Signal()
}

// 偷窥消息
func (m *MessageList) OnForeseeMessage(msg message.BaseMessage) (done bool) {
	logrus.Errorf("never should foresee message %s", msg.String())
	return
}

// 处理消息
func (m *MessageList) OnProcessMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	logrus.Errorf("never should process message %s", msg.String())
	return nil, nil
}

// 偷窥消息的回应
func (m *MessageList) OnForeseeResponse(rsp message.BaseResponse) (done bool) {
	logrus.Errorf("never should foresee response %s", rsp.String())
	return
}

// 处理消息的回应
func (m *MessageList) OnProcessResponse(rsp message.BaseResponse) {
	logrus.Errorf("never should process response %s", rsp.String())
	return
}

// 是否退出
func (m *MessageList) IsExit() bool {
	return m.exit
}

// 发送消息
func (m *MessageList) SendMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	return m.MessageRoute.SendMessage(msg)
}

// 消息队列是否为空
func (m *MessageList) isMessageEmpty() bool {
	m.messageCondition.L.Lock()
	defer m.messageCondition.L.Unlock()

	for i := 0; i < len(m.message); i++ {
		if len(m.message[i]) > 0 {
			return false
		}
	}
	return true
}

// 回应队列是否为空
func (m *MessageList) isResponseEmpty() bool {
	m.responseCondition.L.Lock()
	defer m.responseCondition.L.Unlock()

	for i := 0; i < len(m.response); i++ {
		if len(m.response[i]) > 0 {
			return false
		}
	}
	return true
}

// 取出一个消息
func (m *MessageList) popMessage() (msg message.BaseMessage, err error) {

	// 如果队列中没有消息，则等待
	if m.isMessageEmpty() {

		// 等待通知
		m.messageCondition.L.Lock()
		m.messageCondition.Wait()

	} else {
		m.messageCondition.L.Lock()
	}
	defer m.messageCondition.L.Unlock()

	// 从队列中取出消息
	for priority := 0; priority < int(common.Priority_Count); priority++ {
		if len(m.message[priority]) > 0 {
			msg := m.message[priority][0]
			m.message[priority] = m.message[priority][1:]
			return msg, nil
		}
	}

	return nil, m.emptyErr
}

// 取出一个消息的回应
func (m *MessageList) popResponse() (rsp message.BaseResponse, err error) {
	// 如果队列中没有回应，则等待
	if m.isResponseEmpty() {
		m.responseCondition.L.Lock()
		m.responseCondition.Wait()
		defer m.responseCondition.L.Unlock()

		// 取出回应
		for priority := 0; priority < int(common.Priority_Count); priority++ {
			if len(m.response[priority]) > 0 {
				rsp := m.response[priority][0]
				m.response[priority] = m.response[priority][1:]
				return rsp, nil
			}
		}
	}

	return nil, errors.New("empty response list")
}

// 处理消息的goroutine
func (m *MessageList) processMessageLoop() {
	logrus.Infof("begin process message loop in %s", m.ModuleName)
	for !m.exit {

		// 从队列中取出消息
		msg, err := m.popMessage()
		if err != nil || msg == nil {
			continue
		}

		// 处理消息(这里的消息都是异步消息)
		rsp, err := m.Process.OnProcessMessage(msg)
		if err == nil && rsp != nil {
			m.PushResponse(rsp)
		}
	}
	logrus.Infof("end process message loop in %s", m.ModuleName)
}

// 处理消息回应的goroutine
func (m *MessageList) processResponseLoop() {
	logrus.Infof("begin process response loop in %s", m.ModuleName)
	for !m.exit {

		// 从队列中取出回应
		rsp, err := m.popResponse()
		if err != nil {
			continue
		}

		// 处理回应
		m.Process.OnProcessResponse(rsp)
	}
	logrus.Infof("end process response loop in %s", m.ModuleName)
}
