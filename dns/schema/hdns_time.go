package schema

import (
	"fmt"
	"strings"
	"time"
)

// HdnsTime defines a wrapper for time.Time, to handle the DateTime format(s) used in HCloud DNS API...
type HdnsTime time.Time

func (ht *HdnsTime) UnmarshalJSON(b []byte) error {
	hdns_time_layout_1 := "2006-01-02 15:04:05 -0700 UTC"
	// Real example API return value:
	// [...] "verified": "2020-04-07 01:56:03.196438163 +0000 UTC m=+755.322810452", [...]
	hdns_time_layout_2 := "2006-01-02 15:04:05.000000000 -0700 UTC" //... m=+755.322810452"

	s := strings.Trim(string(b), "\"")

	if len(s) == 0 {
		*ht = HdnsTime{}
		return nil
	}

	t, err_format1 := time.Parse(hdns_time_layout_1, s)
	if err_format1 != nil {

		t_format2, err_format2 := time.Parse(hdns_time_layout_2, strings.Split(s, " m=+")[0])
		t = t_format2

		if err_format2 != nil {

			time_rfc3339, err_rfc3339 := time.Parse(time.RFC3339, s)
			t = time_rfc3339

			if err_rfc3339 != nil {
				return fmt.Errorf("Error while parsing date '%s' with\n;; '%s' (err: %s)\n;; '%s' (err: %s)\n;; default rfc3339 time layout (err: %s)\n", s, hdns_time_layout_1, err_format1, hdns_time_layout_2, err_format2, err_rfc3339)
			}
		}
	}

	*ht = HdnsTime(t)
	return nil
}

func (ht HdnsTime) MarshalJSON() ([]byte, error) {
	return time.Time(ht).MarshalJSON()
}
