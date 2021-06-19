package main

import "time"

func main() {
	from, _ := time.Parse(dateLayout, "2021-01-01")
	to, _ := time.Parse(dateLayout, "2021-02-01")
	screenStrategy := smaStrategy{}
	screener := screener{
		direction:         above,
		periodInDays:      150,
		screeningStrategy: screenStrategy,
	}

	revGrowth12crit := criterion{
		criterionType: revenueGrowth,
		period:        PeriodAnnual,
		weight:        0.5,
		direction:     highest,
	}

	profitGrowth12crit := criterion{
		criterionType: grossProfitGrowth,
		period:        PeriodAnnual,
		weight:        0.5,
		direction:     highest,
	}

	strategy := strategy{
		portfolioSize: 3,
		criteria:      []criterion{revGrowth12crit, profitGrowth12crit},
	}

	backtest := Backtest{screener, strategy}

	backtest.doBacktest([]string{"GOOG", "AAL", "INTC", "MSFT", "NVDA", "VRTX"}, from, to)
}
