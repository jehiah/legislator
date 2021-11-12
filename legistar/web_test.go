package legistar

import (
	"context"
	"net/url"
	"testing"
)

func TestLookupWebURL(t *testing.T) {
	c := NewClient("nyc", "")
	c.LookupURL, _ = url.Parse("https://legistar.council.nyc.gov/gateway.aspx?m=l&id=")
	u, err := c.LookupWebURL(context.Background(), 56676)
	if err != nil {
		t.Fatal(err)
	}
	if u != "https://legistar.council.nyc.gov/LegislationDetail.aspx?ID=3713951&GUID=E7B03ABA-8F42-4341-A0D2-50E2F95320CD" {
		t.Errorf("got %q", u)
	}
}
