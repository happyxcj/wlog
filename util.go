package wlog

const lowerDigit = "0123456789abcdef"

// uin64ToHex returns the hexadecimal value of the given v.
func uin64ToHex(v uint64) []byte {
	// buf is large enough to store %b of an uint64 with prefix "0x".
	var buf [66]byte
	// Printing is easier right-to-left: format v into buf, ending at buf[i].
	i := 66
	for v >= 16 {
		i--
		buf[i] = lowerDigit[v&0xF]
		v >>= 4
	}
	i--
	buf[i] = lowerDigit[v]
	i--
	buf[i] = 'x'
	i--
	buf[i] = '0'
	return buf[i:]
}
