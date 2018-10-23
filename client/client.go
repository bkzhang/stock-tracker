package main

import (
    "flag"
    "log"
    "net/url"
    "os"
    "os/signal"
    "time"

    "github.com/gorilla/websocket"
)

func main() {
    var (
        user = flag.String("user", "", "username")
        function = flag.String("function", "time_series_intraday", "stock function, currently available: times_series_intraday")
    )
    flag.Parse()

    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/user/"+*user+"/function/"+*function}
    log.Println("Connecting to", u.String())

    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        log.Fatal("Dial error:", err)
    }
    defer c.Close()

    done := make(chan struct{})

    go func() {
        defer close(done)
        for {
            _, message, err := c.ReadMessage()
            if err != nil {
                log.Println("Read error:", err)
                return
            }
            log.Println("Received:", string(message))
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
