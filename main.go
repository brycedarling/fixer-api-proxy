package main

import (
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

func main() {
	rf, err := NewRateFetcher()
	if err != nil {
		log.Fatal(err)
	}

	c := cache.New(5*time.Minute, 10*time.Minute)

	crf := NewCachedRateFetcher(rf, c)

	pf := NewPeriodFetcher(crf)

	s := NewServer(pf)

	http.HandleFunc("/exchange-rate", s.ExchangeRateHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
