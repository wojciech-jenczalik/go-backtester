package main

import (
	"reflect"
	"testing"
	"time"
)

// ################# Price determination tests #################

var prices = []Price{
	{Date: "2021-01-10"},
	{Date: "2021-01-09"},
	{Date: "2021-01-06"},
	{Date: "2021-01-05"},
	{Date: "2021-01-04"},
	{Date: "2021-01-03"},
}

func TestDetermine_price_index_on_exact_day(t *testing.T) {
	// Given
	date, _ := time.Parse(dateLayout, "2021-01-09")
	expectedIndex := 1

	// When
	index, err := determinePriceIndexForDate(prices, date)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if index != expectedIndex {
		t.Fatalf("expected index: %d, actual index: %d", expectedIndex, index)
	}
}

func TestDetermine_price_index_on_day_after(t *testing.T) {
	// Given
	date, _ := time.Parse(dateLayout, "2021-01-08")
	expectedIndex := 1

	// When
	index, err := determinePriceIndexForDate(prices, date)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if index != expectedIndex {
		t.Fatalf("expected index: %d, actual index: %d", expectedIndex, index)
	}
}

func TestDetermine_price_index_on_two_days_after(t *testing.T) {
	// Given
	date, _ := time.Parse(dateLayout, "2021-01-07")
	expectedIndex := 1

	// When
	index, err := determinePriceIndexForDate(prices, date)

	// Then
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if index != expectedIndex {
		t.Fatalf("expected index: %d, actual index: %d", expectedIndex, index)
	}
}

func TestDetermine_price_index_on_day_not_present(t *testing.T) {
	// Given
	date, _ := time.Parse(dateLayout, "2021-01-11")

	// When
	index, err := determinePriceIndexForDate(prices, date)

	// Then
	if err == nil {
		t.Fatalf("expected error, but found actual index: %d", index)
	}
}

func TestDetermine_price_index_on_day_out_of_bounds(t *testing.T) {
	// Given
	date, _ := time.Parse(dateLayout, "2021-01-02")

	// When
	index, err := determinePriceIndexForDate(prices, date)

	// Then
	if err == nil {
		t.Fatalf("expected error, but found actual index: %d", index)
	}
}

func TestDetermine_price_index_on_wrong_date_format(t *testing.T) {
	// Given
	date, _ := time.Parse(dateLayout, "BAD_FORMAT")

	// When
	index, err := determinePriceIndexForDate(prices, date)

	// Then
	if err == nil {
		t.Fatalf("expected error, but found actual index: %d", index)
	}
}

// ################# SMA screener tests #################
var screenStrategy = smaStrategy{}
var direction = above
var screenPeriod = 3
var date, _ = time.Parse(dateLayout, "2021-01-20")

var apple = companyInfo{
	symbol: "AAPL",
	historicalPrice: HistoricalPrice{
		Symbol: "AAPL",
		Historical: []Price{
			{
				Date:  "2021-01-20",
				Close: 100.0,
			},
			{
				Date:  "2021-01-19",
				Close: 95.0,
			},
			{
				Date:  "2021-01-18",
				Close: 90.0,
			},
		},
	},
}
var tesla = companyInfo{
	symbol: "TSLA",
	historicalPrice: HistoricalPrice{
		Symbol: "TSLA",
		Historical: []Price{
			{
				Date:  "2021-01-20",
				Close: 90.0,
			},
			{
				Date:  "2021-01-19",
				Close: 95.0,
			},
			{
				Date:  "2021-01-18",
				Close: 100.0,
			},
		},
	},
}
var insufficient = companyInfo{
	symbol: "INSU",
	historicalPrice: HistoricalPrice{
		Symbol: "INSU",
		Historical: []Price{
			{
				Date:  "2021-01-20",
				Close: 100.0,
			},
			{
				Date:  "2021-01-19",
				Close: 95.0,
			},
		},
	},
}
var empty = companyInfo{
	symbol: "EMTY",
	historicalPrice: HistoricalPrice{
		Symbol:     "EMTY",
		Historical: []Price{},
	},
}

func TestScreen_companies_for_sma_above_3(t *testing.T) {
	// Given
	var companies = []companyInfo{apple, tesla}
	expectedCompanies := []companyInfo{apple}

	// When
	companiesAfterScreening := screenStrategy.perform(companies, direction, screenPeriod, date)

	// Then
	if !reflect.DeepEqual(expectedCompanies, companiesAfterScreening) {
		t.Fatalf("expected companies: %+v\\n, actual companies: %+v\\n", companies, companiesAfterScreening)
	}
}

func TestScreen_companies_for_insufficient_company_data(t *testing.T) {
	// Given
	var companies = []companyInfo{apple, insufficient}
	expectedCompanies := []companyInfo{apple}

	// When
	companiesAfterScreening := screenStrategy.perform(companies, direction, screenPeriod, date)

	// Then
	if !reflect.DeepEqual(expectedCompanies, companiesAfterScreening) {
		t.Fatalf("expected companies: %+v\\n, actual companies: %+v\\n", companies, companiesAfterScreening)
	}
}

func TestScreen_companies_for_no_company_data(t *testing.T) {
	// Given
	var companies = []companyInfo{apple, empty}
	expectedCompanies := []companyInfo{apple}

	// When
	companiesAfterScreening := screenStrategy.perform(companies, direction, screenPeriod, date)

	// Then
	if !reflect.DeepEqual(expectedCompanies, companiesAfterScreening) {
		t.Fatalf("expected companies: %+v\\n, actual companies: %+v\\n", companies, companiesAfterScreening)
	}
}
