package asset

import (
    "fmt"
    "github.com/asaskevich/govalidator"
    reader "github.com/pmuaddib/sedb/entity/asset/reader/db"
    "github.com/pmuaddib/sedb/entity/asset/writer/db"
    apperror "github.com/pmuaddib/sedb/entity/error"
    "math"
    "time"
)

const (
    checkFailed = "Send failed, check data"
    dataNotSaved = "Data not saved"
    dateConvertFormat = "2006-01-02 15:01:05"
    dateLayout = "02-Jan-06 15:04:05"
    negType = "withdrawal"
)

// Asset data type represents client request data
type Asset struct {
    ID int `valid:"-"`
    UserID int `json:"user_id,string" valid:"required"`
    Currency string `json:"currency" valid:"in(USD|EUR),required"`
    Amount float64 `json:"amount,float64" valid:"required"`
    Time string `json:"time_placed" valid:"required"`
    Type string `json:"type" valid:"in(withdrawal|deposit), required"`
}

// Add validator
func init() {
    govalidator.CustomTypeTagMap.Set("customDateValidator", func(i interface{}, context interface{}) bool {
        switch v := context.(type) {
        case *Asset:
            if _, err := time.Parse(dateLayout, v.Time); err != nil {
                return false
            }
            return true
        }
        return false
    })
}

// Validate checks incoming data from client
func (a *Asset) Validate() error {
    if ok, err := govalidator.ValidateStruct(a); !ok {
        apperror.ProcessFatal(err)
        return fmt.Errorf("%s", checkFailed)
    }
    dateValid, ok := govalidator.CustomTypeTagMap.Get("customDateValidator")
    if !ok {
        apperror.ProcessFatal("Not found validator")
        return fmt.Errorf("%s", checkFailed)
    }
    if ok := dateValid(nil, a); !ok {
        apperror.ProcessFatal("Date not valid")
        return fmt.Errorf("%s", checkFailed)
    }
    a.convertDate()
    a.convertAmount()

    return nil
}

// GetBalance returns aggregated data for client
func (a Asset) GetBalance() ([]Asset, error) {
    var result []Asset
    dbReader := reader.Reader{}
    res, err:= dbReader.ReadBalance(a.UserID)
    if err != nil {
        return nil, err
    }
    var i Asset
    for k, v := range res {
        i.Amount = k
        i.Currency = v
        result = append(result, i)
    }

    return result, nil
}

// Write saves data
func (a Asset) Write() error {
    dbWriter := db.Writer{}
    // userID, currency, amount, time, typeVal
    err := dbWriter.Write(a.UserID, a.Currency, a.Amount, a.Time, a.Type)
    if err != nil {
        return fmt.Errorf("%s", dataNotSaved)
    }

    return nil
}

// Convert time format
func (a *Asset) convertDate() {
    t, _ := time.Parse(dateLayout, a.Time)
    a.Time = t.Format(dateConvertFormat)
}

// Convert amount value
func (a *Asset) convertAmount() {
    if a.Type == negType {
        am := math.Abs(a.Amount)
        a.Amount = -am
    }
}
