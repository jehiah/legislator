package main

import (
	"html/template"
	"path"
	"strings"

	"github.com/dustin/go-humanize"
)

type Contains interface {
	Contains(needle string) bool
}

func strArrayContains(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
func contains(c Contains, needle string) bool {
	return c.Contains(needle)
}

func commaInt(i int) string {
	return humanize.Comma(int64(i))
}

func compileTemplates(base string) (t *template.Template) {

	t = template.Must(template.New("").Funcs(template.FuncMap{
		"SI":               humanize.SI,
		"Comma":            humanize.Comma,
		"Commai":           commaInt,
		"Commaf":           humanize.Commaf,
		"Bytes":            humanize.Bytes,
		"Time":             humanize.Time,
		"Title":            strings.Title,
		"strArrayContains": strArrayContains,
		"contains":         contains,
	}).ParseGlob(path.Join(base, "*.html")))
	return
}
