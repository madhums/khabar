package time

import "time"

func EpochHours() int64 {
	now := time.Now()
	return 3600 * (now.Unix() / 3600)
}

func EpochNow() int64 {
	now := time.Now()
	return now.UnixNano() / 1000000 //Convert to Milliseconds
}
