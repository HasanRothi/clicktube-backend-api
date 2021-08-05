package services

import (
	"strconv"
	"time"
)

func GenarateSortLink(next int, CAMPUS_CODE string, DEPT_CODE string) string {
	dateSlice := UrlKey()
	nextLink := next + 1
	shortUrl := dateSlice + "-" + CAMPUS_CODE + "-" + DEPT_CODE + "-" + strconv.Itoa(nextLink)
	return shortUrl
}
func UrlKey() string {
	currentDate := time.Now()
	dateString := currentDate.String()
	dateSlice := dateString[0:10]
	return dateSlice
}
