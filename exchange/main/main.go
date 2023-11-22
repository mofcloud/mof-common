package main

import (
	"encoding/json"
	"fmt"
	ex "github.com/mofcloud/mof-common/exchange"
	ptime "github.com/mofcloud/mof-common/time"
	curLib "golang.org/x/text/currency"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	exRateHisApi = "https://v6.exchangerate-api.com/v6/%s/history/USD/%d/%d/%d"
)

func main() {
	res := &ex.Exchange{
		BaseCode: curLib.USD.String(),
		ByMonth:  []*ex.ExchangeByMonth{},
	}

	// 10 year data
	// start date 2000-01-01
	// end data currency
	token := "ac363279c6241e38f74085d0"

	curr, _ := ptime.StringToTime("2010-01-31")
	now := time.Now()
	for curr.Before(now) {
		// call API
		urlStr := fmt.Sprintf(exRateHisApi, token, curr.Year(), curr.Month(), curr.Day())
		rawResp, err := http.Get(urlStr)

		if err != nil {
			panic(err)
		}

		bytes, err := io.ReadAll(rawResp.Body)
		if err != nil {
			panic(err)
		}

		innerType := struct {
			Result          string             `json:"result"`
			BaseCode        string             `json:"base_code"`
			ConversionRates map[string]float64 `json:"conversion_rates"`
		}{}

		err = json.Unmarshal(bytes, &innerType)
		if err != nil {
			panic(err)
		}

		// add to res
		byMonth := res.GetByMonth(ptime.TimeToLayoutMonth(curr))
		for code, rate := range innerType.ConversionRates {
			switch code {
			case "CNY", "USD", "EUR":
				byMonth.AddRate(code, rate)
			}
		}

		curr = nextTime(curr)
	}

	res.Finish()

	// marshal
	bytes, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("exchange.json")
	if err != nil {
		panic(err)
	}

	if _, err := f.Write(bytes); err != nil {
		panic(err)
	}
}

func nextTime(currData time.Time) time.Time {
	next := ptime.NextMonthLayoutMonthTime(currData)
	res := ptime.LastDayOfMonthTime(next)
	return res
}
