package main

import "net/url"

func UpdateInSiteURL(host, InSiteURL string) string {
	u, err := url.Parse(InSiteURL)
	if err != nil {
		return InSiteURL
	}
	u.Host = host
	return u.String()
}
