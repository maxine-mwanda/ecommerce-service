package tests

import (
	"fmt"
	"os"
	"testing"

	"path/filepath"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	rootPath, err := filepath.Abs("..") // go one directory up from /tests
	if err != nil {
		fmt.Println("Failed to determine project root")
		os.Exit(1)
	}

	// Step 2: Construct full path to .env.test
	envPath := filepath.Join(rootPath, ".env.test")

	err = godotenv.Load(envPath)
	if err != nil {
		fmt.Println("Failed to load .env.test:", err)
	} else {
		fmt.Println("Loaded .env.test in tests")
		fmt.Println("TEST_DB_HOST:", os.Getenv("TEST_DB_HOST"))
	}

	os.Exit(m.Run())
}
