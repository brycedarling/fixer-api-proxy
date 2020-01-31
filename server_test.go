package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakePeriodFetcher struct{}

func (fakePeriodFetcher) FetchExchangeRates(startDateStr, endDateStr, base, other string) (map[string]float64, error) {
	rates := make(map[string]float64)
	rates["2013-12-23"] = 0.1
	rates["2013-12-24"] = 0.2
	rates["2013-12-25"] = 0.3
	return rates, nil
}

func TestExchangeRateHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/exchange-rate?start=2013-12-23&end=2013-12-25&other=USD", nil)
	if err != nil {
		t.Fatal(err)
	}

	fake := fakePeriodFetcher{}
	s := NewServer(fake)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.ExchangeRateHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"2013-12-23":0.1,"2013-12-24":0.2,"2013-12-25":0.3}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
