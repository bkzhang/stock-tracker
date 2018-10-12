package core

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type Api struct {
    Key string
}

type StockResult struct {
    Symbol string
    Date time.Time
    TimeZone string
    High float64
    Low float64
    Open float64
    Close float64
    Volume uint64
}

type Result struct {
    StockResult StockResult
    Error error
}

func (api *Api) Function(user User, fnName string) ([]StockResult, []error) {
    var length int
    ch := make(chan Result)
    for _, fn := range user.Functions {
        if fn.Name == "time_series_intraday" {
            length = len(fn.Stocks)
            interval := fn.Interval
            for _, stock := range fn.Stocks {
                go func(s Stock) {
                    ch <- api.TimeSeriesIntraday(s, interval) 
                }(stock)
            }
            break;
        }
    }

    stockResults := make([]StockResult, length)
    errs := make([]error, 0)
    for i := 0; i < length; i++ {
        res := <-ch
        stockResults[i] = res.StockResult
        if res.Error != nil {
            errs = append(errs, res.Error)
        }
    }

    return stockResults, errs
}

func (api *Api) TimeSeriesIntraday(stock Stock, interval uint) Result {
    var res Result
    var stockData StockData
    if stock.Symbol == "" {
        res.Error = fmt.Errorf("Stock object missing symbol")
        return res
    }
    
    uri := "function=TIME_SERIES_INTRADAY&symbol="+stock.Symbol+"&interval="+strconv.FormatUint(uint64(interval), 10)+"min&apikey="+api.Key

    resp, err := http.Get(URL+uri)
    if err != nil {
        res.Error = fmt.Errorf("Could not get function results from Alpha Vantage: %v", err)
        return res
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        res.Error = fmt.Errorf("Error getting the response body: %v", err)
        return res
    }

    if err := json.Unmarshal(body, &stockData); err != nil {
        res.Error = fmt.Errorf("Error unmarshalling stock: %v", err)
        return res
    }

    for k, v := range stockData {
        if k == "Information" {
            res.StockResult.Symbol = v.(map[string]interface{})["Meta Data"].(map[string]interface{})["Symbol"].(string)
            res.Error = fmt.Errorf("Alpha Vantage API error: %v", v)
        }
    }

    metadata := stockData["Meta Data"].(map[string]interface{}) 
    res.StockResult.Symbol = metadata["2. Symbol"].(string)
    for k, v := range stockData {
        if strings.Contains(k, "Time Series") {
            for d, v2 := range v.(map[string]interface{}) {
                date := metadata["3. Last Refreshed"].(string)
                if d == date {
                    timezone := metadata["6. Time Zone"].(string)
                    stockDate, err := TimeSeriesToTime(d, timezone)
                    if err != nil {
                        res.Error = err
                        return res
                    }

                    v3 := v2.(map[string]interface{})

                    res.StockResult.Date = stockDate
                    res.StockResult.TimeZone = timezone //t.Format("2006-1-2 15:04:05") 

                    res.StockResult.Open, err = strconv.ParseFloat(v3["1. open"].(string), 64)
                    if err != nil {
                        res.Error = err
                        return res
                    }

                    res.StockResult.High, err = strconv.ParseFloat(v3["2. high"].(string), 64)
                    if err != nil {
                        res.Error = err
                        return res
                    }

                    res.StockResult.Low, err = strconv.ParseFloat(v3["3. low"].(string), 64)
                    if err != nil {
                        res.Error = err
                        return res
                    }

                    res.StockResult.Close, err = strconv.ParseFloat(v3["4. close"].(string), 64)
                    if err != nil {
                        res.Error = err
                        return res
                    }

                    res.StockResult.Volume, err = strconv.ParseUint(v3["5. volume"].(string), 10, 64)
                    if err != nil {
                        res.Error = err
                        return res
                    }

                    break;
                }
            }
            break;
        }
    }

    return res
}

func TimeSeriesToTime(d string, timezone string) (time.Time, error) {
    date := strings.Split(d[:10], "-")
    hourlytime := strings.Split(d[11:], ":") 
    year, err := strconv.Atoi(date[0])
    if err != nil {
        return time.Time{}, fmt.Errorf("String to int conversion error: %v", err)
    }
    month, err := strconv.Atoi(date[1])
    if err != nil {
        return time.Time{}, fmt.Errorf("String to int conversion error: %v", err)
    }
    day, err := strconv.Atoi(date[2])
    if err != nil {
        return time.Time{}, fmt.Errorf("String to int conversion error: %v", err)
    }
    hour, err := strconv.Atoi(hourlytime[0])
    if err != nil {
        return time.Time{}, fmt.Errorf("String to int conversion error: %v", err)
    }
    min, err := strconv.Atoi(hourlytime[1])
    if err != nil {
        return time.Time{}, fmt.Errorf("String to int conversion error: %v", err)
    }
    sec, err := strconv.Atoi(hourlytime[2])
    if err != nil {
        return time.Time{}, fmt.Errorf("String to int conversion error: %v", err)
    }
    zone, err := time.LoadLocation(timezone)
    if err != nil {
        return time.Time{}, fmt.Errorf("String to *time.Location conversion error: %v", err)
    }
    t := time.Date(year, time.Month(month), day, hour, min, sec, 0, zone)
    return t, nil
}
