package main

type companyInfo struct {
	symbol          string
	profile         Profile
	historicalPrice HistoricalPrice
	ratios          []FinancialRatio
	growth          []FinancialGrowth
}

type Company struct {
	Symbol         string
	Name           string
	DateFirstAdded string
}

type Profile struct {
	IpoDate     string
	CompanyName string
}

type FinancialGrowth struct {
	Symbol            string
	Date              string
	RevenueGrowth     float64
	GrossProfitGrowth float64
	NetIncomeGrowth   float64
}

type FinancialRatio struct {
	Symbol string
	Date   string
	Period string
}

type HistoricalPrice struct {
	Symbol     string
	Historical []Price
}

type Price struct {
	Date  string
	Open  float64
	Close float64
	Low   float64
	High  float64
}
