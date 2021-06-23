package main

import (
	"errors"
	"log"
	"math"
	"sort"
	"time"
)

type strategy struct {
	criteria []criterion
}

type criterion struct {
	criterionType string
	period        string
	weight        float64
	direction     string
}

type criteriaEvaluationResult struct {
	companySymbol string
	results       []float64
	error         error
}

const (
	// Criteria
	revenueGrowth     = "REVENUE_GROWTH"
	grossProfitGrowth = "GROSS_PROFIT_GROWTH"

	// Direction
	lowest  = "LOWEST"
	highest = "HIGHEST"
)

func (s *strategy) evaluateTopCompanies(companies []companyInfo, date time.Time, portfolioSize int) []companyInfo {
	evaluationResult := s.evaluateCriteria(companies, date)
	evaluationResult = filterOutErrorResults(evaluationResult)
	evaluationResult = s.normalizeResultsTo01(evaluationResult)
	finalResults := s.calculateFinalResults(evaluationResult)
	return s.selectTopCompanies(finalResults, companies, portfolioSize)
}

var criteriaUnknown = errors.New("unknown criteria for company evaluation were provided")

func (s *strategy) evaluateCriteria(companies []companyInfo, date time.Time) []criteriaEvaluationResult {
	criteriaEvaluationResults := make([]criteriaEvaluationResult, len(companies))

	for i, company := range companies {
		results := make([]float64, len(s.criteria))

		for j, criterion := range s.criteria {
			switch criterion.criterionType {
			case revenueGrowth:
				result, err := getRevenueGrowth(company, criterion.period, date)
				if err != nil {
					criteriaEvaluationResults[i] = criteriaEvaluationResult{error: err}
					break
				}
				results[j] = result
			case grossProfitGrowth:
				result, err := getGrossProfitGrowth(company, criterion.period, date)
				if err != nil {
					criteriaEvaluationResults[i] = criteriaEvaluationResult{error: err}
					break
				}
				results[j] = result
			default:
				criteriaEvaluationResults[i] = criteriaEvaluationResult{error: criteriaUnknown}
			}
		}

		criteriaEvaluationResults[i] = criteriaEvaluationResult{
			companySymbol: company.symbol,
			results:       results,
		}
	}

	return criteriaEvaluationResults
}

func filterOutErrorResults(results []criteriaEvaluationResult) []criteriaEvaluationResult {
	filtered := make([]criteriaEvaluationResult, 0)

	for _, result := range results {
		if result.error == nil {
			filtered = append(filtered, result)
		}
	}

	return filtered
}

func (s *strategy) normalizeResultsTo01(companyResults []criteriaEvaluationResult) []criteriaEvaluationResult {
	if len(companyResults) == 0 {
		return make([]criteriaEvaluationResult, 0)
	}

	minResultsOfGivenCriteria := make([]float64, 0)
	maxResultsOfGivenCriteria := make([]float64, 0)

	for _, companyResult := range companyResults {
		for i, resultOfGivenCriterion := range companyResult.results {
			if len(minResultsOfGivenCriteria) == i {
				minResultsOfGivenCriteria = append(minResultsOfGivenCriteria, resultOfGivenCriterion)
			} else if resultOfGivenCriterion < minResultsOfGivenCriteria[i] {
				minResultsOfGivenCriteria[i] = resultOfGivenCriterion
			}
			if len(maxResultsOfGivenCriteria) == i {
				maxResultsOfGivenCriteria = append(maxResultsOfGivenCriteria, resultOfGivenCriterion)
			} else if resultOfGivenCriterion > maxResultsOfGivenCriteria[i] {
				maxResultsOfGivenCriteria[i] = resultOfGivenCriterion
			}
		}
	}

	for i, companyResult := range companyResults {
		for j, criterionResult := range companyResult.results {
			normalizedValue, err := normalizeValueTo01(
				criterionResult,
				minResultsOfGivenCriteria[j],
				maxResultsOfGivenCriteria[j],
				s.criteria[j].direction)
			if err != nil {
				companyResults[i].results[j] = 0
			}
			companyResults[i].results[j] = normalizedValue
		}
	}

	return companyResults
}

func (s *strategy) calculateFinalResults(results []criteriaEvaluationResult) map[string]float64 {
	var finalEvaluationResults map[string]float64
	finalEvaluationResults = make(map[string]float64, 0)

	for _, result := range results {
		finalResult := 0.0
		for j, singleResult := range result.results {
			finalResult += singleResult * s.criteria[j].weight
		}
		finalEvaluationResults[result.companySymbol] = finalResult
	}

	return finalEvaluationResults
}

func (s *strategy) selectTopCompanies(result map[string]float64, companies []companyInfo, portfSize int) []companyInfo {
	symbols := make([]string, 0, len(result))
	for symbol := range result {
		symbols = append(symbols, symbol)
	}

	sort.Slice(symbols, func(i, j int) bool {
		return result[symbols[i]] > result[symbols[j]]
	})

	amount := int(math.Min(float64(portfSize), float64(len(result))))
	topCompaniesSymbols := symbols[0:amount]

	topCompanies := make([]companyInfo, 0)

	for _, symbol := range topCompaniesSymbols {
		topCompany, err := findBySymbol(companies, symbol)
		if err != nil {
			log.Println(err)
			continue
		}
		topCompanies = append(topCompanies, topCompany)
	}

	return topCompanies
}

var companyNotFound = errors.New("could not find company by symbol")

func findBySymbol(companies []companyInfo, symbol string) (companyInfo, error) {
	for _, company := range companies {
		if company.symbol == symbol {
			return company, nil
		}
	}
	log.Println(companyNotFound, " "+symbol)
	return companyInfo{}, companyNotFound
}

func getRevenueGrowth(company companyInfo, period string, date time.Time) (float64, error) {
	growthReport, err := getGrowthReport(company, period, date)
	if err != nil {
		return 0.0, err
	}
	return growthReport.RevenueGrowth, nil
}

func getGrossProfitGrowth(company companyInfo, period string, date time.Time) (float64, error) {
	growthReport, err := getGrowthReport(company, period, date)
	if err != nil {
		return 0.0, err
	}
	return growthReport.GrossProfitGrowth, nil
}

var periodNotSupported = errors.New("period not supported. Supported periods are: periodAnnual")
var companyGrowthNotFound = errors.New("could not find company growth report for given Date")

func getGrowthReport(company companyInfo, period string, date time.Time) (FinancialGrowth, error) {
	if period != periodAnnual {
		return FinancialGrowth{}, periodNotSupported
	}
	for _, growthReport := range company.growth {
		reportDate, err := time.Parse(dateLayout, growthReport.Date)
		if err != nil {
			return FinancialGrowth{}, timeParseError
		}
		reportDatePlus1Year := reportDate.AddDate(1, 0, 0)

		if reportDate.Before(date) && reportDatePlus1Year.After(date) {
			return growthReport, nil
		}
	}

	return FinancialGrowth{}, companyGrowthNotFound
}

var unsupportedDirection = errors.New("unknown direction type")

func normalizeValueTo01(val float64, min float64, max float64, direction string) (float64, error) {
	normalizedTo01 := (val - min) / (max - min)
	if direction == highest {
		return normalizedTo01, nil
	} else if direction == lowest {
		return 1 - normalizedTo01, nil
	}
	log.Println(unsupportedDirection, " "+direction)
	return 0, unsupportedDirection
}
