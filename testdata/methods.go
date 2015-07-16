package testdata

import (
	"html/template"
	"regexp"
	"time"

	"github.com/lib/pq"
)

var summaryExp *regexp.Regexp = regexp.MustCompile("<p>.*?</p>")

type nonStruct int

type Post struct {
	ID        int
	SemURL    string
	Created   time.Time
	Modified  time.Time
	Published pq.NullTime
	Draft     bool
	Title     string
	Body      string
}

func (p Post) BodyHTML() template.HTML {
	return template.HTML(p.Body)
}

func (p Post) Fmt(t time.Time, f string) string {
	switch f {
	case "date-num":
		return t.Format("2006-01-02")
	case "datetime-num":
		return t.Format("2006-01-02 15:04:05")
	case "date":
		return t.Format("02 January 2006")
	case "datetime":
		return t.Format("02 January 2006 15:04:05")
	}

	return ""
}

func (p Post) ModSincePub() bool {
	return p.Modified.After(p.Published.Time)
}

func (p Post) Summary() template.HTML {
	return template.HTML(summaryExp.FindString(p.Body))
}
