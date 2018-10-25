package core

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "math"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
)

type Controller struct {
    DB *Database
    Router *mux.Router
}

func (c *Controller) User(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    user, err := c.DB.GetUser(vars["user"])
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    data, err := json.Marshal(user)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}

func (c *Controller) AddUser(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    if err := r.Body.Close(); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    var user User
    if err := json.Unmarshal(body, &user); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    log.Println(user)

    if err := c.DB.AddUser(user); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
}

func (c *Controller) Quote(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    symbol := vars["symbol"]
    quote, err := c.DB.GetQuote(symbol)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    data, err := json.Marshal(quote)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}

func (c *Controller) IntraDay(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    upgrader := websocket.Upgrader{}
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("upgrader error:", err)
    }
    defer conn.Close()

    user, err := c.DB.GetUser(vars["user"])
    if err != nil {
        if e := conn.WriteMessage(websocket.TextMessage, []byte("user error: " + err.Error())); e != nil {
            log.Println("write error:", e)
            return
        }
        log.Println("user error:", err)
        return
    }

    type Earnings struct {
        Quotes []Quote
        GainsLosses map[string]float64
    }

    cache := make(map[string]float64)
    numShares := make(map[string]int)

    for symbol, shares := range user.Stocks {
        for _, share := range shares { 
            cache[symbol] += share.Price * float64(share.Shares)
            numShares[symbol] += share.Shares
        }
    }

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            return
        }

        res, errs := c.DB.Api.IntraDayOneMin(user)
        if errs != nil {
            for _, err := range errs {
                errstring := vars["function"] + " error: " + err.Error()
                if e := conn.WriteMessage(websocket.TextMessage, []byte(errstring)); e != nil {
                    log.Println("write error:", e)
                    return
                }
                log.Println(errstring)
            }
        }

        earning := &Earnings{
            Quotes: res,
            GainsLosses: make(map[string]float64),
        }

        for _, v := range res {
            earning.GainsLosses[v.Symbol] = RoundFloat64((v.Open+v.Close)/2.0*float64(numShares[v.Symbol]) - cache[v.Symbol], 0.01)
        }

        data, err := json.Marshal(earning)
        if err != nil {
            if e := conn.WriteMessage(websocket.TextMessage, []byte("marshalling error: " + err.Error())); e != nil {
                log.Println("write error:", e)
                return
            }
            log.Println("marshalling error: " + err.Error())
        }

        if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
            log.Println("write error:", err)
            return
        }
    }
}

func RoundFloat64(f, decimalPlace float64) float64 {
    return math.Round(f / decimalPlace) * decimalPlace
}
