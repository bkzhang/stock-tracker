package core

import (
    "encoding/json"
    //"io/ioutil"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
)

type Controller struct {
    ApiKey *Api
    DB *Database
    Router *mux.Router
}

/*func (c *Controller) Function(w http.ResponseWriter, r *http.Request) {
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
}*/

func (c *Controller) Function(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    upgrader := websocket.Upgrader{}
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("upgrader error:", err)
    }
    defer conn.Close()

    user, err := c.DB.User(vars["user"])
    if err != nil {
        if e := conn.WriteMessage(websocket.TextMessage, []byte("user error: " + err.Error())); e != nil {
            log.Println("write error:", e)
            return
        }
        log.Println("user error:", err)
        return
    }

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            return
        }

        res, errs := c.ApiKey.Function(user, vars["function"])
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

        data, err := json.Marshal(res)
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
