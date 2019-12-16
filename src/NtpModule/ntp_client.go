package NtpModule

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

/*
时间了作为ntp客户端，向ntp服务端获取时间
*/
type packet struct {
	Settings       uint8  // leap yr indicator, ver number, and mode
	Stratum        uint8  // stratum of local clock
	Poll           int8   // poll exponent
	Precision      int8   // precision exponent
	RootDelay      uint32 // root delay
	RootDispersion uint32 // root dispersion
	ReferenceID    uint32 // reference id
	RefTimeSec     uint32 // reference timestamp sec
	RefTimeFrac    uint32 // reference timestamp fractional
	OrigTimeSec    uint32 // origin time secs
	OrigTimeFrac   uint32 // origin time fractional
	RxTimeSec      uint32 // receive time secs
	RxTimeFrac     uint32 // receive time frac
	TxTimeSec      uint32 // transmit time secs
	TxTimeFrac     uint32 // transmit time frac
}

// 从ntp服务器获取时间
func ntpTime(ip string, port int32) (now time.Time, err error) {

	// 连接ntp服务器
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", ip, port), 3*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()
	if err = conn.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
		return
	}

	// 向ntp服务器发送请求
	req := &packet{Settings: 0x1B}
	if err = binary.Write(conn, binary.BigEndian, req); err != nil {
		return
	}

	// 读取ntp服务器发送过来的回应
	rsp := &packet{}
	if err = binary.Read(conn, binary.BigEndian, rsp); err != nil {
		return
	}

	// 解析接收到的时间信息
	const ntpEpochOffset = 2208988800
	secs := float64(rsp.TxTimeSec) - ntpEpochOffset
	nano := (int64(rsp.TxTimeFrac) * 1e9) >> 32
	now = time.Unix(int64(secs), nano)
	return
}
