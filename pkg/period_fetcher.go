package fixer

import (
	"fmt"
	"time"
)

type PeriodFetcher interface {
	FetchExchangeRates(startDateStr, endDateStr, base, other string) (map[string]float64, error)
}

type periodFetcher struct {
	rateFetcher RateFetcher
}

func NewPeriodFetcher(rf RateFetcher) PeriodFetcher {
	return periodFetcher{rf}
}

func (pf periodFetcher) FetchExchangeRates(startDateStr, endDateStr, base, other string) (map[string]float64, error) {
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return nil, err
	}

	if endDate.Before(startDate) {
		return nil, fmt.Errorf("end date %s is before start date %s", endDate, startDate)
	}

	rates := make(map[string]float64)
	for date := rangeDate(startDate, endDate); ; {
		d := date()

		if d.IsZero() {
			break
		}

		formattedDate := d.Format(layout)

		rate, err := pf.rateFetcher.FetchExchangeRate(formattedDate, base, other)
		if err != nil {
			return nil, err
		}

		rates[formattedDate] = rate
	}

	return rates, nil
}

func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}
