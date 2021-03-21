package entrypoint

import (
    "encoding/json"
    "fmt"
    "io"
    "sedb/entity"
    "sedb/entity/asset"
    apperror "sedb/entity/error"
    "strconv"
)

// Get used as endpoint
type Get struct {
    Output io.Writer
}

// Execute retrieves data
func (e Get) Execute(data entity.Data) error {
    userId, _ := strconv.Atoi(data[0])

    item := &asset.Asset{UserID: userId}
    res, err := item.GetBalance()
    if err != nil {
        return err
    }
    var o []struct{
        Currency string
        Balance string
    }
    for _, v := range res {
        o = append(o, struct {
            Currency string
            Balance  string
        }{Currency: v.Currency, Balance: fmt.Sprintf("%.3f", v.Amount)})
    }
    j, err := json.Marshal(o)
    if err != nil {
        apperror.ProcessFatal(err)
        return fmt.Errorf("%s", "Can't process output")
    }
    fmt.Fprintln(e.Output, string(j))

    return nil
}

// SetWriter stores io.Writer
func (e *Get) SetWriter(w io.Writer) {
    e.Output = w
}
