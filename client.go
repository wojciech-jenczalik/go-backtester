package main

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	FmpUrl        = "https://financialmodelingprep.com/api/v3"
	FmpApiKeyKey  = "apikey"
	periodAnnual  = "annual"
	PeriodQuarter = "quarter"
)

func GetNasdaqConstituent100() ([]Company, error) {
	url := "/nasdaq_constituent"
	res, err := get(url)
	if err != nil {
		return nil, err
	}

	cmp := make([]Company, 0)
	err = json.Unmarshal(res, &cmp)
	if err != nil {
		return nil, err
	}

	return cmp, nil
}

func GetProfile(symbol string) (Profile, error) {
	url := fmt.Sprintf("/profile/%s", symbol)
	res, err := get(url)

	var profile []Profile
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(res, &profile)
	if err != nil {
		panic(err)
	}

	return profile[0], nil
}

func GetHistoricalPrices(symbol string, from time.Time, to time.Time) (HistoricalPrice, error) {
	strfrom := from.Format(dateLayout)
	strto := to.Format(dateLayout)

	url := fmt.Sprintf("/historical-price-full/%s?from=%s&to=%s", symbol, strfrom, strto)
	res, err := get(url)

	var prices HistoricalPrice
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(res, &prices)
	if err != nil {
		panic(err)
	}

	return prices, nil
}

func GetFinancialRatios(symbol string, from time.Time, to time.Time) ([]FinancialRatio, error) {
	period, limit := convertTimeToQuarters(from, to)
	url := fmt.Sprintf("/ratios/%s?period=%s&limit=%d", symbol, period, limit)
	res, err := get(url)
	if err != nil {
		return nil, err
	}
	ratios := make([]FinancialRatio, 0)
	err = json.Unmarshal(res, &ratios)
	if err != nil {
		return nil, err
	}

	return ratios, nil
}

func GetFinancialGrowthYearly(symbol string, from time.Time, to time.Time) ([]FinancialGrowth, error) {
	period, limit := convertTimeToYears(from, to)
	url := fmt.Sprintf("/financial-growth/%s?period=%s&limit=%d", symbol, period, limit)
	res, err := get(url)
	if err != nil {
		return nil, err
	}

	cmp := make([]FinancialGrowth, 0)
	err = json.Unmarshal(res, &cmp)
	if err != nil {
		return nil, err
	}

	return cmp, nil
}
