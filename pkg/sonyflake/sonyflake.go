package sonyflake

import (
	"github.com/sony/sonyflake"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake init panic")
	}
}

func ID() (uint64, error) {
	return sf.NextID()
}
