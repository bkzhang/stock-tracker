package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"

    "github.com/bkzhang/stock/core"
    "github.com/gorilla/mux"
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

    r := mux.NewRouter()
    c := &core.Controller{
        ApiKey: *apikey,
        DB: &core.Database {
            Server: *dbserver,
            Name: *dbname,
            Collection: *dbcollection,
        },
        Router: r,
    }

    //write tests
    r.Methods("GET", "POST").Path("/user/{user}").HandlerFunc(c.UserStocks).Name("UserStocks")
    r.Methods("GET").Path("/user/{user}/function/{function}").HandlerFunc(c.Function).Name("Function")

    fmt.Println("Serving to http://localhost:", *port)
    log.Fatal(http.ListenAndServe(":"+*port, r))
}
