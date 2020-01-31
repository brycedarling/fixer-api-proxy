package fixer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type RateFetcher interface {
	FetchExchangeRate(string, string, string) (float64, error)
}

type rateFetcher struct {
	APIKey string
}

func NewRateFetcher() (RateFetcher, error) {
	apiKey := os.Getenv("FIXER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("missing FIXER_API_KEY env var")
	}

	return rateFetcher{apiKey}, nil
}

type exchangeRateResponse struct {
	Success    bool               `json:"success"`
	Timestamp  int                `json:"timestamp"`
	Historical bool               `json:"historical"`
	Base       string             `json:"base"`
	Date       string             `json:"date"`
	Rates      map[string]float64 `json:"rates"`
}

func (f rateFetcher) FetchExchangeRate(date, base, other string) (float64, error) {
	url := fmt.Sprintf("http://data.fixer.io/api/%s?access_key=%s&base=%s&other=%s", date, f.APIKey, base, other)
	fmt.Println(url)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	data := exchangeRateResponse{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	if amount, ok := data.Rates[other]; ok {
		return amount, nil
	}

	return 0, fmt.Errorf("rate not found for other %s", other)
}
