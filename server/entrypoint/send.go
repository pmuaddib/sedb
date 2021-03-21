package entrypoint

import (
    "encoding/json"
    "fmt"
    "sedb/entity"
    "sedb/entity/asset"
    apperror "sedb/entity/error"
)

// Send type used for saving data
type Send struct {
}

// Execute endpoint saves data
func (e Send) Execute(data entity.Data) error {
    item := &asset.Asset{}
    for _, v := range data {
        err := json.Unmarshal([]byte(v), item)
        if err != nil {
            msg := fmt.Sprintf("can't unmarshal: %s", err)
            apperror.ProcessFatal(msg)
            return fmt.Errorf("%s", "Send failed, check data")
        }
    }
    err := item.Validate()
    if err != nil {
        return err
    }

    err = item.Write()
    if err != nil {
        return err
    }

    return nil
}
