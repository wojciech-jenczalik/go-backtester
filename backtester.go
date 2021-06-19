package main

import (
	"fmt"
	"log"
	"time"
)

type Backtest struct {
	screener
	strategy
	//portfolio
}

func (b *Backtest) doBacktest(symbols []string, from time.Time, to time.Time) {
	companies := prepareData(symbols, from, to, b.screener.periodInDays)

	currentBacktestDate := from

	for currentBacktestDate.Before(to) {
		screenedCompanies := b.screener.screen(companies, currentBacktestDate)
		topCompanies := b.strategy.evaluateTopCompanies(screenedCompanies, currentBacktestDate)
		fmt.Printf("%+v\n", topCompanies)
		break
	}
}

func prepareData(symbols []string, from time.Time, to time.Time, screeningPeriod int) []companyInfo {
	// The NYSE and NASDAQ average about 253 trading days a year.
	// This is from 365.25 (days on average per year) * 5/7 (proportion work days per week)
	// - 6 (weekday holidays) - 3*5/7 (fixed Date holidays) = 252.75 ≈ 253.
	tradingDaysInYearRatio := 1.44
	safeOffset := 10.0

	screeningPeriod = int(float64(screeningPeriod)*tradingDaysInYearRatio + safeOffset)
	from = from.AddDate(0, 0, -screeningPeriod)

	return gatherInfo(symbols, from, to)
}

func gatherInfo(symbols []string, from time.Time, to time.Time) []companyInfo {
	cmps := make([]companyInfo, len(symbols))

	for i, tckr := range symbols {
		log.Printf("Gathering information for %s... \n", tckr)
		histPrice, err := GetHistoricalPrices(tckr, from, to)
		if err != nil {
			panic(err)
		}
		//profile, err := GetProfile(tckr)
		//if err != nil {
		//	panic(err)
		//}
		//finRatios, err := GetFinancialRatios(tckr, from, to)
		//if err != nil {
		//	panic(err)
		//}
		finGrowth, err := GetFinancialGrowthYearly(tckr, from, to)
		if err != nil {
			panic(err)
		}
		cmps[i] = companyInfo{
			tckr,
			Profile{},
			histPrice,
			nil,
			finGrowth,
		}
		log.Printf("Gathering information for %s completed.\n", tckr)
	}

	return cmps
}
