package ipnet

import (
	"net"
)

func NextStart(ipNet *net.IPNet) net.IP {
	ip := LastIP(ipNet)
	return NextIP(ip)
}

// MaskLess CIDR大小比较
// IPv4 < IPv6
// ones越大，子网范围越小
func MaskLess(a, b *net.IPNet) bool {
	if lenA, lenB := len(a.IP), len(b.IP); lenA != lenB {
		return lenA < lenB
	}
	aOnes, _ := a.Mask.Size()
	bOnes, _ := b.Mask.Size()
	return aOnes > bOnes
}

// LastIP IPNet最后一个IP
func LastIP(ipNet *net.IPNet) net.IP {
	ip, mask := ipNet.IP, ipNet.Mask
	ipLen := len(ip)
	res := make(net.IP, ipLen)
	for i := 0; i < ipLen; i++ {
		res[i] = ip[i] | ^mask[i]
	}
	return res
}

// Contains 检查IP区间内是否包含IP
func Contains(start, end, ip net.IP) bool {
	return !IPLess(ip, start) && IPLess(ip, NextIP(end))
}
