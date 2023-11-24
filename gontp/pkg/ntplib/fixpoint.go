package ntplib

import "time"

type ShortFixedPoint struct {
	IntPart  uint16
	FracPart uint16
}

type LongFixedPoint struct {
	IntPart  uint32
	FracPart uint32
}

func CreateLongFixedPoint(time time.Time) LongFixedPoint {
	totalUSec := time.UnixMicro()
	sec := totalUSec / 1000000
	uSec := totalUSec % 1000000
	new := LongFixedPoint{
		IntPart:  uint32(sec + 2208988800),
		FracPart: uint32(float64(uSec+1) * float64(int64(1)<<32) * 1.0e-6),
	}
	return new
}

func (timestamp *LongFixedPoint) ToUnixMicro() int64 {
	sec := timestamp.IntPart - 2208988800
	uSec := uint32(float64(timestamp.FracPart) * 1e6 / float64(int64(1)<<32))
	return int64(sec)*1000000 + int64(uSec)
}

func (timestamp *LongFixedPoint) ToTime() time.Time {
	pSec := timestamp.FracPart * 232
	nSec := float64(pSec) / 1000
	return time.Unix(int64(timestamp.IntPart-2208988800), int64(nSec))
}
