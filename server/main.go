package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"

    "github.com/bkzhang/stock-tracker/server/core"
    "github.com/gorilla/mux"
)

func main() {
    var (
        apikey = flag.String("apikey", "demo", "Alpha Vantage api key")
        dbuser = flag.String("dbuser", "root", "mysql database user")
        dbpw = flag.String("dbpw", "", "mysql database user password")
        dbname = flag.String("dbname", "stock_tracker", "mysql database name")
        port = flag.String("port", "8080", "port")
    )
    flag.Parse()

    r := mux.NewRouter()
    c := &core.Controller{
        DB: &core.Database {
            User: *dbuser,
            Password: *dbpw,
            Name: *dbname,
            Api: &core.Api{
                Key: *apikey,
            },
        },
        Router: r,
    }

    //write tests
    r.Methods("POST").Path("/user").HandlerFunc(c.AddUser).Name("AddUser")
    r.Methods("GET").Path("/user/{user}").HandlerFunc(c.User).Name("User")
    r.Methods("GET").Path("/quote/{symbol}").HandlerFunc(c.Quote).Name("Quote")
    r.Methods("POST").Path("/user/{user}/buy").HandlerFunc(c.BuyStocks).Name("BuyStocks")

    r.PathPrefix("/user/{user}/function/intraday").HandlerFunc(c.IntraDay).Name("IntraDay")

    fmt.Println("Serving to http://localhost:" + *port)
    log.Fatal(http.ListenAndServe(":"+*port, r))
}
