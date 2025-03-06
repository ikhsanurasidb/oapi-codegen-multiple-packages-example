package mysql

import (
    "database/sql"
    "fmt"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

type Config struct {
    Host     string
    Port     int
    User     string
    Password string
    Database string
}

func NewConnection(cfg Config) (*sql.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    if err = db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping MySQL: %w", err)
    }

    return db, nil
}