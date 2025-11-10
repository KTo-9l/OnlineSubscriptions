package main

import "time"

func parseTime(stringTime string) (t time.Time, err error) {
	return time.Parse("01-2006", stringTime)
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.After(b) {
		return b
	}
	return a
}
