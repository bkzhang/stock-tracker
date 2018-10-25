package core

import (
    "database/sql"
    "fmt"
    "strings"
    "time"

    "github.com/go-sql-driver/mysql"
)

type Database struct {
    User string
    Password string
    Name string
    Api *Api
}

func (db *Database) GetUser(username string) (User, error) {
    var user User

    sqldb, err := sql.Open("mysql", db.User+":"+db.Password+"@/"+db.Name)
    if err != nil {
        return user, fmt.Errorf("Failed to connect to the database server: %v", err)
    }
    defer sqldb.Close()

    var id string
    checkUserExists := sqldb.QueryRow("SELECT id FROM user where username = ?", username)
    err = checkUserExists.Scan(&id)
    if err == sql.ErrNoRows {
        return user, fmt.Errorf("user does not exist")
    } else if err != nil {
        return user, fmt.Errorf("mysql query error:", err)
    }

    query := `SELECT o.symbol, o.date, o.timezone, o.price, o.shares FROM user u
    INNER JOIN user_ownedstock uo
    ON u.id = uo.user_id
    INNER JOIN ownedstock o
    ON o.id = uo.ownedstock_id
    WHERE u.id = ?`

    user.Name = username
    rows, err := sqldb.Query(query, id)
    if err != nil {
        return user, fmt.Errorf("error retrieving user:", err)
    }
    defer rows.Close()

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
    err = row.Scan(&quote.Symbol, &nt, &ntz, &nh, &nl, &no, &nc, &nv)
    if err == sql.ErrNoRows {
        quote, err2 := db.UpsertQuote(symbol)
        if err2 != nil {
            return quote, fmt.Errorf("error updating quote:", err2)
        }
        return quote, nil
    } else if err != nil {
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

    storedQuoteDate := time.Date(quote.Date.Year(), quote.Date.Month(), quote.Date.Day(), quote.Date.Hour(), quote.Date.Minute(), 0, 0, quote.Date.Location())
    now := time.Now()
    t := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
    
    if storedQuoteDate != t && !(storedQuoteDate.Year() == t.Year() && storedQuoteDate.Month() == t.Month() && storedQuoteDate.Day() == t.Day() && (storedQuoteDate.Hour() >= 16 || storedQuoteDate.Hour() <= 9)) {
        quote, err = db.UpsertQuote(symbol)
        if err != nil {
            return quote, fmt.Errorf("error updating quote:", err)
        }
    }


    return quote, nil
}

func (db *Database) UpsertQuote(symbol string) (Quote, error) {
    quote, err := db.Api.GetQuote(symbol)
    if err != nil {
        return quote, err
    }

    sqldb, err := sql.Open("mysql", db.User+":"+db.Password+"@/"+db.Name)
    if err != nil {
        return quote, fmt.Errorf("Failed to connect to the database server: %v", err)
    }
    defer sqldb.Close()

    upsert := `INSERT INTO quote (symbol, date, timezone, high, low, open, close, volume)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    ON DUPLICATE KEY UPDATE symbol = ?, date = ?, timezone = ?, high = ?, low = ?, open = ?, close = ?, volume = ?`
    stmtUps, err := sqldb.Prepare(upsert)
    if err != nil {
        return quote, fmt.Errorf("error preparing statement for mysql:", err)
    }

    if _, err := stmtUps.Exec(quote.Symbol, quote.Date, quote.TimeZone, quote.High, quote.Low, quote.Open, quote.Close, quote.Volume, quote.Symbol, quote.Date, quote.TimeZone, quote.High, quote.Low, quote.Open, quote.Close, quote.Volume); err != nil {
        return quote, fmt.Errorf("error upserting into mysql:", err)
    }
    
    return quote, nil
}

func (db *Database) BuyStocks(buy Buy, username string) error { 
    sqldb, err := sql.Open("mysql", db.User+":"+db.Password+"@/"+db.Name)
    if err != nil {
        return fmt.Errorf("Failed to connect to the database server: %v", err)
    }
    defer sqldb.Close()

    stmts := []string{ 
        "INSERT INTO ownedstock (symbol, date, timezone, price, shares) VALUES (?, ?, ?, ?, ?)", 
        "INSERT INTO user_ownedstock (user_id, ownedstock_id) VALUES ((select id from user where username = ?), ?)",
        "INSERT INTO ownedstock_quote (ownedstock_id, quote_id) Values (?, (select id from quote where symbol = ?))",
    }

    tx, err := sqldb.Begin()
    if err != nil {
        return fmt.Errorf("error starting mysql transaction:", err)
    }

    defer tx.Rollback()

    for k, v := range buy {
        quote, err := db.Api.GetQuote(k)
        if err != nil {
            return err
        }

        price := RoundFloat64((quote.High+quote.Low)/2.0, 0.01)

        res, err := tx.Exec(stmts[0], strings.ToUpper(k), quote.Date, quote.TimeZone, price, v)
        if err != nil {
            return fmt.Errorf("mysql transaction error:", err)
        }

        ownedstock_id, err := res.LastInsertId()
        if err != nil {
            return fmt.Errorf("mysql transaction error:", err)
        }

        res, err = tx.Exec(stmts[1], username, ownedstock_id)  
        if err != nil {
            return fmt.Errorf("mysql transaction error:", err)
        }

        res, err = tx.Exec(stmts[2], ownedstock_id, strings.ToUpper(k))  
        if err != nil {
            return fmt.Errorf("mysql transaction error:", err)
        }
    }

    return tx.Commit() 
}
