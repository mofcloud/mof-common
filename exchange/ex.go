package pex

import (
	_ "embed"
	"encoding/json"
	"sort"
)

//go:embed codeToFlag.json
var codeToFlagRaw []byte

var CodeToFlagByCode map[string]*CodeToFlag

func init() {
	flagList := make([]*CodeToFlag, 0)
	if err := json.Unmarshal(codeToFlagRaw, &flagList); err != nil {
		panic(err)
	}

	CodeToFlagByCode = make(map[string]*CodeToFlag)

	for i := range flagList {
		CodeToFlagByCode[flagList[i].Code] = flagList[i]
	}
}

type CodeToFlag struct {
	Code        string `json:"code"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Flag        string `json:"flag"`
}

type Exchange struct {
	BaseCode string             `json:"baseCode"`
	ByMonth  []*ExchangeByMonth `json:"byMonth"`
}

func (e *Exchange) Copy() *Exchange {
	res := &Exchange{
		BaseCode: e.BaseCode,
	}

	for i := range e.ByMonth {
		res.ByMonth = append(res.ByMonth, e.ByMonth[i].Copy())
	}

	return res
}

func (e *Exchange) GetByMonth(month string) *ExchangeByMonth {
	for i := range e.ByMonth {
		if e.ByMonth[i].Month == month {
			return e.ByMonth[i]
		}
	}

	res := NewExchangeByMonth(month)
	e.ByMonth = append(e.ByMonth, res)

	return res
}

func (e *Exchange) Finish() {
	sort.Slice(e.ByMonth, func(i, j int) bool {
		return e.ByMonth[i].Month > e.ByMonth[j].Month
	})

	for i := range e.ByMonth {
		m := e.ByMonth[i]
		sort.Slice(m.RateList, func(i, j int) bool {
			return m.RateList[i].CurrencyCode > m.RateList[j].CurrencyCode
		})
	}
}

func (e *Exchange) Convert(ts string, origin float64, from, to string) (float64, bool) {

}

func NewExchangeByMonth(month string) *ExchangeByMonth {
	return &ExchangeByMonth{
		Month:    month,
		RateList: []*Rate{},
	}
}

type ExchangeByMonth struct {
	Month    string  `json:"month"`
	RateList []*Rate `json:"rateList"`
}

func (e *ExchangeByMonth) Copy() *ExchangeByMonth {
	res := &ExchangeByMonth{
		Month: e.Month,
	}

	for i := range e.RateList {
		res.RateList = append(res.RateList, e.RateList[i].Copy())
	}

	return res
}

func (e *ExchangeByMonth) AddRate(code string, rate float64) {
	r := &Rate{
		CurrencyCode: code,
		Rate:         rate,
	}
	e.RateList = append(e.RateList, r)

	// get flag
	if v, ok := CodeToFlagByCode[code]; ok {
		r.CountryName = v.Country
		r.CountryCode = v.CountryCode
	}
}

type Rate struct {
	CurrencyCode string  `json:"currencyCode"`
	CountryName  string  `json:"countryName"`
	CountryCode  string  `json:"countryCode"`
	Rate         float64 `json:"rate"`
}

func (r *Rate) Copy() *Rate {
	return &Rate{
		CurrencyCode: r.CurrencyCode,
		CountryName:  r.CountryName,
		CountryCode:  r.CountryCode,
		Rate:         r.Rate,
	}
}
