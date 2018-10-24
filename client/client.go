package main

import (
    "encoding/json"
    "flag"
    "log"
    "net/url"
    "os"
    "os/signal"
    "time"

    "github.com/bkzhang/stock-tracker/client/model"
    "github.com/gorilla/websocket"
)

func main() {
    var (
        user = flag.String("user", "", "username")
        function = flag.String("function", "intraday", "stock function, currently available: intraday")
    )
    flag.Parse()

    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/user/"+*user+"/function/"+*function}
    log.Println("connecting to", u.String(), "\n")

    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        log.Fatal("dial error:", err)
    }
    defer c.Close()

    done := make(chan struct{})

    go func() {
        defer close(done)
        for {
            _, message, err := c.ReadMessage()
            if err != nil {
                log.Println("read error:", err)
                return
            }

            var stocks model.Stocks
            if err := json.Unmarshal(message, &stocks); err != nil {
                log.Println("json unmarshalling error:", err)
            }

            log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
            for _, v := range stocks {
                t := v.Date
                log.Println(t.Format("2006/1/2 15:04:05"), v.TimeZone, v.Symbol, "high:", v.High, "low:", v.Low, "open:", v.Open, "close:", v.Close, "volume:", v.Volume)
            }
            log.SetFlags(log.Flags() |^ (log.Ldate | log.Ltime))
        }
    }()

    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()

    if err := c.WriteMessage(websocket.TextMessage, []byte("")); err != nil {
        log.Println("Write error:", err)
        return
    }

    for {
        select {
        case <-done:
            return
        case t := <-ticker.C:
            if err := c.WriteMessage(websocket.TextMessage, []byte(t.String())); err != nil {
                log.Println("Write error:", err)
                return
            }
        case <-interrupt:
            log.Println("Interrupt")

            if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
                log.Println("Write error:", err)
                return
            }
            select {
            case <-done:
            case <-time.After(time.Second):
            }
            return
        }
    }
}
