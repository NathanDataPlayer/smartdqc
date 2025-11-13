package store

import (
    "database/sql"
    "time"
)

type AlertEvent struct {
    ID int `json:"id"`
    Level string `json:"level"`
    Message string `json:"message"`
    Status string `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
    AckedAt sql.NullTime `json:"ackedAt"`
}

func initAlertSchema() error {
    _, err := db.Exec(`CREATE TABLE IF NOT EXISTS alert_events(
        id INT PRIMARY KEY AUTO_INCREMENT,
        level VARCHAR(32) NOT NULL,
        message TEXT NOT NULL,
        status VARCHAR(32) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        acked_at TIMESTAMP NULL
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)
    return err
}

func CreateAlert(level, message string) (AlertEvent, error) {
    res, err := db.Exec("INSERT INTO alert_events(level,message,status) VALUES(?,?,?)", level, message, "open")
    if err != nil { return AlertEvent{}, err }
    id, _ := res.LastInsertId()
    return GetAlert(int(id))
}

func ListAlerts() ([]AlertEvent, error) {
    rows, err := db.Query("SELECT id,level,message,status,created_at,acked_at FROM alert_events ORDER BY id DESC")
    if err != nil { return nil, err }
    defer rows.Close()
    var out []AlertEvent
    for rows.Next() {
        var a AlertEvent
        rows.Scan(&a.ID,&a.Level,&a.Message,&a.Status,&a.CreatedAt,&a.AckedAt)
        out = append(out,a)
    }
    return out, nil
}

func AckAlert(id int) (AlertEvent, error) {
    _, err := db.Exec("UPDATE alert_events SET status='ack', acked_at=NOW() WHERE id=?", id)
    if err != nil { return AlertEvent{}, err }
    return GetAlert(id)
}

func GetAlert(id int) (AlertEvent, error) {
    var a AlertEvent
    err := db.QueryRow("SELECT id,level,message,status,created_at,acked_at FROM alert_events WHERE id=?", id).Scan(&a.ID,&a.Level,&a.Message,&a.Status,&a.CreatedAt,&a.AckedAt)
    return a, err
}

