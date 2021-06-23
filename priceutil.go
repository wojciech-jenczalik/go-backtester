package main

import "time"

func determinePriceIndexForDate(priceHistory []Price, date time.Time) (int, error) {
	for index, price := range priceHistory {
		priceDateFormatted, err := time.Parse(dateLayout, price.Date)
		if err != nil {
			break
		}
		if priceDateFormatted == date {
			return index, nil
		}
		if index == len(priceHistory)-1 {
			break
		}
		nextPriceDateFormatted, err := time.Parse(dateLayout, priceHistory[index+1].Date)
		if err != nil {
			break
		}
		if date.Before(priceDateFormatted) && date.After(nextPriceDateFormatted) {
			return index, nil
		}
	}
	return -1, dateIndexNotFound
}
