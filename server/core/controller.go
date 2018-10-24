package core

import (
    "encoding/json"
    "io/ioutil"
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
    quote, err := c.DB.GetQuote(vars["symbol"])
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

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            return
        }

        res, errs := c.ApiKey.IntraDayOneMin(user)
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
