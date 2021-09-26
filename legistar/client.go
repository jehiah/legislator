// legistar package is an API client for http://webapi.legistar.com/
package legistar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const apiBase = "https://webapi.legistar.com/v1/"

type Filter struct {
	Skip   int    // $skip
	Filter string // $filter
	Top    int    // $top=n (like limit?)
}

type Client struct {
	Client string // i.e. Client in http://webapi.legistar.com/v1/{Client}

	Token      string
	HttpClient *http.Client
}

type Filters interface {
	Paramters() url.Values
}

func dateTimeFilter(field string, t time.Time) url.Values {
	if t.IsZero() {
		return url.Values{}
	}
	v := fmt.Sprintf("%s gt datetime'%s'", field, t.Format("2006-01-02T15:04:05.999999999"))
	return url.Values{"$filter": []string{v}}
}

// Return all Persons
func (c Client) Persons(f Filters) (Persons, error) {
	var data Persons
	var p url.Values
	if f != nil {
		p = f.Paramters()
	}
	return data, c.Call("/Persons", p, &data)
}

func (c Client) Person(ID int) (Person, error) {
	var p Person
	return p, c.Call(fmt.Sprintf("/Persons/%d", ID), nil, &p)
}
func (c Client) PersonVotes(ID int) (Votes, error) {
	// TODO: page
	var v Votes
	return v, c.Call(fmt.Sprintf("/Persons/%d/Votes", ID), nil, &v)
}
func (c Client) PersonOfficeRecords(ID int) (OfficeRecords, error) {
	var v OfficeRecords
	return v, c.Call(fmt.Sprintf("/Persons/%d/OfficeRecords", ID), nil, &v)
}

func (c Client) OfficeRecords(f Filters) (OfficeRecords, error) {
	var v OfficeRecords
	var p url.Values
	if f != nil {
		p = f.Paramters()
	}
	return v, c.Call("/OfficeRecords", p, &v)
}

// VoteTypes
// http://webapi.legistar.com/Help/Api/GET-v1-Client-VoteTypes
func (c Client) VoteTypes() (VoteTypes, error) {
	var v VoteTypes
	return v, c.Call("/VoteTypes", nil, &v)
}

func (c Client) MatterTypes() (MatterTypes, error) {
	var v MatterTypes
	return v, c.Call("/MatterTypes", nil, &v)
}

func (c Client) MatterIndexes() (MatterIndexes, error) {
	var v MatterIndexes
	return v, c.Call("/MatterIndexes", nil, &v)
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

func (c Client) Call(endpoint string, params url.Values, data interface{}) error {
	h := c.HttpClient
	if h == nil {
		h = http.DefaultClient
	}
	u := apiBase + c.Client + endpoint
	if params == nil {
		params = url.Values{}
	}
	params.Set("token", c.Token)
	if strings.Contains(u, "?") {
		u += "&" + params.Encode()
	} else {
		u += "?" + params.Encode()
	}
	log.Printf("GET %s", u)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
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
