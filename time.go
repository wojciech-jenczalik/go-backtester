package main

import (
	"errors"
	"math"
	"time"
)

const (
	dateLayout = "2006-01-02"
)

var timeParseError = errors.New("error while parsing Date. Date should have layout: " + dateLayout)

func convertTimeToQuarters(from time.Time, to time.Time) (period string, limit int) {
	duration := to.Sub(from).Hours() / 24.0
	limit = int(math.Ceil(duration / 30.0 / 3))
	return PeriodQuarter, limit
}

func convertTimeToYears(from time.Time, to time.Time) (period string, limit int) {
	duration := to.Sub(from).Hours() / 24.0
	limit = int(math.Ceil(duration / 30.0 / 12))
	return periodAnnual, limit
}
