package utility

import (
	"fmt"
	"strings"
	"time"
)

const ZendeskTimeFormat = "2006-01-02T15:04:05 -07:00"

var (
	nilTime = (time.Time{}).UnixNano()
)

type ZendeskTime struct {
	time.Time
}

func (ct *ZendeskTime) UnmarshalJSON(b []byte) (err error) {
	// Strip json stuff
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return nil
	}

	ct.Time, err = time.Parse(ZendeskTimeFormat, s)
	return err
}

func (ct *ZendeskTime) MarshalJSON() ([]byte, error) {
	if ct.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(ZendeskTimeFormat))), nil
}

func (ct *ZendeskTime) IsSet() bool {
	return ct.UnixNano() != nilTime
}
