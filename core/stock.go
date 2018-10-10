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
    Data StockData
    Error error
}

func (api *Api) Function(user User, fnName string) ([]StockData, []error) {
    var length int
    ch := make(chan StockResult)
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

    stockData := make([]StockData, length)
    errs := make([]error, 0)
    for i := 0; i < length; i++ {
        res := <-ch
        if res.Error != nil {
            errs = append(errs, res.Error)
        } else {
            stockData[i] = res.Data
        }
    }

    if len(errs) > 0 {
        return stockData, errs
    }
    return stockData, nil
}

func (api *Api) TimeSeriesIntraday(stock Stock, interval uint) StockResult {
    var res StockResult
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

    if err := json.Unmarshal(body, &res.Data); err != nil {
        res.Error = fmt.Errorf("Error unmarshalling stock: %v", err)
        return res
    }

    for k, v := range res.Data {
        if k == "Information" {
            res.Data = make(StockData)
            res.Error = fmt.Errorf("Alpha Vantage API error: %v", v)
        }
    }

    for k, v := range res.Data {
        if strings.Contains(k, "Time Series") {
            // need to sort since we only want most recent
            for d := range v.(map[string]interface{}) {
                date, err := TimeSeriesToTime(d, res.Data["Meta Data"].(map[string]interface{})["6. Time Zone"].(string))
                if err != nil {
                    res.Error = err
                    return res
                }
            }
            break;
        }
    }

    return res
}

func TimeSeriesToTime(d string, timezone string) (string, error) {
    date := strings.Split(d[:10], "-")
    hourlytime := strings.Split(d[11:], ":") 
    year, err := strconv.Atoi(date[0])
    if err != nil {
        return "", fmt.Errorf("String to int conversion error: %v", err)
    }
    month, err := strconv.Atoi(date[1])
    if err != nil {
        return "", fmt.Errorf("String to int conversion error: %v", err)
    }
    day, err := strconv.Atoi(date[2])
    if err != nil {
        return "", fmt.Errorf("String to int conversion error: %v", err)
    }
    hour, err := strconv.Atoi(hourlytime[0])
    if err != nil {
        return "", fmt.Errorf("String to int conversion error: %v", err)
    }
    min, err := strconv.Atoi(hourlytime[1])
    if err != nil {
        return "", fmt.Errorf("String to int conversion error: %v", err)
    }
    sec, err := strconv.Atoi(hourlytime[2])
    if err != nil {
        return "", fmt.Errorf("String to int conversion error: %v", err)
    }
    zone, err := time.LoadLocation(timezone)
    if err != nil {
        return "", fmt.Errorf("String to *time.Location conversion error: %v", err)
    }
    t := time.Date(year, time.Month(month), day, hour, min, sec, 0, zone)
    return t.Format("2006-1-2 15:04:05"), nil
}
