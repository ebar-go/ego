package number

import (
	"fmt"
	"hash/crc32"
	"math"
	"math/rand"
	"strconv"
	"time"
)

//Min return min number
func Min(a, b int) int {
	if a < b {
		return a
	}
	if b < a {
		return b
	}
	return a
}

//Max return max number
func Max(a, b int) int {
	if a > b {
		return a
	}
	if b > a {
		return b
	}
	return a
}

// Div return a/b
func Div(a, b int) int {
	if b == 0 || a == 0 {
		return 0
	}

	return int(math.Ceil(float64(a) / float64(b)))
}

// Round rounding-off method
func Round(f float64) int {
	return int(math.Floor(f + 0.5))
}

// DefaultInt return default value if v is zero
func DefaultInt(v, dv int) int {
	if v == 0 {
		return dv
	}

	return v
}

// FloatValue
type FloatValue float64

// Int return integer
func (f FloatValue) Int() int {
	str := fmt.Sprintf("%.f", float64(f))
	if result, err := strconv.Atoi(str); err == nil {
		return result
	}

	return 0
}

// Round rounding-off method
func (f FloatValue) Round() int {
	return int(math.Floor(float64(f) + 0.5))
}

// HashCode
func HashCode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

// RandInt 取随机数
func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
