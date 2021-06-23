package main

import (
	"errors"
	"log"
	"time"
)

// Screen direction
const (
	below = "below"
	above = "above"
)

type screeningStrategy interface {
	perform(cmpInfs []companyInfo, direction string, periodDays int, date time.Time) []companyInfo
}

type screener struct {
	direction    string
	periodInDays int
	screeningStrategy
}

func (s screener) screen(companyInfos []companyInfo, date time.Time) []companyInfo {
	return s.screeningStrategy.perform(companyInfos, s.direction, s.periodInDays, date)
}

type smaStrategy struct{}

func (s smaStrategy) perform(companyInfos []companyInfo, direction string, screeningPeriodDays int, date time.Time) []companyInfo {
	result := make([]companyInfo, 0)

	for _, company := range companyInfos {
		priceHistory := company.historicalPrice.Historical
		startingIndex, err := determinePriceIndexForDate(priceHistory, date)

		// Skip company the Price index of which could not be determined
		// (e.g. IPO was later than we try to calculate SMA from)
		if err != nil {
			log.Printf("error while screening %s: %s \n", company.symbol, err)
			continue
		}
		if startingIndex+screeningPeriodDays > len(priceHistory) {
			log.Printf("error while screening %s: screening period out of bounds", company.symbol)
			continue
		}

		priceHistoryForScreening := make([]float64, 0)

		for _, price := range priceHistory[startingIndex : startingIndex+screeningPeriodDays] {
			priceHistoryForScreening = append(priceHistoryForScreening, price.Close)
		}

		sma := Sma(priceHistoryForScreening...)

		shouldAppend :=
			direction == above && priceHistoryForScreening[0] > sma ||
				direction == below && priceHistoryForScreening[0] < sma

		if shouldAppend {
			result = append(result, company)
		}
	}

	return result
}

var dateIndexNotFound = errors.New("Date index could not be determined")
