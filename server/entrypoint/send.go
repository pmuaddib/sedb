package entrypoint

import (
    "encoding/json"
    "fmt"
    "github.com/pmuaddib/sedb/entity"
    "github.com/pmuaddib/sedb/entity/asset"
    apperror "github.com/pmuaddib/sedb/entity/error"
)

type Send struct {
}

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
