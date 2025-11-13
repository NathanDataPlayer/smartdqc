package meta

import (
    "database/sql"
    "os"
    _ "github.com/lib/pq"
)

type HiveClient struct{ db *sql.DB }

func New() (*HiveClient, error) {
    dsn := os.Getenv("METASTORE_DSN")
    if dsn == "" {
        // Example: postgres://root:root@127.0.0.1:5432/metastore?sslmode=disable
        dsn = os.Getenv("METASTORE_DEFAULT_DSN")
    }
    if dsn == "" { return nil, sql.ErrConnDone }
    x, err := sql.Open("postgres", dsn)
    if err != nil { return nil, err }
    if err = x.Ping(); err != nil { return nil, err }
    return &HiveClient{db:x}, nil
}

func (hc *HiveClient) Close(){ if hc.db != nil { hc.db.Close() } }

func (hc *HiveClient) Databases() ([]string, error) {
    rows, err := hc.db.Query("SELECT name FROM DBS ORDER BY name")
    if err != nil { return nil, err }
    defer rows.Close()
    var out []string
    for rows.Next(){ var n string; rows.Scan(&n); out = append(out,n) }
    return out, nil
}

func (hc *HiveClient) Tables(db string) ([]string, error) {
    rows, err := hc.db.Query("SELECT T.TBL_NAME FROM TBLS T JOIN DBS D ON T.DB_ID=D.DB_ID WHERE D.NAME=$1 ORDER BY T.TBL_NAME", db)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []string
    for rows.Next(){ var n string; rows.Scan(&n); out = append(out,n) }
    return out, nil
}

type Partition struct{ Name string }

func (hc *HiveClient) Partitions(db, table string) ([]Partition, error) {
    rows, err := hc.db.Query("SELECT P.PART_NAME FROM PARTITIONS P JOIN TBLS T ON P.TBL_ID=T.TBL_ID JOIN DBS D ON T.DB_ID=D.DB_ID WHERE D.NAME=$1 AND T.TBL_NAME=$2 ORDER BY P.PART_NAME", db, table)
    if err != nil { return nil, err }
    defer rows.Close()
    var out []Partition
    for rows.Next(){ var n string; rows.Scan(&n); out = append(out, Partition{Name:n}) }
    return out, nil
}
