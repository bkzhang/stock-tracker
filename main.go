package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"

    "github.com/bkzhang/stock-tracker/core"
    "github.com/gorilla/mux"
    //"github.com/gorilla/websocket"
)

func main() {
    var (
        apikey = flag.String("apikey", "demo", "Alpha Vantage api key")
        dbcollection = flag.String("dbcollection", "stocks", "mongodb database collection name")
        dbname = flag.String("dbname", "stocksapi", "mongodb database name")
        dbserver = flag.String("dbserver", "mongodb://127.0.0.1:27017", "mongodb server")
        port = flag.String("port", "8080", "port")
    )
    flag.Parse()

    //ADD WEBSOCKET (OR CHANGE TO) FOR CLI

    r := mux.NewRouter()
    c := &core.Controller{
        ApiKey: &core.Api{
            Key: *apikey,
        },
        DB: &core.Database {
            Server: *dbserver,
            Name: *dbname,
            Collection: *dbcollection,
        },
        Router: r,
    }

    //write tests
    //r.Methods("GET", "POST").Path("/user/{user}").HandlerFunc(c.UserStocks).Name("UserStocks")
    r.Methods("GET").Path("/user/{user}/function/{function}").HandlerFunc(c.Function).Name("Function")

    r.PathPrefix("/test").HandlerFunc(c.FunctionSocket).Name("FunctionSocket")

    fmt.Println("Serving to http://localhost:", *port)
    log.Fatal(http.ListenAndServe(":"+*port, r))
}
