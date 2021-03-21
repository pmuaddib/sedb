package dbase

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/pmuaddib/sedb/entity/config/yaml"
    apperror "github.com/pmuaddib/sedb/entity/error"
)

const maxConns = 10

// DB store connection
type DB struct {
    Conn *sql.DB
}

// InitConnection initialize DB connection
func (d *DB) InitConnection() error {
    conn, err := sql.Open("mysql", getSource())
    if err != nil {
        apperror.ProcessFatal(err)
        return fmt.Errorf("%s", "Can't save data, try again")
    }
    conn.SetMaxIdleConns(maxConns)
    d.Conn = conn

    return nil
}

// Return DB source
func getSource() string {
    config := yaml.Storage{}

    u, ok := config.GetVal("db", "user")
    if !ok {
        apperror.ProcessFatal("db user not found in config")
    }

    p, ok := config.GetVal("db", "pass")
    if !ok {
        apperror.ProcessFatal("db pass not found in config")
    }

    h, ok := config.GetVal("db", "host")
    if !ok {
        apperror.ProcessFatal("db host not found in config")
    }

    port, ok := config.GetVal("db", "port")
    if !ok {
        apperror.ProcessFatal("db port not found in config")
    }

    dbname, ok := config.GetVal("db", "dbname")
    if !ok {
        apperror.ProcessFatal("db name not found in config")
    }

    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", u, p, h, port, dbname)
}
