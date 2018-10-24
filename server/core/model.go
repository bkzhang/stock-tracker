package core

import "time"

type User struct {
    Name string `json:"name"`
    Stocks map[string][]OwnedStock `json:"stocks"`
}

type OwnedStock struct {
    Date time.Time `json:"date"`
    TimeZone string `json:"timezone"`
    Price float64 `json:"price"` // price bought at
    Shares int `json:"shares"` // number of shares bought
}

type StockData map[string]interface{}

type MetaData struct {
    Information string `json:"1. Information"`
    Symbol string `json:"2. Symbol"`
    LastRefreshed string `json:"3. Last Refreshed"`
    Interval string `json:"4. Interval"`
    OutputSize string `json:"5. Output Size"`
    TimeZone string `json:"6. Time Zone"`
}

type TimeSeries map[string]TimeSeriesInfo

type TimeSeriesInfo struct {
    Open float64 `json:"1. open,string"`
    High float64 `json:"2. high,string"`
    Low float64 `json:"3. low,string"`
    Close float64 `json:"4. close,string"`
    Volume uint `json:"5. volume,string"`
}

type Quote struct {
    Symbol string `json:"symbol"`
    Date time.Time `json:"date"`
    TimeZone string `json:"timezone"`
    High float64 `json:"high"`
    Low float64 `json:"low"`
    Open float64 `json:"open"`
    Close float64 `json:"close"`
    Volume int64 `json:"volume"` //mysql driver doesn't support uint64
}

type StockQuery struct {
    Quote Quote 
    Error error
}
