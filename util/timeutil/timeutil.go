package timeutil

import (
	"time"
)

//Now time round by second
func Now() time.Time {
	return time.Now().Round(time.Second)
}

//NanoTime from nano
func NanoTime(nano int64) time.Time {
	return time.Unix(0, nano)
}
