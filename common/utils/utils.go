package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func MustParseInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func MustParseFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func CreateOrderSn() string {
	currentTime := time.Now()
	date := currentTime.Format("20060102")
	uid := fmt.Sprintf("%08d", rand.Intn(99999999)) // 8位随机数
	seconds := currentTime.Unix() % 10000000000     // 10位秒级时间戳
	return fmt.Sprintf("%s%s%010d", date, uid, seconds)
}

func ScheduleNextTask(hours []int, early time.Duration) time.Time {
	now := time.Now()
	var earliestNext time.Time

	// 找到下一个最近的执行时间
	for _, hour := range hours {
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location()).Add(-early)
		if now.After(next) {
			// 如果当前时间已经过了预定时间，计算下一天的同一时间
			next = next.Add(24 * time.Hour)
		}
		if earliestNext.IsZero() || next.Before(earliestNext) {
			earliestNext = next
		}
	}

	return earliestNext
}
