package core

import (
    "database/sql"
    "fmt"
    "strings"
//    "time"

    "github.com/go-sql-driver/mysql"
)

type Database struct {
    User string
    Password string
    Name string
}

func (db *Database) GetUser(username string) (User, error) {
    var user User

    sqldb, err := sql.Open("mysql", db.User+":"+db.Password+"@/"+db.Name)
    if err != nil {
        return user, fmt.Errorf("Failed to connect to the database server: %v", err)
    }
    defer sqldb.Close()

    query := fmt.Sprintf(`SELECT o.symbol, o.date, o.timezone, o.price, o.shares FROM user u
    INNER JOIN user_ownedstock uo
    ON u.id = uo.user_id
    INNER JOIN ownedstock o
    ON o.id = uo.ownedstock_id
    WHERE u.username = '%s'`, username)

    rows, err := sqldb.Query(query)
    if err != nil {
        return user, fmt.Errorf("error retrieving user:", err)
    }

    user.Name = username
    for rows.Next() {
        var stock OwnedStock
        var symbol string
        var nt mysql.NullTime
        var ntz sql.NullString
        if err := rows.Scan(&symbol, &nt, &ntz, &stock.Price, &stock.Shares); err != nil {
            return user, fmt.Errorf("error retrieving user's tracked stocks:", err)
        }
        
        if nt.Valid {
            stock.Date = nt.Time
        }
        if ntz.Valid {
            stock.TimeZone = ntz.String
        }

        if user.Stocks == nil {
            user.Stocks = make(map[string][]OwnedStock)
        }
        user.Stocks[symbol] = append(user.Stocks[symbol], stock)
    }
    if err := rows.Err(); err != nil {
        return user, fmt.Errorf("mysql row error:", err)
    }

    return user, nil
}


func (db *Database) AddUser(user User) error { 
    sqldb, err := sql.Open("mysql", db.User+":"+db.Password+"@/"+db.Name)
    if err != nil {
        return fmt.Errorf("Failed to connect to the database server: %v", err)
    }
    defer sqldb.Close()

    stmtIns, err := sqldb.Prepare("INSERT INTO user (username) VALUES (?)")
    if err != nil {
        return fmt.Errorf("mysql statement prepare error:", err)
    }
    defer stmtIns.Close()

    if _, err := stmtIns.Exec(user.Name); err != nil {
        return fmt.Errorf("error adding user:", err)
    }

    return nil
}

func (db *Database) GetQuote(symbol string) (Quote, error) {
    var quote Quote 
    symbol = strings.ToUpper(symbol)

    sqldb, err := sql.Open("mysql", db.User+":"+db.Password+"@/"+db.Name)
    if err != nil {
        return quote, fmt.Errorf("Failed to connect to the database server: %v", err)
    }
    defer sqldb.Close()

    query := fmt.Sprintf("SELECT symbol, date, timezone, high, low, open, close, volume FROM quote where symbol = '%s'", symbol) 

    var nt mysql.NullTime
    var ntz sql.NullString
    var nh, nl, no, nc sql.NullFloat64 
    var nv sql.NullInt64

    row := sqldb.QueryRow(query)
    if err := row.Scan(&quote.Symbol, &nt, &ntz, &nh, &nl, &no, &nc, &nv); err != nil {
        return quote, fmt.Errorf("error retrieving user's tracked stocks:", err)
    }
    
    if nt.Valid {
        quote.Date = nt.Time
    }
    if ntz.Valid {
        quote.TimeZone = ntz.String
    }
    if nh.Valid {
        quote.High = nh.Float64
    }
    if nl.Valid {
        quote.Low = nl.Float64
    }
    if no.Valid {
        quote.Open = no.Float64
    }
    if nc.Valid {
        quote.Close = nc.Float64
    }
    if nv.Valid {
        quote.Volume = nv.Int64
    }

    return quote, nil
}
