package util

import (
	"fmt"
	"time"
)

func GetCurrentUtcTime(value int) time.Time {
	fixedZone := time.FixedZone(fmt.Sprintf("UTC%d", value), value*60*60)
	return time.Now().UTC().In(fixedZone)
}
