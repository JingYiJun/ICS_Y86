package utils

import "strconv"

func Hex2uint64(b []byte) uint64 {
	v, err := strconv.ParseInt(string(b), 16, 16)
	if err != nil {
		panic(err)
	}
	return uint64(v)
}

const (
	ZF = iota
	SF
	OF
)
