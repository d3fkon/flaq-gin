package conversions

import "errors"

// Ticker struct
type Ticker struct {
	TickerName string  `json:"TickerName"`
	Conversion float64 `json:"Conversion"`
	Name       string  `json:"Name"`
}

// Get all conversions
// TODO make this a microservice call, to a service which uses a 3rd party api and redis
func GetTickers() []Ticker {
	tickers := []Ticker{
		{
			TickerName: "USDT",
			Name:       "USD Tether",
			Conversion: 80.0,
		},
		{
			TickerName: "BNB",
			Name:       "Binance Token",
			Conversion: 19864.0,
		},
		{
			TickerName: "MATIC",
			Name:       "Polygon",
			Conversion: 59.23,
		},
	}
	return tickers
}

// Get a ticker by a given name
// Iterate over all the available tickers and find a ticker by name
// TODO: Use hashmaps to find tickers instead of a linear search
func GetTickerByName(name string) (Ticker, error) {
	tickers := GetTickers()
	for _, t := range tickers {
		if t.TickerName == name {
			return t, nil
		}
	}
	return Ticker{}, errors.New("cannot find ticker with the given name")
}
