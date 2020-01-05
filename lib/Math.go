package lib

import (
	"github.com/ethereum/go-ethereum/common/math"
	"math/rand"
	"strconv"
)

const (
	B  uint64 = 1
	KB        = B << 10
	MB        = KB << 10
	GB        = MB << 10
	TB        = GB << 10
	PB        = TB << 10
	EB        = PB << 10
)

//RandomBetween Create random number between two ranges
func RandomBetween(min, max int) int {
	return rand.Intn(max-min) + min
}

func ParseSafeInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func ParseSafeInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func ParseSafeFloat(s string) float64 {
	i, _ := strconv.ParseFloat(s, 10)
	return i
}

func MaxInt(args ...int) int {
	x := math.MinInt32
	for _, arg := range args {
		if arg > x {
			x = arg
		}
	}
	return x
}

func MaxUInt(args ...uint32) uint32 {
	var x uint32
	for _, arg := range args {
		if arg > x {
			x = arg
		}
	}
	return x
}

func MaxInt64(args ...int64) int64 {
	var x int64
	x = math.MinInt64
	for _, arg := range args {
		if arg > x {
			x = arg
		}
	}
	return x
}

func MinInt(args ...int) int {
	x := math.MaxInt32
	for _, arg := range args {
		if arg < x {
			x = arg
		}
	}
	return x
}

func MinUInt(args ...uint) uint {
	var x uint
	x = math.MaxUint32
	for _, arg := range args {
		if arg < x {
			x = arg
		}
	}
	return x
}

func MinInt64(args ...int64) int64 {
	var x int64
	x = math.MaxInt64
	for _, arg := range args {
		if arg < x {
			x = arg
		}
	}
	return x
}

func MinUInt64(args ...uint64) uint64 {
	var x uint64
	x = math.MaxUint64
	for _, arg := range args {
		if arg < x {
			x = arg
		}
	}
	return x
}

func AvgInt(args ...int) float64 {
	return float64(SumInt(args...)) / float64(len(args))
}

func AvgInt64(args ...int64) float64 {
	return float64(SumInt64(args...)) / float64(len(args))
}

func AvgUInt(args ...uint) float64 {
	return float64(SumUInt(args...)) / float64(len(args))
}

func AvgUInt64(args ...uint64) float64 {
	return float64(SumUInt64(args...)) / float64(len(args))
}

func SumInt(args ...int) int {
	var sum int
	for _, arg := range args {
		sum += arg
	}
	return sum
}

func SumInt64(args ...int64) int64 {
	var sum int64
	for _, arg := range args {
		sum += arg
	}
	return sum
}

//SumUInt Sum of uint32
func SumUInt(args ...uint) uint {
	var sum uint
	for _, arg := range args {
		sum += arg
	}
	return sum
}

func SumUInt64(args ...uint64) uint64 {
	var sum uint64
	for _, arg := range args {
		sum += arg
	}
	return sum
}
