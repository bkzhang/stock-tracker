package model 

import "time"

type Stock struct {
    Symbol string
    Date time.Time
    TimeZone string
    High float64
    Low float64
    Open float64
    Close float64
    Volume uint64
}

type Stocks []Stock
