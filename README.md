Set FIXER_API_KEY environment variable to your Fixer API key first before running:

```
export FIXER_API_KEY=changeme
```

Build program:

```
go build
```

Run server:

```
./fixer-api-proxy
```

Request exchange rates:

```
curl localhost:8000/exchange-rate?start=2013-12-23&end=2013-12-26&other=GBP
```

The query params for the `/exchange-rate` endpoint are:

- `start` - period start date
- `end` - period end date
- `base` - base currency for conversion rates
- `other` - other currency for conversion rates


The `/exchange-rate` endpoint will request the Fixer API for every date in that period and return all of the exchange rates between the two given currencies:

```
{"2013-12-23":0.837729,"2013-12-24":0.835788,"2013-12-25":0.836701,"2013-12-26":0.834049}
```

---

In order to run tests:

```
go test
```

However, due to time constraints there is only one test for the API using a fake period fetcher.

---

*Oops*, while doing the 2nd part of the challenge, I noticed that the Fixer API has an `Time-Series Endpoint` instead of the `Historical Rates Endpoint` that I used. This would've saved me from having to make N API requests.
