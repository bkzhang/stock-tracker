package core

type User struct {
    Name string
    Function Functions
}

type Stock struct {
    Symbol string
    Amount uint
}

type Function struct {
    Name string
    Stocks Stocks
    Interval uint
}

type Stocks []Stock
type Functions []Function
