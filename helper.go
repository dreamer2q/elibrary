package main

import (
	"ebook/douban"
	"strings"
	"time"
)

func parseTime(d string) (t time.Time) {
	if strings.Count(d, "-") == 2 {
		t, _ = time.Parse("2006-1-2", d)
	} else {
		t, _ = time.Parse("2006-1", d)
	}
	return
}
func intWrapper(i int, _ error) int {
	return i
}

func wrapBookTags(tags []douban.BookTag) (retTags []BookTag) {
	for _, t := range tags {
		retTags = append(retTags, BookTag{
			Title: t.Title,
			Name:  t.Name,
			Count: t.Count,
		})
	}
	return
}
