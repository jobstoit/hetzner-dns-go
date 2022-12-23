package schema

import (
	"fmt"
	"strings"
	"time"
)

// HdnsTime defines a wrapper for time.Time, to handle the DateTime format(s) used in HCloud DNS API...
type HdnsTime time.Time

func (ht *HdnsTime) UnmarshalJSON(b []byte) error {

	if len(b) == 0 {
		*ht = HdnsTime{}
		return nil
	}

	s := strings.Trim(string(b), "\"")

	if len(s) == 0 {
		*ht = HdnsTime{}
		return nil
	}

	hdns_time_layout_1 := "2006-01-02 15:04:05 -0700 UTC"

	if t, err := time.Parse(hdns_time_layout_1, s); err == nil {
		*ht = HdnsTime(t)
		return nil
	}

	// Real example API return value:
	// [...] "verified": "2020-04-07 01:56:03.196438163 +0000 UTC m=+755.322810452", [...]
	hdns_time_layout_2 := "2006-01-02 15:04:05.000000000 -0700 UTC" //... m=+755.322810452"

	if t, err := time.Parse(hdns_time_layout_2, strings.Split(s, " m=+")[0]); err == nil {
		*ht = HdnsTime(t)
		return nil
	}

	if t, err := time.Parse(time.RFC3339, s); err == nil {
		*ht = HdnsTime(t)
		return nil
	}

	return fmt.Errorf("Error while parsing date '%s'\n", s)
}

func (ht HdnsTime) MarshalJSON() ([]byte, error) {
	return time.Time(ht).MarshalJSON()
}
