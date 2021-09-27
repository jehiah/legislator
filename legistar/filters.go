package legistar

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Filters interface {
	Paramters() url.Values
}

func DateTimeFilter(field, eq string, t time.Time) url.Values {
	if t.IsZero() {
		return url.Values{}
	}
	v := fmt.Sprintf("%s %s datetime'%s'", field, eq, t.Format("2006-01-02T15:04:05.999999999"))
	return url.Values{"$filter": []string{v}}
}

func StringFilter(field, value string) url.Values {
	if value == "" {
		return url.Values{}
	}
	return url.Values{"$filter": []string{fmt.Sprintf("%s eq '%s'", field, value)}}
}

func AndFilters(f ...Filters) Filters {
	return andFilters{f: f}
}

type andFilters struct {
	f []Filters
}

func (f andFilters) Paramters() url.Values {
	v := url.Values{}
	for _, ff := range f.f {
		vv := ff.Paramters()
		for k, vvv := range vv {
			for _, vvvv := range vvv {
				v.Add(k, vvvv)
			}
		}
	}
	for k, vv := range v {
		if len(vv) > 1 {
			v[k] = []string{strings.Join(vv, " and ")}
		}
	}
	return v
}
