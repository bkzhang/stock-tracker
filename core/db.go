package core

import (
    "fmt"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type Database struct {
    Server string
    Name string
    Collection string
}

func (db *Database) User(username string) (User, error) {
    session, err := mgo.Dial(db.Server) 
    if err != nil {
        return nil, fmt.Errorf("Failed to connect to the database server: %v", err)
    }
    defer session.Close()

    c := session.DB(db.Name).C(db.Collection)

    var user User
    if err := c.Find(bson.M{"name": username}).One(&user); err != nil {
        return nil, fmt.Errorf("Failed to get user info: %v", err)
    }

    return user, nil
}

func (db *Database) Function(username string, fnName string) (Stocks err) {
    user, err := db.User(username)
    if err != nil {
        return nil, err
    }

    ch := make(chan Stock)
    for _, fn := range user.Functions {
        if fn.Name == "times_series_intraday" {
            for _, stock := range fn.Stocks {
                go func(s Stock) {
                   ch <- GetStock(s) // implement and move this function over to new file, implement w/ timer
                }(stock)
            }
            close(ch) // might close before everything sends
            break;
        }
    }
    
    stocks := make(Stock)
    for {
        select {
        case stock := <-ch:
            stocks = append(stocks, stock)
        default:
            return stocks, nil
        }
    }
}

// change this since model changed
func (db *Database) Stocks(username string) (Stocks, error) {
    session, err := mgo.Dial(db.Server)
    if err != nil {
        return nil, fmt.Errorf("Failed to connect to the database server: %v", err)
    }
    defer session.Close()

    c := session.DB(db.Name).C(db.Collection)

    var user User 
    if err := c.Find(bson.M{"name": username}).One(&user); err != nil {
        return nil, fmt.Errorf("Failed to get the list of stocks saved: %v", err)
    }

    return user.Stocks, nil
}

// change this since model changed
func (db *Database) AddStock(user string, stock Stock) error {
    stocks, err := db.Stocks(user)
    if err != nil {
        return err
    }

    stocks = append(stocks, stock)

    session, err := mgo.Dial(db.Server)
    if err != nil {
        return fmt.Errorf("Failed to connect to the database server: %v", err)
    }
    defer session.Close()

    c := session.DB(db.Name).C(db.Collection)

    if err := c.Update(bson.M{"name": user}, bson.M{"$set": bson.M{"stocks": stocks}}); err != nil {
        return fmt.Errorf("Failed to insert stock into the database: %v", err)
    }
    
    fmt.Println(stock.Symbol, "added to the database")
    return nil
}
