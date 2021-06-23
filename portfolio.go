package main

import (
	"errors"
	"time"
)

const (
	buy  = "BUY"
	sell = "SELL"
	hold = "HOLD"
)

type portfolio struct {
	commision commision
	capital   float64
	size      int
	positions []position
}

type position struct {
	company        companyInfo
	amountOfShares int
	atPrice        float64
}

type commision struct {
	fixed    float64
	perShare float64
}

type signal struct {
	date           time.Time
	company        companyInfo
	price          float64
	amountOfShares int
	action         string
}

func (p *portfolio) generateSignals(newPositions []position, date time.Time) []signal {
	signals := make([]signal, 0)

	for _, position := range newPositions {
		contains, atIndex := contains(p.positions, position)
		var sig signal
		newSharesAmount := position.amountOfShares
		if contains {
			sharesCurrentlyHeldAmount := p.positions[atIndex].amountOfShares
			if sharesCurrentlyHeldAmount > newSharesAmount {
				sig = signal{
					date:           date,
					company:        position.company,
					price:          position.atPrice,
					amountOfShares: sharesCurrentlyHeldAmount - newSharesAmount,
					action:         sell,
				}
			} else if sharesCurrentlyHeldAmount == newSharesAmount {
				sig = signal{
					date:           date,
					company:        position.company,
					price:          position.atPrice,
					amountOfShares: sharesCurrentlyHeldAmount,
					action:         hold,
				}
			} else if sharesCurrentlyHeldAmount < newSharesAmount {
				sig = signal{
					date:           date,
					company:        position.company,
					price:          position.atPrice,
					amountOfShares: newSharesAmount - sharesCurrentlyHeldAmount,
					action:         buy,
				}
			}
		} else {
			sig = signal{
				date:           date,
				company:        position.company,
				price:          position.atPrice,
				amountOfShares: newSharesAmount,
				action:         buy,
			}
		}
		signals = append(signals, sig)
	}

	for _, position := range p.positions {
		contains, _ := contains(newPositions, position)

		if !contains {
			sig := signal{
				date:           date,
				company:        position.company,
				price:          position.atPrice,
				amountOfShares: position.amountOfShares,
				action:         sell,
			}
			signals = append(signals, sig)
		}
	}

	return signals
}

func (p *portfolio) patchPortfolio(signals []signal) error {
	for _, signal := range signals {
		p.performSignalAction(signal)
	}

	return nil
}

func (p *portfolio) performSignalAction(signal signal) {
	switch signal.action {
	case sell:
		p.capital += float64(signal.amountOfShares)*(signal.price-p.commision.perShare) - p.commision.fixed
		indexAt := indexAt(p.positions, signal.company.symbol)
		if p.positions[indexAt].amountOfShares == signal.amountOfShares {
			p.positions = remove(p.positions, indexAt)
		} else {
			p.positions[indexAt] = position{
				company:        signal.company,
				amountOfShares: p.positions[indexAt].amountOfShares - signal.amountOfShares,
				atPrice:        signal.price,
			}
		}
	case buy:
		p.capital -= float64(signal.amountOfShares)*(signal.price+p.commision.perShare) + p.commision.fixed
		containsSymbol, indexAt := containsSymbol(p.positions, signal.company.symbol)
		if containsSymbol {
			p.positions[indexAt] = position{
				company:        signal.company,
				amountOfShares: p.positions[indexAt].amountOfShares + signal.amountOfShares,
				atPrice:        signal.price,
			}
		} else {
			p.positions = append(p.positions, position{
				company:        signal.company,
				amountOfShares: signal.amountOfShares,
				atPrice:        signal.price,
			})
		}
	}
}

func indexAt(positions []position, symbol string) int {
	for index, p := range positions {
		if p.company.symbol == symbol {
			return index
		}
	}
	return -1
}

func remove(s []position, index int) []position {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

func (p *portfolio) calculateNewPositions(topCompanies []companyInfo, date time.Time) ([]position, error) {
	portfolioValue, err := p.calculatePortfolioValue(date)
	if err != nil {
		return nil, err
	}
	positions := make([]position, 0)
	for _, topCompany := range topCompanies {
		amountOfShares, price := p.calculateAmountAndPriceOfShares(topCompany, portfolioValue, date)
		positions = append(positions, position{
			company:        topCompany,
			amountOfShares: amountOfShares,
			atPrice:        price,
		})
	}

	return positions, nil
}

func contains(positions []position, position position) (bool, int) {
	for index, p := range positions {
		if p.company.symbol == position.company.symbol {
			return true, index
		}
	}

	return false, 0
}

func containsSymbol(positions []position, symbol string) (bool, int) {
	for index, p := range positions {
		if p.company.symbol == symbol {
			return true, index
		}
	}

	return false, 0
}

var portfolioCalculationError = errors.New("error while evaluating portfolio value")

func (p *portfolio) calculatePortfolioValue(date time.Time) (float64, error) {
	var positionsValue float64

	for _, position := range p.positions {
		priceIndex, err := determinePriceIndexForDate(position.company.historicalPrice.Historical, date)
		if err != nil {
			return 0, portfolioCalculationError
		}
		positionsValue += position.company.historicalPrice.Historical[priceIndex].Close * float64(position.amountOfShares)
	}

	return positionsValue + p.capital, nil
}

func (p *portfolio) calculateAmountAndPriceOfShares(company companyInfo, portfolioValue float64, date time.Time) (int, float64) {
	valueGrantedPerCompany := portfolioValue / float64(p.size)
	priceIndex, _ := determinePriceIndexForDate(company.historicalPrice.Historical, date)
	price := company.historicalPrice.Historical[priceIndex].Close

	return int(valueGrantedPerCompany / price), price
}
