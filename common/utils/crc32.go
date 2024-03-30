package utils

import (
	"hash/crc32"
)

func validCrc32sum(crc32sum uint32, data ...[]byte) bool {
	crc32hash := crc32.NewIEEE()
	for _, d := range data {
		_, _ = crc32hash.Write(d)
	}
	return crc32sum == crc32hash.Sum32()
}
