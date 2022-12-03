package legistar

import "time"

const iso8601 = "2006-01-02T15:04:05.999999999"

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t.Time, err = time.Parse(`"`+iso8601+`"`, string(data))
	return err
}

type ShortTime struct {
	time.Time
}

func (t *ShortTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	// log.Printf("Shorttime %s", string(data))
	// Fractional seconds are handled implicitly by Parse.
	var err error
	t.Time, err = time.Parse(`"3:04 PM"`, string(data))
	return err
}

func (s ShortTime) Set(t time.Time, tz *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), s.Hour(), s.Minute(), s.Second(), s.Nanosecond(), tz)
}
