package main

import (
	"testing"
	"time"
)

func TestGiven_interval_13_months_when_converted_period_should_be_5_quarters(t *testing.T) {
	from, _ := time.Parse(dateLayout, "2020-01-01")
	to, _ := time.Parse(dateLayout, "2021-02-01")

	expectedPeriod, expectedLimit := PeriodQuarter, 5
	period, limit := convertTimeToQuarters(from, to)

	if period != expectedPeriod || limit != expectedLimit {
		t.Fatalf(
			"From: %s, To: %s. Wanted period: %s, limit: %d. Actual period %s, limit: %d",
			from.String(),
			to.String(),
			expectedPeriod,
			expectedLimit,
			period,
			limit)
	}
}

func TestGiven_interval_1_day_when_converted_period_should_be_1_quarter(t *testing.T) {
	from, _ := time.Parse(dateLayout, "2021-01-01")
	to, _ := time.Parse(dateLayout, "2021-01-02")

	expectedPeriod, expectedLimit := PeriodQuarter, 1
	period, limit := convertTimeToQuarters(from, to)

	if period != expectedPeriod || limit != expectedLimit {
		t.Fatalf(
			"From: %s, To: %s. Wanted period: %s, limit: %d. Actual period %s, limit: %d",
			from.String(),
			to.String(),
			expectedPeriod,
			expectedLimit,
			period,
			limit)
	}
}

func TestGiven_interval_45_days_when_converted_period_should_be_1_quarter(t *testing.T) {
	from, _ := time.Parse(dateLayout, "2021-01-01")
	to, _ := time.Parse(dateLayout, "2021-02-15")

	expectedPeriod, expectedLimit := PeriodQuarter, 1
	period, limit := convertTimeToQuarters(from, to)

	if period != expectedPeriod || limit != expectedLimit {
		t.Fatalf(
			"From: %s, To: %s. Wanted period: %s, limit: %d. Actual period %s, limit: %d",
			from.String(),
			to.String(),
			expectedPeriod,
			expectedLimit,
			period,
			limit)
	}
}

func TestGiven_interval_3690_days_when_converted_period_should_be_41_quarters(t *testing.T) {
	from, _ := time.Parse(dateLayout, "2021-01-01")
	to, _ := time.Parse(dateLayout, "2031-02-08")

	expectedPeriod, expectedLimit := PeriodQuarter, 41
	period, limit := convertTimeToQuarters(from, to)

	if period != expectedPeriod || limit != expectedLimit {
		t.Fatalf(
			"From: %s, To: %s. Wanted period: %s, limit: %d. Actual period %s, limit: %d",
			from.String(),
			to.String(),
			expectedPeriod,
			expectedLimit,
			period,
			limit)
	}
}

func TestGiven_interval_1_day_when_converted_period_should_be_1_year(t *testing.T) {
	from, _ := time.Parse(dateLayout, "2021-01-01")
	to, _ := time.Parse(dateLayout, "2021-01-02")

	expectedPeriod, expectedLimit := PeriodAnnual, 1
	period, limit := convertTimeToYears(from, to)

	if period != expectedPeriod || limit != expectedLimit {
		t.Fatalf(
			"From: %s, To: %s. Wanted period: %s, limit: %d. Actual period %s, limit: %d",
			from.String(),
			to.String(),
			expectedPeriod,
			expectedLimit,
			period,
			limit)
	}
}

func TestGiven_interval_3690_days_when_converted_period_should_be_11_years(t *testing.T) {
	from, _ := time.Parse(dateLayout, "2021-01-01")
	to, _ := time.Parse(dateLayout, "2031-02-08")

	expectedPeriod, expectedLimit := PeriodAnnual, 11
	period, limit := convertTimeToYears(from, to)

	if period != expectedPeriod || limit != expectedLimit {
		t.Fatalf(
			"From: %s, To: %s. Wanted period: %s, limit: %d. Actual period %s, limit: %d",
			from.String(),
			to.String(),
			expectedPeriod,
			expectedLimit,
			period,
			limit)
	}
}
