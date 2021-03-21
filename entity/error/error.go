package error

import (
    "fmt"
    "io"
    "log"
)

const (
    server = "Something went wrong, please check logs."
    Client = "Something went wrong, please try again later."
)

// ProcessServErr used to store logs and notify client
func ProcessServErr(err interface{}) {
    fmt.Println(server)
    log.Println(err)
}

// ProcessFatal prints generic message
func ProcessFatal(msg interface{}) {
    log.Println(msg)
}

// OutErr prints message by using io.Writer
func OutErr(o io.Writer, err interface{}) {
    fmt.Fprintln(o, err)
}
