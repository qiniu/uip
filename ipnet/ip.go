package ipnet

import (
	"bytes"
	"net"
)

// Uint32ToIPv4 uint32 转换为 IP
func Uint32ToIPv4(uip uint32) net.IP {
	return net.IPv4(byte(uip>>24&0xFF), byte(uip>>16&0xFF), byte(uip>>8&0xFF), byte(uip&0xFF))
}

// IPv4ToUint32 IPv4 转换为 uint32
func IPv4ToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}

// IPv4StrToUint32 IPv4字符串 转换为 uint32 (syntactic sugar)
func IPv4StrToUint32(ipStr string) uint32 {
	if ip := net.ParseIP(ipStr).To4(); ip != nil {
		return IPv4ToUint32(ip)
	}
	return 0
}

// Uint64ToIP uint64 转换为 IP
func Uint64ToIP(uip uint64) net.IP {
	return net.IP{
		byte(uip >> 56 & 0xFF), byte(uip >> 48 & 0xFF), byte(uip >> 40 & 0xFF), byte(uip >> 32 & 0xFF),
		byte(uip >> 24 & 0xFF), byte(uip >> 16 & 0xFF), byte(uip >> 8 & 0xFF), byte(uip & 0xFF),
		0, 0, 0, 0, 0, 0, 0, 0,
	}
}

// Uint64ToIP2 uint64 转换为 IP
func Uint64ToIP2(high, low uint64) net.IP {
	return net.IP{
		byte(high >> 56 & 0xFF), byte(high >> 48 & 0xFF), byte(high >> 40 & 0xFF), byte(high >> 32 & 0xFF),
		byte(high >> 24 & 0xFF), byte(high >> 16 & 0xFF), byte(high >> 8 & 0xFF), byte(high & 0xFF),
		byte(low >> 56 & 0xFF), byte(low >> 48 & 0xFF), byte(low >> 40 & 0xFF), byte(low >> 32 & 0xFF),
		byte(low >> 24 & 0xFF), byte(low >> 16 & 0xFF), byte(low >> 8 & 0xFF), byte(low & 0xFF),
	}
}

func allFF(b []byte) bool {
	for _, c := range b {
		if c != 0xff {
			return false
		}
	}
	return true
}

func IsAllMask(ip net.IP) bool {
	return allFF(ip)
}

// NextIP 下一个IP
func NextIP(ip net.IP) net.IP {
	res := make(net.IP, len(ip))
	for i := len(ip) - 1; i >= 0; i-- {
		res[i] = ip[i] + 1
		if res[i] != 0 {
			copy(res, ip[0:i])
			break
		}
	}
	return res
}

// IPLess IP大小比较
func IPLess(a, b net.IP) bool {
	left := a
	right := b
	if lenA, lenB := len(a), len(b); lenA != lenB {
		if lenA < lenB {
			left = a.To16()
		} else {
			right = b.To16()
		}
	}
	return bytes.Compare(left, right) < 0
}

// IsFirstIP 是否是起始IP
func IsFirstIP(ip net.IP, ipv6 bool) bool {
	if ipv6 {
		return ip.Equal(net.IPv6zero)
	}
	if len(ip) == net.IPv6len {
		return ip[12] == 0 && ip[13] == 0 && ip[14] == 0 && ip[15] == 0
	}
	return ip.Equal(net.IPv4zero)
}
