package xstring

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

var timeBase = time.Date(1582, time.October, 15, 0, 0, 0, 0, time.UTC).Unix()
var hardwareAddr []byte
var clockSeq uint32
var randInstance *rand.Rand

func init() {
	randInstance = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GenerateUUID simply generates an unique UID.
func GenerateUUID(seedTime time.Time) string {
	var u [16]byte
	utcTime := seedTime.In(time.UTC)
	t := uint64(utcTime.Unix()-timeBase)*10000000 + uint64(utcTime.Nanosecond()/100)

	u[0], u[1], u[2], u[3] = byte(t>>24), byte(t>>16), byte(t>>8), byte(t)
	u[4], u[5] = byte(t>>40), byte(t>>32)
	u[6], u[7] = byte(t>>56)&0x0F, byte(t>>48)

	clock := atomic.AddUint32(&clockSeq, 1)
	u[8] = byte(clock >> 8)
	u[9] = byte(clock)

	copy(u[10:], hardwareAddr)

	u[6] |= 0x10 // set version to 1 (time based uuid)
	u[8] &= 0x3F // clear variant
	u[8] |= 0x80 // set to IETF variant

	var offsets = [...]int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30}
	const hexString = "0123456789abcdef"
	r := make([]byte, 32)
	for i, b := range u {
		r[offsets[i]] = hexString[b>>4]
		r[offsets[i]+1] = hexString[b&0xF]
	}
	return string(r)
}

// GenerateID simply generates an ID.
func GenerateID() string {
	return fmt.Sprintf("%016x", uint64(randInstance.Int63()))
}
