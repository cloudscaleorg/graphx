package unixtime

import (
	"strconv"
	"time"
)

// ParseString will take a unix string and convert it to a time.Time value
func ParseString(epoch string) (time.Time, error) {
	// do we have a decimal
	var decimal bool

	for _, r := range []rune(epoch) {
		if r == '.' {
			decimal = true
		}
	}

	var i int64
	if decimal {
		f, err := strconv.ParseFloat(epoch, 64)
		if err != nil {
			return time.Time{}, err
		}
		i = int64(f)
	} else {
		ii, err := strconv.ParseInt(epoch, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		i = int64(ii)
	}

	t := time.Unix(i, 0)
	t = t.UTC()
	return t, nil

}

func ParseInt64(epoch int64) (time.Time, error) {
	t := time.Unix(epoch, 0)
	t = t.UTC()
	return t, nil
}
