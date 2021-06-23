package main

import (
	"time"
)

type Backtest struct {
	screener
	strategy
	portfolio
}

func (b *Backtest) doBacktest(symbols []string, from time.Time, to time.Time, iterateForDays int) {
	companies := prepareData(symbols, from, to, b.screener.periodInDays)

	currentBacktestDate := from

	for currentBacktestDate.Before(to) {
		screenedCompanies := b.screener.screen(companies, currentBacktestDate)
		topCompanies := b.strategy.evaluateTopCompanies(screenedCompanies, currentBacktestDate, b.portfolio.size)
		newPositions, err := b.portfolio.calculateNewPositions(topCompanies, currentBacktestDate)
		if err != nil {
			currentBacktestDate = currentBacktestDate.AddDate(0, 0, iterateForDays)
			continue
		}
		signals := b.portfolio.generateSignals(newPositions, currentBacktestDate)
		err = b.portfolio.patchPortfolio(signals)
		if err != nil {
			currentBacktestDate = currentBacktestDate.AddDate(0, 0, iterateForDays)
			continue
		}
		currentBacktestDate = currentBacktestDate.AddDate(0, 0, iterateForDays)
		//log.Println(b.portfolio.calculatePortfolioValue(to))
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
	}

	return cmps
}
