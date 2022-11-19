package schema

import (
	"fmt"
	"strings"
	"time"
)

// HdnsTime defines a wrapper for time.Time, to handle the DateTime format used in HCloud DNS API
type HdnsTime time.Time

func (ht *HdnsTime) UnmarshalJSON(b []byte) (err error) {
	hdns_time_layout := "2006-01-02 15:04:05 -0700 UTC"

	s := strings.Trim(string(b), "\"")

	if len(s) == 0 {
		*ht = HdnsTime{}
		return nil
	}

	t, err := time.Parse(hdns_time_layout, s)
	if err != nil {
		return fmt.Errorf("Error while parsing date '%s' with time layout '%s': %s\n", s, hdns_time_layout, err)
	}

	*ht = HdnsTime(t)
	return nil
}

func (ht HdnsTime) MarshalJSON() ([]byte, error) {
	return time.Time(ht).MarshalJSON()
}
