package main

import (
    "fmt"
    "github.com/pmuaddib/sedb/entity"
    "github.com/pmuaddib/sedb/entity/config/yaml"
    apperror "github.com/pmuaddib/sedb/entity/error"
    entry "github.com/pmuaddib/sedb/server/entrypoint"
    "io/ioutil"
    "log"
    "net/http"
)

// {"user_id": "134256", "currency": "EUR", "amount": 1000, "time_placed" : "24-JAN-20 10:27:44", "type": "deposit"}
// curl -X POST -d '{"user_id": "134256", "currency": "EUR", "amount": 1000, "time_placed" : "24-JAN-20 10:27:44", "type": "deposit"}' 127.0.0.1:8082/send

const (
    success = "Ok"
    reqUserName = "userId"
)

// Main entrypoint
func main() {
    conf := &yaml.Storage{}
    mux := registerMux()
    addr := getAddr(conf)

    log.Fatal(http.ListenAndServe(addr, mux))
}

// Register handlers
func registerMux() *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/send", Send)
    mux.HandleFunc("/get", Get)

    return mux
}

// Return server address
func getAddr(v entity.Config) string {
    storage := v.(*yaml.Storage)

    h, ok := storage.GetVal("server","host")
    if !ok {
        apperror.ProcessFatal(fmt.Errorf("%s", "Can't get host value"))
    }
    p, ok := storage.GetVal("server", "port")
    if !ok {
        apperror.ProcessFatal(fmt.Errorf("%s", "Can't get port value"))
    }

    return fmt.Sprintf("%s:%s", h, p)
}

// Send handler
func Send(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusBadGateway)
        apperror.ProcessFatal("endpoint is not POST method")
        apperror.OutErr(w, apperror.Client)
        return
    }
    var data []string
    s, _ := ioutil.ReadAll(r.Body)
    data = append(data, string(s))

    sendEntrypoint := entry.Send{}
    err := sendEntrypoint.Execute(data)
    if err != nil {
        apperror.OutErr(w, apperror.Client)
        return
    }
    fmt.Fprintf(w, "%s\n", success)
}

// Get handler
func Get(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Redirect(w, r, "/", http.StatusBadGateway)
        apperror.ProcessFatal("endpoint is not GET method")
        apperror.OutErr(w, apperror.Client)
        return
    }
    err := r.ParseForm()
    if err != nil {
        http.Redirect(w, r, "/", http.StatusBadGateway)
        apperror.ProcessFatal("can't parse data")
        apperror.OutErr(w, apperror.Client)
        return
    }
    var data []string
    for k, v := range r.Form {
        if k == reqUserName {
            data = v
            break
        }
    }
    if nil == data || len(data) < 1 {
        apperror.OutErr(w, apperror.Client)
        return
    }
    get := entry.Get{}
    get.SetWriter(w)
    err = get.Execute(data)
    if err != nil {
        apperror.OutErr(w, apperror.Client)
        return
    }
}
