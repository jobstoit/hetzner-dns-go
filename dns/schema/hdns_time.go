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
		t2, err2 := time.Parse(time.RFC3339, s)

		if err2 != nil {
			return fmt.Errorf("Error while parsing date '%s' with default rfc3339 time layout and '%s': %s\n", s, hdns_time_layout, err)
		}
		t = t2
	}

	*ht = HdnsTime(t)
	return nil
}

func (ht HdnsTime) MarshalJSON() ([]byte, error) {
	return time.Time(ht).MarshalJSON()
}
