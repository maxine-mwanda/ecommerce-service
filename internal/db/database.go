// internal/db/database.go
package db

import (
	"database/sql"
	"ecommerce-service/internal/config"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ConnectMySQL(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=5s&readTimeout=5s&writeTimeout=5s&multiStatements=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func RunSQLFile(db *sql.DB, filePath string) error {
	sqlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read SQL file: %v", err)
		return err
	}

	sqlStatements := string(sqlBytes)

	_, err = db.Exec(sqlStatements)
	if err != nil {
		log.Printf("Failed to execute SQL statements: %v", err)
		return err
	}

	fmt.Println("SQL schema applied successfully.")
	return nil
}
