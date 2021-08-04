package services

import (
	"strconv"
	"time"
)

func GenarateSortLink(next int) string {
	dateSlice := UrlKey()
	CAMPUS_CODE := "UITS"
	DEPT_CODE := "IT"
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
