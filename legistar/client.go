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
)

const apiBase = "https://webapi.legistar.com/v1/"

type Client struct {
	Client string // i.e. Client in http://webapi.legistar.com/v1/{Client}

	Token      string
	HttpClient *http.Client
}

// Return all Persons
func (c Client) Persons() (Persons, error) {
	var data Persons
	return data, c.Call("/Persons", nil, &data)
}

func (c Client) Person(ID int) (Person, error) {
	var p Person
	return p, c.Call(fmt.Sprintf("/Persons/%d", ID), nil, &p)
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
		return fmt.Errorf("HTTP %s", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(&data)
}
