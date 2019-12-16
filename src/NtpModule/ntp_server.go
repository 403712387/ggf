package NtpModule

import (
	"github.com/btfak/sntp/netapp"
	"github.com/btfak/sntp/netevent"
)

// 启动ntp服务
func StartNtpServer(port int) {
	var handler = netapp.GetHandler()
	netevent.Reactor.ListenUdp(port, handler)
	netevent.Reactor.Run()
}
