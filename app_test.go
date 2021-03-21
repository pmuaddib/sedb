package main

import (
    "bytes"
    "github.com/pmuaddib/sedb/entity/dbase"
    apperror "github.com/pmuaddib/sedb/entity/error"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

// Test case with wrong http method
func TestSendWrongMethod(t *testing.T) {
    var jsonStr = []byte(`{"user_id": "9000", "currency": "EUR", "amount": 100, "time_placed" : "20-JAN-20 10:00:00", "type": "deposit"}`)
    req, err := http.NewRequest("GET", "/send", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Send)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusBadGateway {
        t.Errorf("Send returned wrong status code: got %v want %v",
            status, http.StatusBadGateway)
    }
}

// Test case with correct data
func TestSendCorrectMethod(t *testing.T) {
    var jsonStr = []byte(`{"user_id": "9000", "currency": "EUR", "amount": 100, "time_placed" : "20-JAN-20 10:00:00", "type": "deposit"}`)
    req, err := http.NewRequest("POST", "/send", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Send)
    handler.ServeHTTP(rr, req)
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Send returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := success
    s := strings.Trim(rr.Body.String(), "\n")
    if s != expected {
       t.Errorf("handler returned unexpected body: got %v want %v", s, expected)
    }

    cleanDbTable(9000, t)
}

// Test case with missed currency field
func TestSendMissedField(t *testing.T) {
    var jsonStr = []byte(`{"user_id": "9000", "amount": 100, "time_placed" : "20-JAN-20 10:00:00", "type": "deposit"}`)
    req, err := http.NewRequest("POST", "/send", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Send)
    handler.ServeHTTP(rr, req)
    ok := strings.Contains(rr.Body.String(), apperror.Client)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Send returned wrong status code: got %v want %v", status, http.StatusOK)
    }
    if !ok {
        t.Errorf("Field validation failed")
    }
}

// Test case with getting item
func TestGetItem(t *testing.T) {
    var jsonStr = []byte(`{"user_id": "9000", "currency": "EUR", "amount": 100, "time_placed" : "20-JAN-20 10:00:00", "type": "deposit"}`)
    req, err := http.NewRequest("POST", "/send", bytes.NewBuffer(jsonStr))
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Send)
    handler.ServeHTTP(rr, req)

    expected := success
    s := strings.Trim(rr.Body.String(), "\n")
    if s != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", s, expected)
    }
    req, err = http.NewRequest("GET", "/get?userId=9000", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr = httptest.NewRecorder()
    handler = http.HandlerFunc(Get)
    handler.ServeHTTP(rr, req)
    expected = `[{"Currency":"EUR","Balance":"100.000"}]`
    got := strings.Trim(rr.Body.String(), "\n")
    if got != expected {
        t.Errorf("can't get item: got %v, expected %v", got, expected)
    }

    cleanDbTable(9000, t)
}

// Test case with getting several items
func TestGetItems(t *testing.T) {
    data := []string{
        `{"user_id": "9000", "currency": "EUR", "amount": -100.99999, "time_placed" : "20-JAN-20 10:00:00", "type": "withdrawal"}`,
        `{"user_id": "9000", "currency": "EUR", "amount": 100, "time_placed" : "20-JAN-20 10:00:00", "type": "deposit"}`,
        `{"user_id": "9000", "currency": "EUR", "amount": 100, "time_placed" : "20-JAN-20 10:00:00", "type": "deposit"}`,
        `{"user_id": "9000", "currency": "USD", "amount": -100.0101, "time_placed" : "20-JAN-20 10:00:00", "type": "withdrawal"}`,
        `{"user_id": "9000", "currency": "USD", "amount": 100, "time_placed" : "20-JAN-20 10:00:00", "type": "deposit"}`,
    }
    for _, v := range data {
        req, err := http.NewRequest("POST", "/send", bytes.NewBuffer([]byte(v)))
        if err != nil {
            t.Fatal(err)
        }
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(Send)
        handler.ServeHTTP(rr, req)
    }

    req, err := http.NewRequest("GET", "/get?userId=9000", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Get)
    handler.ServeHTTP(rr, req)
    expectedUsd := `"Currency":"USD","Balance":"-0.010"`
    expectedEur := `"Currency":"EUR","Balance":"99.000"`

    if !strings.Contains(rr.Body.String(), expectedUsd) {
        t.Errorf("item not exist,: got %v, expected %v", rr.Body.String(), expectedUsd)
    }
    if !strings.Contains(rr.Body.String(), expectedEur) {
        t.Errorf("item not exist,: got %v, expected %v", rr.Body.String(), expectedEur)
    }

    cleanDbTable(9000, t)
}

// Clean data after test
func cleanDbTable(userId int, t *testing.T) {
    d := dbase.DB{}
    err := d.InitConnection()
    if err != nil {
        t.Errorf("Can't connect from DB")
    }
    defer d.Conn.Close()
    _, err = d.Conn.Exec("DELETE FROM assets WHERE user_id = ?", userId)
    if err != nil {
        t.Errorf("Can't remove test item from DB")
    }
}
