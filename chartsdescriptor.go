package graphx

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ChartsDescriptor is a client request to begin streaming metrics from the
// configured chart names.
type ChartsDescriptor struct {
	// an optional timesestamp indicating the client would like historical metrics
	Fill TimeStamp `json:"fill"`
	// the named charts a user has configured graphx with.
	ChartNames []string `json:"chart_names" validate:"required,min=1"`
	// list of names to return the above chart metrics for
	Names []string `json:"names" validate:"required,min=1"`
	// the value we poll the backend datastore and provide the client with updated metrics
	PollInterval Duration `json:"poll_interval" validate:"required"`
}

// ChartName is a type faciliating marshaling and unmarshaling a string to our ChartName type
type ChartName string

func (cn ChartName) MarshalJSON() ([]byte, error) {
	str := string(cn)
	qs := strconv.Quote(str)
	return []byte(qs), nil
}

func (cn *ChartName) UnmarshalJSON(value []byte) error {
	s := value
	str := strings.Trim(string(s), "\"")
	*cn = ChartName(str)
	return nil
}

// TimeStamp is a type which faciliates marshaling and unmarshaling of rfc3339 timestamps
type TimeStamp time.Time

func (t TimeStamp) MarshalJSON() ([]byte, error) {
	s := time.Time(t).Format(time.RFC3339)
	return []byte(s), nil
}

func (t *TimeStamp) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	ts, err := time.Parse(time.RFC3339, str)
	*t = TimeStamp(ts)
	if err != nil {
		return err
	}

	return nil
}

// Duration is a type which faciliates unmarshaling and marshaling a duration format string such as 5s
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%v", time.Duration(d))
	return []byte(s), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	td, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	*d = Duration(td)
	return nil
}
