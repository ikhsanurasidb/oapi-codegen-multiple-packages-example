package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oapi-codegen-multiple-packages-example/config"
)

func main() {
	// Load configuration
	cfg := config.Get()

	// Connect to MySQL without selecting a database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true",
		cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Create database if not exists
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.MySQL.Database))
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// Use the database
	_, err = db.Exec(fmt.Sprintf("USE %s", cfg.MySQL.Database))
	if err != nil {
		log.Fatalf("Failed to use database: %v", err)
	}

	// Create orders table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
            id BIGINT AUTO_INCREMENT PRIMARY KEY,
            pet_id BIGINT,
            quantity INT,
            ship_date DATETIME,
            status VARCHAR(50),
            complete BOOLEAN DEFAULT false,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
    `)
	if err != nil {
		log.Fatalf("Failed to create orders table: %v", err)
	}

	log.Println("Migration completed successfully!")
}
