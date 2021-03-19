package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)
// {"user_id": "134256", "currency": "EUR", "amount": 1000, "time_placed" : "24-JAN-20 10:27:44", "type": "deposit"}

type asset struct {
    UserID int `json:"user_id,string"`
    Currency string `json:"currency"`
    Amount float64 `json:"amount,float64"`
    Time string `json:"time_placed"`
    Type string `json:"type"`
}
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/send", send)
    mux.HandleFunc("/get", get)

    log.Fatal(http.ListenAndServe(":8082", mux))
}

func send(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusBadGateway)
        log.Print("send got not POST")
        return
    }
    err := r.ParseForm()
    if err != nil {
        http.Redirect(w, r, "/", http.StatusBadGateway)
        log.Print("send can't parse data")
        return
    }
    var items []*asset
    item := &asset{}
    for k, v := range r.Form {
        err := json.Unmarshal([]byte(k), item)
        if err != nil {
            log.Println("can't unmarshal", err)
        }
        fmt.Fprintf(w, "[key '%+v'] %v\n", item, v)
        items = append(items, item)
    }
    fmt.Fprintln(w, items)
    //fmt.Fprint(w, "Send", r.Form)
}

func get(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Get")
}