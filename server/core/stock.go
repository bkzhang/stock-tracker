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

const URL = "https://www.alphavantage.co/query?"

type Api struct {
    Key string
}

func (api *Api) GetQuote(symbol string) (Quote, error) {
    quote := &Quote{Symbol: symbol}
    res := api.QueryIntraday(quote, 1)
    if res.Error != nil {
        return res.Quote, res.Error 
    }

    return res.Quote, nil
}

func (api *Api) IntraDayOneMin(user User) ([]Quote, []error) {
    quotes, errs := api.IntraDay(user, 1)
    return quotes, errs 
}

func (api *Api) IntraDay(user User, interval uint) ([]Quote, []error) {
    ch := make(chan StockQuery)
    for symbol := range user.Stocks {
        quote := &Quote{Symbol: symbol}
        go func(q *Quote) {
            ch <- api.QueryIntraday(q, 1)
        }(quote)
    }

    n := len(user.Stocks)
    quotes := make([]Quote, n)
    errs := make([]error, 0)
    for i := 0; i < n; i++ {
        query := <-ch
        if query.Error != nil {
            errs = append(errs, query.Error)
        } else {
            quotes[i] = query.Quote
        }
    }
    
    if len(errs) > 0 {
        return quotes, errs
    }
    return quotes, nil 
}

func (api *Api) QueryIntraday(quote *Quote, interval uint) StockQuery {
    var query StockQuery 
    var stockData StockData
    if quote.Symbol == "" {
        query.Error = fmt.Errorf("Stock object missing symbol")
        return query 
    }
    
    uri := "function=TIME_SERIES_INTRADAY&symbol="+quote.Symbol+"&interval="+strconv.FormatUint(uint64(interval), 10)+"min&apikey="+api.Key

    resp, err := http.Get(URL+uri)
    if err != nil {
        query.Error = fmt.Errorf("Could not get function results from Alpha Vantage: %v", err)
        return query 
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        query.Error = fmt.Errorf("Error getting the response body: %v", err)
        return query 
    }

    if err := json.Unmarshal(body, &stockData); err != nil {
        query.Error = fmt.Errorf("Error unmarshalling stock: %v", err)
        return query 
    }

    for k, v := range stockData {
        if k == "Information" {
            query.Quote.Symbol = v.(map[string]interface{})["Meta Data"].(map[string]interface{})["Symbol"].(string)
            query.Error = fmt.Errorf("Alpha Vantage API error: %v", v)
        }
    }

    metadata := stockData["Meta Data"].(map[string]interface{}) 
    query.Quote.Symbol = metadata["2. Symbol"].(string)
    for k, v := range stockData {
        if strings.Contains(k, "Time Series") {
            for d, v2 := range v.(map[string]interface{}) {
                date := metadata["3. Last Refreshed"].(string)
                if d == date {
                    timezone := metadata["6. Time Zone"].(string)
                    stockDate, err := ToTime(d, timezone)
                    if err != nil {
                        query.Error = err
                        return query 
                    }

                    v3 := v2.(map[string]interface{})

                    query.Quote.Date = stockDate
                    query.Quote.TimeZone = timezone //t.Format("2006-1-2 15:04:05") 

                    query.Quote.Open, err = strconv.ParseFloat(v3["1. open"].(string), 64)
                    if err != nil {
                        query.Error = err
                        return query 
                    }

                    query.Quote.High, err = strconv.ParseFloat(v3["2. high"].(string), 64)
                    if err != nil {
                        query.Error = err
                        return query 
                    }

                    query.Quote.Low, err = strconv.ParseFloat(v3["3. low"].(string), 64)
                    if err != nil {
                        query.Error = err
                        return query 
                    }

                    query.Quote.Close, err = strconv.ParseFloat(v3["4. close"].(string), 64)
                    if err != nil {
                        query.Error = err
                        return query 
                    }

                    query.Quote.Volume, err = strconv.ParseInt(v3["5. volume"].(string), 10, 64)
                    if err != nil {
                        query.Error = err
                        return query 
                    }

                    break;
                }
            }
            break;
        }
    }

    return query 
}

func ToTime(d string, timezone string) (time.Time, error) {
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
