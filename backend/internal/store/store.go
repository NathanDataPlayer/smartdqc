package store

import (
    "database/sql"
    "os"
    "time"
    _ "github.com/go-sql-driver/mysql"
)

func NullFloat64FromFloat(f float64) sql.NullFloat64 { return sql.NullFloat64{Float64:f, Valid:true} }

type Rule struct {
    ID int
    Name string
    Table string
    Type string
    Status string
    LastRun string
    DB string
    Partition string
    Threshold sql.NullFloat64
}

var db *sql.DB

func Init() error {
    dsn := os.Getenv("DQC_MYSQL_DSN")
    if dsn == "" {
        dsn = "dqc:dqc@tcp(127.0.0.1:3306)/dqc?parseTime=true&charset=utf8mb4"
    }
    x, err := sql.Open("mysql", dsn)
    if err != nil { return err }
    x.SetConnMaxLifetime(30 * time.Minute)
    x.SetMaxOpenConns(10)
    x.SetMaxIdleConns(5)
    if err = x.Ping(); err != nil { return err }
    db = x
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS rules(
        id INT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        table_name VARCHAR(255) NOT NULL,
        type VARCHAR(64) NOT NULL,
        status VARCHAR(64) NOT NULL,
        last_run VARCHAR(64) NULL,
        db_name VARCHAR(255) NULL,
        partition_key VARCHAR(255) NULL,
        threshold DOUBLE NULL
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)
    if err != nil { return err }
    ensureColumn("rules", "db_name", "ALTER TABLE rules ADD COLUMN db_name VARCHAR(255) NULL")
    ensureColumn("rules", "partition_key", "ALTER TABLE rules ADD COLUMN partition_key VARCHAR(255) NULL")
    ensureColumn("rules", "threshold", "ALTER TABLE rules ADD COLUMN threshold DOUBLE NULL")
    if err = initAlertSchema(); err != nil { return err }
    return nil
}

func ensureColumn(table, column, alterSQL string) {
    var cnt int
    row := db.QueryRow("SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = ?", table, column)
    if row.Scan(&cnt) == nil && cnt == 0 {
        db.Exec(alterSQL)
    }
}

func ListRules() ([]Rule, error) {
    rows, err := db.Query("SELECT id,name,table_name,type,status,last_run,db_name,partition_key,threshold FROM rules ORDER BY id")
    if err != nil { return nil, err }
    defer rows.Close()
    var out []Rule
    for rows.Next() {
        var r Rule
        rows.Scan(&r.ID,&r.Name,&r.Table,&r.Type,&r.Status,&r.LastRun,&r.DB,&r.Partition,&r.Threshold)
        out = append(out,r)
    }
    return out, nil
}

func CreateRule(in Rule) (Rule, error) {
    if in.Status == "" { in.Status = "enabled" }
    if in.LastRun == "" { in.LastRun = time.Now().Format("15:04") }
    res, err := db.Exec("INSERT INTO rules(name,table_name,type,status,last_run,db_name,partition_key,threshold) VALUES(?,?,?,?,?,?,?,?)", in.Name, in.Table, in.Type, in.Status, in.LastRun, in.DB, in.Partition, in.Threshold)
    if err != nil { return Rule{}, err }
    id, _ := res.LastInsertId()
    in.ID = int(id)
    return in, nil
}

func UpdateRule(id int, in Rule) (Rule, error) {
    cur, err := GetRule(id)
    if err != nil { return Rule{}, err }
    if in.Name != "" { cur.Name = in.Name }
    if in.Table != "" { cur.Table = in.Table }
    if in.Type != "" { cur.Type = in.Type }
    if in.Status != "" { cur.Status = in.Status }
    _, err = db.Exec("UPDATE rules SET name=?, table_name=?, type=?, status=?, db_name=?, partition_key=?, threshold=? WHERE id=?", cur.Name, cur.Table, cur.Type, cur.Status, cur.DB, cur.Partition, cur.Threshold, id)
    if err != nil { return Rule{}, err }
    return cur, nil
}

func DeleteRule(id int) error {
    _, err := db.Exec("DELETE FROM rules WHERE id=?", id)
    return err
}

func GetRule(id int) (Rule, error) {
    var r Rule
    err := db.QueryRow("SELECT id,name,table_name,type,status,last_run,db_name,partition_key,threshold FROM rules WHERE id=?", id).Scan(&r.ID,&r.Name,&r.Table,&r.Type,&r.Status,&r.LastRun,&r.DB,&r.Partition,&r.Threshold)
    return r, err
}

func EnsureSeed() {
    rows, err := db.Query("SELECT COUNT(1) FROM rules")
    if err != nil { return }
    defer rows.Close()
    var c int
    if rows.Next() { rows.Scan(&c) }
    if c == 0 {
        CreateRule(Rule{Name:"分区完整性检查", Table:"ods.user_event", Type:"partition", Status:"enabled"})
        CreateRule(Rule{Name:"空值率阈值", Table:"dwd.order_detail", Type:"null_rate", Status:"enabled"})
        CreateRule(Rule{Name:"唯一性检查", Table:"dim.product", Type:"unique", Status:"paused"})
    }
}
