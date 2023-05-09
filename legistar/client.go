// legistar package is an API client for http://webapi.legistar.com/
package legistar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const apiBase = "https://webapi.legistar.com/v1/"

var userAgent = fmt.Sprintf("Go-http-client/%s (https://github.com/jehiah/legislator)", strings.TrimPrefix(runtime.Version(), "go"))

type Filter struct {
	Skip   int    // $skip
	Filter string // $filter
	Top    int    // $top=n (like limit?)
}

type Client struct {
	Client string // i.e. Client in http://webapi.legistar.com/v1/{Client}

	LookupURL *url.URL

	Token      string
	HttpClient *http.Client
	Limiter    *rate.Limiter
}

func NewClient(client, token string) *Client {
	h := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return &Client{
		Client:     client,
		HttpClient: h,
		Token:      token,
		Limiter:    rate.NewLimiter(rate.Every(20*time.Millisecond), 20),
	}
}

// Return all Persons
func (c Client) Persons(ctx context.Context, f Filters) (Persons, error) {
	var data Persons
	var p url.Values
	if f != nil {
		p = f.Paramters()
	}
	return data, c.Call(ctx, "/Persons", p, &data)
}

func (c Client) Person(ctx context.Context, ID int) (Person, error) {
	var p Person
	return p, c.Call(ctx, fmt.Sprintf("/Persons/%d", ID), nil, &p)
}
func (c Client) PersonVotes(ctx context.Context, ID int, f Filters) (Votes, error) {
	var p url.Values
	if f != nil {
		p = f.Paramters()
	}
	var v Votes
	return v, c.Call(ctx, fmt.Sprintf("/Persons/%d/Votes", ID), p, &v)
}
func (c Client) PersonOfficeRecords(ctx context.Context, ID int) (OfficeRecords, error) {
	var v OfficeRecords
	return v, c.Call(ctx, fmt.Sprintf("/Persons/%d/OfficeRecords", ID), nil, &v)
}

func (c Client) PersonRollCalls(ctx context.Context, ID int, f Filters) (RollCalls, error) {
	var p url.Values
	if f != nil {
		p = f.Paramters()
	}
	var v RollCalls
	return v, c.Call(ctx, fmt.Sprintf("/Persons/%d/RollCalls", ID), p, &v)
}

func (c Client) OfficeRecords(ctx context.Context, f Filters) (OfficeRecords, error) {
	var v OfficeRecords
	var p url.Values
	if f != nil {
		p = f.Paramters()
	}
	return v, c.Call(ctx, "/OfficeRecords", p, &v)
}

func (c Client) Matters(ctx context.Context, f Filters) (Matters, error) {
	var v Matters
	var p url.Values
	if f != nil {
		p = f.Paramters()
	}
	return v, c.Call(ctx, "/Matters", p, &v)
}

func (c Client) Matter(ctx context.Context, ID int) (Matter, error) {
	var v Matter
	return v, c.Call(ctx, fmt.Sprintf("/Matters/%d", ID), nil, &v)
}

func (c Client) MatterSponsors(ctx context.Context, ID int) (MatterSponsors, error) {
	var v MatterSponsors
	err := c.Call(ctx, fmt.Sprintf("/Matters/%d/Sponsors", ID), nil, &v)
	sort.Slice(v, func(i, j int) bool { return v[i].Sequence < v[j].Sequence })
	return v, err
}

func (c Client) MatterHistories(ctx context.Context, ID int) (MatterHistories, error) {
	var v MatterHistories
	err := c.Call(ctx, fmt.Sprintf("/Matters/%d/Histories?AgendaNote=1&MinutesNote=1", ID), nil, &v)
	sort.Slice(v, func(i, j int) bool {
		if v[i].ActionDate.Time.Equal(v[j].ActionDate.Time) {
			return v[i].ID < v[j].ID
		}
		return v[i].ActionDate.Time.Before(v[j].ActionDate.Time)
	})
	return v, err
}

func (c Client) MatterText(ctx context.Context, matterID, textID int) (MatterText, error) {
	var v MatterText
	if textID == 0 {
		return v, errors.New("got textID 0")
	}
	return v, c.Call(ctx, fmt.Sprintf("/Matters/%d/Texts/%d", matterID, textID), nil, &v)
}

func (c Client) MatterTextVersions(ctx context.Context, matterID int) (MatterTextVersions, error) {
	var v MatterTextVersions
	return v, c.Call(ctx, fmt.Sprintf("/Matters/%d/Versions", matterID), nil, &v)
}

func (c Client) MatterAttachments(ctx context.Context, matterID int) (MatterAttachments, error) {
	var v MatterAttachments
	err := c.Call(ctx, fmt.Sprintf("/Matters/%d/Attachments", matterID), nil, &v)
	sort.Slice(v, func(i, j int) bool { return v[i].Sort < v[j].Sort })
	return v, err
}

// VoteTypes
// http://webapi.legistar.com/Help/Api/GET-v1-Client-VoteTypes
func (c Client) VoteTypes(ctx context.Context) (VoteTypes, error) {
	var v VoteTypes
	return v, c.Call(ctx, "/VoteTypes", nil, &v)
}

func (c Client) MatterTypes(ctx context.Context) (MatterTypes, error) {
	var v MatterTypes
	return v, c.Call(ctx, "/MatterTypes", nil, &v)
}

func (c Client) MatterIndexes(ctx context.Context) (MatterIndexes, error) {
	var v MatterIndexes
	return v, c.Call(ctx, "/MatterIndexes", nil, &v)
}

// Events
// http://webapi.legistar.com/Help/Api/GET-v1-Client-Events
func (c Client) Events(ctx context.Context, f Filters) (Events, error) {
	var p url.Values
	if f != nil {
		p = f.Paramters()
	}
	var v Events
	return v, c.Call(ctx, "/Events", p, &v)
}

// Event
// http://webapi.legistar.com/Help/Api/GET-v1-Client-Events-EventId_EventItems_AgendaNote_MinutesNote_EventItemAttachments
func (c Client) Event(ctx context.Context, ID int) (Event, error) {
	var v Event
	return v, c.Call(ctx, fmt.Sprintf("/Events/%d", ID), nil, &v)
}


// EventItems
// http://127.0.0.1:7001/Events/379233/EventItems?AgendaNote=1&MinutesNote=1&Attachments=1
func (c Client) EventItems(ctx context.Context, ID int) (EventItems, error) {
	var v EventItems
	err := c.Call(ctx, fmt.Sprintf("/Events/%d/EventItems?AgendaNote=1&MinutesNote=1&Attachments=1", ID), nil, &v)
	// TODO: sort
	// sort.Slice(v, func(i, j int) bool {
	// 	if v[i].ActionDate.Time.Equal(v[j].ActionDate.Time) {
	// 		return v[i].ID < v[j].ID
	// 	}
	// 	return v[i].ActionDate.Time.Before(v[j].ActionDate.Time)
	// })
	return v, err
}

type apiError struct {
	code    int
	message string
}

func (e apiError) Error() string {
	return fmt.Sprintf("HTTP %s", e.message)
}

func IsNotFoundError(err error) bool {
	if e, ok := err.(apiError); ok && e.code == http.StatusNotFound {
		return true
	}
	return false
}

func (c Client) Call(ctx context.Context, endpoint string, params url.Values, data interface{}) error {
	err := c.Limiter.Wait(ctx)
	if err != nil {
		return err
	}

	h := c.HttpClient
	if h == nil {
		h = http.DefaultClient
	}
	u := apiBase + c.Client + endpoint
	if params == nil {
		params = url.Values{}
	}
	if strings.Contains(u, "?") {
		log.Printf("GET %s", u+"&"+params.Encode())
	} else {
		log.Printf("GET %s", u+"?"+params.Encode())
	}
	params.Set("token", c.Token)
	if strings.Contains(u, "?") {
		u += "&" + params.Encode()
	} else {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := h.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("got %q %q for %s", resp.Status, string(body), resp.Request.URL.String())
		return apiError{code: resp.StatusCode, message: resp.Status}
	}
	return json.NewDecoder(resp.Body).Decode(&data)
}
