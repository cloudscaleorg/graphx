package unixtime

import (
	"testing"
	"time"
)

const (
	DefaultTSFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

var UnitTimeTT = []struct {
	epoch            string
	expectedTSString string
	shouldError      bool
}{
	{
		epoch:            "1548785535.104",
		expectedTSString: "2019-01-29 18:12:15 +0000 UTC",
		shouldError:      false,
	},
	{
		epoch:            "1548785535",
		expectedTSString: "2019-01-29 18:12:15 +0000 UTC",
		shouldError:      false,
	},
}

func TestParseString(t *testing.T) {
	for _, tt := range UnitTimeTT {

		t.Run("", func(t *testing.T) {
			ts, err := ParseString(tt.epoch)
			if err != nil {
				if !tt.shouldError {
					t.Fatalf("failed to parse epoch: %v", err)
				}
			}

			expectedTS, err := time.Parse(DefaultTSFormat, tt.expectedTSString)
			if err != nil {
				t.Fatalf("failed to parse expected time stamp string to time.Time: %v", err)
			}

			if ts != expectedTS {
				t.Fatalf("expected ts: %v got ts: %v", expectedTS, ts)
			}
		})

	}
}
