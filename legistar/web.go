package legistar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// LookupWebURL returns the Web URL for a MatterID
func (c Client) LookupWebURL(ctx context.Context, MatterID int) (string, error) {
	if c.LookupURL == nil {
		return "", fmt.Errorf("not configured")
	}

	err := c.Limiter.Wait(ctx)
	if err != nil {
		return "", err
	}

	h := c.HttpClient
	if h == nil {
		// should disable redirects
		h = http.DefaultClient
	}

	u := c.LookupURL.String() + url.QueryEscape(strconv.Itoa(MatterID))

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := h.Do(req)
	if err != nil {
		return "", err
	}
	switch resp.StatusCode {
	case 301, 302:
		r, err := url.Parse(resp.Header.Get("Location"))
		if err != nil {
			return "", err
		}
		return c.LookupURL.ResolveReference(r).String(), nil
	default:
		return "", fmt.Errorf("unexpected response code %d", resp.StatusCode)
	}
}
