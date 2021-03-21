package db

import (
    "fmt"
    "sedb/entity/dbase"
    apperror "sedb/entity/error"
)

// Table name to work with
const tableName = "assets"

// Writer used to communicate with db to save data
type Writer struct {
}

// Write saves data
func (w Writer) Write(userID int, currency string, amount float64, time string, typeVal string) error {
    db := dbase.DB{}
    err := db.InitConnection()
    if err != nil {
        apperror.ProcessFatal(err)
        return fmt.Errorf("%s", "Failed to save data")
    }
    defer db.Conn.Close()

    ins := fmt.Sprintf("INSERT INTO %s (user_id, currency, amount, time_placed, type) VALUES(?, ?, ?, ?, ?)", tableName)
    res, err := db.Conn.Exec(ins, userID, currency, amount, time, typeVal)

    if err != nil {
        apperror.ProcessFatal(err)
        return fmt.Errorf("%s", "Failed to save data")
    }
    _, err = res.LastInsertId()
    if err != nil {
        apperror.ProcessFatal(err)
        return fmt.Errorf("%s", "Failed to save data")
    }

    return nil
}
