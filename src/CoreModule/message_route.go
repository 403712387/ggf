package core

import (
	"CommonModule/message"
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// 处理消息的模块
type MessageRoute struct {
	observer     []ProcessInterface // 消息的观察者
	observerLock sync.RWMutex       // 消息观察者的锁
	exit         bool               // 是否退出程序
}

// 添加观察者
func (o *MessageRoute) AddProcess(process ProcessInterface) (result bool) {
	o.observerLock.Lock()
	defer o.observerLock.Unlock()

	//  判断是否已经添加进观察者
	for _, p := range o.observer {
		if p == process {
			return false
		}
	}

	// 添加观察者
	logrus.Infof("add process %s", process.GetModuleName())
	o.observer = append(o.observer, process)
	return true
}

// 发送消息
func (o *MessageRoute) SendMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {

	//  如果消息被偷窥后返回true,则消息不会继续向下传递
	if done := o.processForeseeMessage(msg); done {
		return nil, errors.New("foresee message done")
	}

	if msg.TransType() == message.Trans_Sync { // 同步消息
		return o.processSyncMessage(msg)
	} else if msg.TransType() == message.Trans_Async { // 异步消息
		o.processAsyncMessage(msg)
	}
	return nil, nil
}

// 发送消息的回应
func (o *MessageRoute) SendResponse(rsp message.BaseResponse) {

	// 调用各个模块的OnForceseeResponse
	o.foreseeResponse(rsp)

	// 调用各个模块的OnProcessesponse
	if rsp.Message().TransType() == message.Trans_Sync { // 同步消息
		o.processSyncResponse(rsp)
	} else if rsp.Message().TransType() == message.Trans_Async { // 异步消息
		o.processAsyncResponse(rsp)
	}
}

//  偷窥消息
func (o *MessageRoute) processForeseeMessage(msg message.BaseMessage) bool {
	o.observerLock.RLock()
	defer o.observerLock.RUnlock()

	// 把消息转发给各个模块的OnForeseeMessage
	done := false
	for _, p := range o.observer {
		if ret := p.OnForeseeMessage(msg); ret {
			done = ret
		}
	}
	return done
}

// 同步消息
func (o *MessageRoute) processSyncMessage(msg message.BaseMessage) (rsp message.BaseResponse, err error) {
	o.observerLock.RLock()
	defer o.observerLock.RUnlock()

	// 调用各个模块的消息处理函数
	for _, p := range o.observer {
		if rsp, err := p.OnProcessMessage(msg); rsp != nil {
			// 向各个模块发送消息的回应
			o.SendResponse(rsp)
			return rsp, err
		}
	}

	return nil, errors.New("not have module process message")
}

// 异步消息
func (o *MessageRoute) processAsyncMessage(message message.BaseMessage) {
	o.observerLock.RLock()
	defer o.observerLock.RUnlock()

	// 把消息放入队列中
	for _, p := range o.observer {
		p.PushMessage(message)
	}
}

// 偷窥回应
func (o *MessageRoute) foreseeResponse(response message.BaseResponse) {
	o.observerLock.RLock()
	defer o.observerLock.RUnlock()

	for _, p := range o.observer {
		p.OnForeseeResponse(response)
	}
}

// 同步消息的回应
func (o *MessageRoute) processSyncResponse(response message.BaseResponse) {
	o.observerLock.RLock()
	defer o.observerLock.RUnlock()

	for _, p := range o.observer {
		p.OnProcessResponse(response)
	}
}

// 异步消息的回应
func (o *MessageRoute) processAsyncResponse(response message.BaseResponse) {
	o.observerLock.RLock()
	defer o.observerLock.RUnlock()

	for _, p := range o.observer {
		p.PushResponse(response)
	}
}

// 启动 message route
func (o *MessageRoute) Beginwork() {
	o.observerLock.RLock()
	defer o.observerLock.RUnlock()

	// 开启消息处理队列
	for _, p := range o.observer {
		p.BeginProcessLoop()
	}

	// 休眠一秒,以免调用p.Init的时候，模块的loop还没启动
	time.Sleep(1 * time.Second)

	// 初始化
	for _, p := range o.observer {
		p.Init()
	}

	// 休眠一秒
	time.Sleep(1 * time.Second)

	// 开始工作
	for _, p := range o.observer {
		p.BeginWork()
	}

	// 开始循环
	for !o.exit {
		time.Sleep(time.Minute)
	}
}

// 停止 message route
func (o *MessageRoute) StopWork() {
	o.observerLock.RLock()
	defer o.observerLock.Unlock()

	// 停止消息处理队列
	for _, p := range o.observer {
		p.StopProcessLoop()
	}

	// 停止工作
	for _, p := range o.observer {
		p.StopWork()
	}

	// 反初始化
	for _, p := range o.observer {
		p.Uninit()
	}
}
