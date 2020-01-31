package main

import (
	"log"
	"net/http"
	"time"

	fixer "github.com/brycedarling/fixer-api-proxy/pkg"
	"github.com/patrickmn/go-cache"
)

func main() {
	rf, err := fixer.NewRateFetcher()
	if err != nil {
		log.Fatal(err)
	}

	c := cache.New(5*time.Minute, 10*time.Minute)

	crf := fixer.NewCachedRateFetcher(rf, c)

	pf := fixer.NewPeriodFetcher(crf)

	s := NewServer(pf)

	http.HandleFunc("/exchange-rate", s.ExchangeRateHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
