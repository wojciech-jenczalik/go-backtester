package main

import "time"

func main() {
	from, _ := time.Parse(dateLayout, "2017-01-01")
	to, _ := time.Parse(dateLayout, "2021-06-21")
	screenStrategy := smaStrategy{}
	screener := screener{
		direction:         above,
		periodInDays:      150,
		screeningStrategy: screenStrategy,
	}

	revGrowth12crit := criterion{
		criterionType: revenueGrowth,
		period:        periodAnnual,
		weight:        0.5,
		direction:     highest,
	}

	profitGrowth12crit := criterion{
		criterionType: grossProfitGrowth,
		period:        periodAnnual,
		weight:        0.5,
		direction:     highest,
	}

	strategy := strategy{
		criteria: []criterion{revGrowth12crit, profitGrowth12crit},
	}

	// Degiro commision for US stocks
	commision := commision{
		fixed:    0.5,
		perShare: 0.0034,
	}

	portfolio := portfolio{
		commision: commision,
		capital:   10000,
		size:      3,
		positions: make([]position, 0),
	}

	backtest := Backtest{screener, strategy, portfolio}

	backtest.doBacktest([]string{"GOOG", "AAL", "INTC", "MSFT", "NVDA", "VRTX"}, from, to, 30)
}
