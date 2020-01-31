package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type Server interface {
	ExchangeRateHandler(http.ResponseWriter, *http.Request)
}

type server struct {
	periodFetcher PeriodFetcher
}

func NewServer(pf PeriodFetcher) Server {
	return &server{pf}
}

func (s server) ExchangeRateHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	start := q.Get("start")
	end := q.Get("end")
	base := q.Get("base")
	other := q.Get("other")

	if start == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "start date required")
		return
	}
	if end == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "end date required")
		return
	}
	if base == "" {
		base = "EUR"
	}
	if other == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "other currency required")
		return
	}

	rates, err := s.periodFetcher.FetchExchangeRates(start, end, base, other)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	jsonString, err := json.Marshal(rates)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}
