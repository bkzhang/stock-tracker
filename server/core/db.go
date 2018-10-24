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

    query := fmt.Sprintf(`SELECT s.symbol, s.date, s.timezone, s.high, s.low, s.open, s.close, s.volume FROM user u
    INNER JOIN user_stock us
    ON u.id = us.user_id
    INNER JOIN stock s
    ON s.id = us.stock_id
    WHERE u.username = '%s'`, username)

    rows, err := sqldb.Query(query)
    if err != nil {
        return user, fmt.Errorf("error retrieving user:", err)
    }

    user.Name = username
    for rows.Next() {
        var stock Stock
        var nt mysql.NullTime
        var ntz sql.NullString
        var nh, nl, no, nc sql.NullFloat64 
        var nv sql.NullInt64
        if err := rows.Scan(&stock.Symbol, &nt, &ntz, &nh, &nl, &no, &nc, &nv); err != nil {
            return user, fmt.Errorf("error retrieving user's tracked stocks:", err)
        }
        
        if nt.Valid {
            stock.Date = nt.Time
        }
        if ntz.Valid {
            stock.TimeZone = ntz.String
        }
        if nh.Valid {
            stock.High = nh.Float64
        }
        if nl.Valid {
            stock.Low = nl.Float64
        }
        if no.Valid {
            stock.Open = no.Float64
        }
        if nc.Valid {
            stock.Close = nc.Float64
        }
        if nv.Valid {
            stock.Volume = nv.Int64
        }

        user.Stocks = append(user.Stocks, stock)
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

func (db *Database) GetStock(symbol string) (Stock, error) {
    var stock Stock
    symbol = strings.ToUpper(symbol)

    sqldb, err := sql.Open("mysql", db.User+":"+db.Password+"@/"+db.Name)
    if err != nil {
        return stock, fmt.Errorf("Failed to connect to the database server: %v", err)
    }
    defer sqldb.Close()

    query := fmt.Sprintf("SELECT symbol, date, timezone, high, low, open, close, volume FROM stock where symbol = '%s'", symbol) 

    row := sqldb.QueryRow(query)
    if err := row.Scan(&stock.Symbol, &stock.Date, &stock.TimeZone, &stock.High, &stock.Low, &stock.Open, &stock.Close, &stock.Volume); err != nil { 
        return stock, fmt.Errorf("error retrieving stock:", err)
    }

    return stock, nil
}

func ScanStock(row *sql.Row) Stock {
    var stock Stock
    var nt mysql.NullTime
    var ntz sql.NullString
    var nh, nl, no, nc sql.NullFloat64 
    var nv sql.NullInt64
    if err := row.Scan(&stock.Symbol, &nt, &ntz, &nh, &nl, &no, &nc, &nv); err != nil {
        return user, fmt.Errorf("error retrieving user's tracked stocks:", err)
    }
    
    if nt.Valid {
        stock.Date = nt.Time
    }
    if ntz.Valid {
        stock.TimeZone = ntz.String
    }
    if nh.Valid {
        stock.High = nh.Float64
    }
    if nl.Valid {
        stock.Low = nl.Float64
    }
    if no.Valid {
        stock.Open = no.Float64
    }
    if nc.Valid {
        stock.Close = nc.Float64
    }
    if nv.Valid {
        stock.Volume = nv.Int64
    }

    return stock
}
