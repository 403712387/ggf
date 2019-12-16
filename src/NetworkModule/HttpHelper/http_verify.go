package HttpHelper

import (
	"CommonModule"
	"fmt"
	"regexp"
)

// 校验IP设置
func VerifyNetworkConfigure(conf common.NetworkInterface) (err error) {

	//校验IPv4
	if conf.IPv4.Enable && !conf.IPv4.AutoConfig {

		// 校验ip
		if !IsValidIPv4(conf.IPv4.IP) {
			err = fmt.Errorf("invalid Ipv4: %s", conf.IPv4.IP)
			return
		}

		// 校验网关
		if !IsValidIPv4(conf.IPv4.GateWay) {
			err = fmt.Errorf("invalid GateWay: %s", conf.IPv4.GateWay)
			return
		}

		// 校验子网掩码
		if !IsValidIPv4(conf.IPv4.NetMask) {
			err = fmt.Errorf("invalid NetMask: %s", conf.IPv4.NetMask)
			return
		}
	}

	// 校验IPv6
	if conf.IPv6.Enable && !conf.IPv6.AutoConfig {

		// 校验ip
		if !IsValidIPv6(conf.IPv6.IP) {
			err = fmt.Errorf("invalid Ipv6: %s", conf.IPv6.IP)
			return
		}

		// 校验网关
		if !IsValidIPv6(conf.IPv6.GateWay) {
			err = fmt.Errorf("invalid GateWay: %s", conf.IPv6.GateWay)
			return
		}
	}

	// 校验DNS
	if conf.DNS.MajorDNS != "" {
		if !IsValidIPv4(conf.DNS.MajorDNS) {
			err = fmt.Errorf("invalid DNS: %s", conf.DNS.MajorDNS)
			return
		}

	}
	if conf.DNS.MinorDNS != "" {
		if !IsValidIPv4(conf.DNS.MinorDNS) {
			err = fmt.Errorf("invalid DNS: %s", conf.DNS.MinorDNS)
			return
		}
	}
	return
}

func IsValidIPv4(ip string) (b bool) {
	if m, _ := regexp.MatchString(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`, ip); !m {
		return false
	}
	return true
}

func IsValidIPv6(ip string) (b bool) {
	if m, err := regexp.MatchString(`^([\da-fA-F]{1,4}:){7}[\da-fA-F]{1,4}|:((:[\da−fA−F]1,4)1,6|:)|:((:[\da−fA−F]1,4)1,6|:)|^[\da-fA-F]{1,4}:((:[\da-fA-F]{1,4}){1,5}|:)|([\da−fA−F]1,4:)2((:[\da−fA−F]1,4)1,4|:)|([\da−fA−F]1,4:)2((:[\da−fA−F]1,4)1,4|:)|^([\da-fA-F]{1,4}:){3}((:[\da-fA-F]{1,4}){1,3}|:)|([\da−fA−F]1,4:)4((:[\da−fA−F]1,4)1,2|:)|([\da−fA−F]1,4:)4((:[\da−fA−F]1,4)1,2|:)|^([\da-fA-F]{1,4}:){5}:([\da-fA-F]{1,4})?|([\da−fA−F]1,4:)6:|([\da−fA−F]1,4:)6:^([\da-fA-F]{1,4}:){7}[\da-fA-F]{1,4}|:((:[\da−fA−F]1,4)1,6|:)|:((:[\da−fA−F]1,4)1,6|:)|^[\da-fA-F]{1,4}:((:[\da-fA-F]{1,4}){1,5}|:)|([\da−fA−F]1,4:)2((:[\da−fA−F]1,4)1,4|:)|([\da−fA−F]1,4:)2((:[\da−fA−F]1,4)1,4|:)|^([\da-fA-F]{1,4}:){3}((:[\da-fA-F]{1,4}){1,3}|:)|([\da−fA−F]1,4:)4((:[\da−fA−F]1,4)1,2|:)|([\da−fA−F]1,4:)4((:[\da−fA−F]1,4)1,2|:)|^([\da-fA-F]{1,4}:){5}:([\da-fA-F]{1,4})?|([\da−fA−F]1,4:)6:|([\da−fA−F]1,4:)6:`, ip); !m {
		return false
	} else if err != nil {
		return false
	}
	return true
}
