package core

import (
    "encoding/json"
    //"io/ioutil"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
)

const URL = "https://www.alphavantage.co/query?"

type Controller struct {
    ApiKey *Api 
    DB *Database
    Router *mux.Router
}

func (c *Controller) FunctionSocket(w http.ResponseWriter, r *http.Request) {
    upgrader := websocket.Upgrader{}
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrader error:", err)
    }
    defer conn.Close()

    for {
        mt, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            break
        }
        log.Println("Received message:", message)
        if err := conn.WriteMessage(mt, message); err != nil {
            log.Println("Write error:", err)
            break
        }
    }
}

func (c *Controller) Function(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    user, err := c.DB.User(vars["user"])
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    switch r.Method {
    case "GET":
        c.GetFunction(w, r, user, vars["function"])
    }
}

func (c *Controller) GetFunction(w http.ResponseWriter, r *http.Request, user User, fn string) {
    res, errs := c.ApiKey.Function(user, fn) 
    if errs != nil {
        for _, err := range errs {
            w.Write([]byte(err.Error()))
        }
    }

    data, err := json.Marshal(res)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}

/*
// change this since model changed
func (c *Controller) UserStocks(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    switch r.Method {
    case "GET":
        c.GetStocks(w, r, vars["user"])
    case "POST":
        c.PostStock(w, r, vars["user"])
    }
}

func (c *Controller) GetStocks(w http.ResponseWriter, r *http.Request, user string) {
    stocks, err := c.DB.Stocks(user)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    data, err := json.Marshal(stocks)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}

// change this isnce model changed
func (c *Controller) PostStock(w http.ResponseWriter, r *http.Request, user string) {
    var stock Stock

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

    if err := json.Unmarshal(body, &stock); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    if err := c.DB.AddStock(user, stock); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
}
*/
