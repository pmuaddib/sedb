package db

import (
    "fmt"
    "github.com/pmuaddib/sedb/entity/dbase"
    apperror "github.com/pmuaddib/sedb/entity/error"
)

// Table name to work with
const tableName = "assets"

// Reader used to retrieve data
type Reader struct {
}

// ReadBalance returns balance information
func (r Reader) ReadBalance(userId int) (map[float64]string, error) {
    db := dbase.DB{}
    err := db.InitConnection()
    if err != nil {
        return nil, fmt.Errorf("%s", "Failed to get data")
    }
    defer db.Conn.Close()

    q := fmt.Sprintf(
        "SELECT t.currency, SUM(amount) AS balance FROM %s AS t WHERE user_id = %d GROUP BY currency",
        tableName, userId)
    res, err := db.Conn.Query(q)
    if err != nil {
        apperror.ProcessFatal(err)
        return nil, fmt.Errorf("%s", "Failed to get data")
    }

    var i struct{
        Currency string
        Balance float64
    }
    balanceList := make(map[float64]string)

    for res.Next() {
       e := res.Scan(&i.Currency, &i.Balance)
       if e != nil {
           apperror.ProcessFatal(e)
       }
       balanceList[i.Balance] = i.Currency
    }

    return balanceList, nil
}
