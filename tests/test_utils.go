package tests

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func findProjectRoot() string {
	// Start here and walk upward
	dir, _ := os.Getwd()
	for i := 0; i < 10; i++ {
		envPath := filepath.Join(dir, ".env.test")
		if _, err := os.Stat(envPath); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	panic("Could not find .env.test")
}

func SetupTestDB(t *testing.T) *sql.DB {
	projectRoot := findProjectRoot()
	envPath := filepath.Join(projectRoot, ".env.test")
	err := godotenv.Load(envPath)
	if err != nil {
		t.Fatalf("Failed to load .env.test from %s: %v", envPath, err)
	}

	fmt.Println("TEST_DB_USER:", os.Getenv("TEST_DB_USER"))
	fmt.Println("TEST_DB_HOST:", os.Getenv("TEST_DB_HOST"))

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_NAME"),
	)
	fmt.Println("DSN:", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	//drop tables before creating new ones
	_, _ = db.Exec(`
		DROP TABLE IF EXISTS order_items;
		DROP TABLE IF EXISTS orders;
		DROP TABLE IF EXISTS products;
		DROP TABLE IF EXISTS categories;
		DROP TABLE IF EXISTS customers;
	`)

	// Initialize test schema
	if _, err := db.Exec(testSchema); err != nil {
		t.Fatalf("Failed to initialize test schema: %v", err)
	}

	t.Cleanup(func() {
		db.Exec("DROP TABLE IF EXISTS order_items, orders, products, categories, customers")
		db.Close()
	})

	return db
}

const testSchema = `
CREATE TABLE customers (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id VARCHAR(36)
);

CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    description TEXT,
    category_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders (
    id VARCHAR(36) PRIMARY KEY,
    customer_id VARCHAR(36) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
	status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
	order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL
);
`
