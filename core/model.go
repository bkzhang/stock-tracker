package core

type User struct {
    Name string
    Functions Functions
}

type Stock struct {
    Symbol string
    Amount uint
}

type Stocks []Stock

type Function struct {
    Name string
    Stocks Stocks
    Interval uint
}

type Functions []Function

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
