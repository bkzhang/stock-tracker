package core

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
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

    temp := strings.SplitAfter(string(body), "{\n\t\"Information\": \"Thank you for using Alpha Vantage! Please visit https://www.alphavantage.co/premium/ if you would like to have a higher API call volume.\n}")
    if len(temp) > 0 {
        body = []byte(temp[0])
    }

    if err := json.Unmarshal(body, &res.Data); err != nil {
        res.Error = fmt.Errorf("Error unmarshalling stock: %v", err)
        return res
    }

    return res
}
