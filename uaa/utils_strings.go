package uaa

import (
	"hash/crc32"
)

// resourceStringHash -
func resourceStringHash(si interface{}) int {
	v := int(crc32.ChecksumIEEE([]byte(si.(string))))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
