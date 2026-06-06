package data

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	loadEnv("../../.env")

	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		panic("TEST_DB_DSN not set")
	}

	var err error
	testDB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic("could not open test db: " + err.Error())
	}
	defer testDB.Close()

	if err = testDB.Ping(); err != nil {
		panic("could not ping test db: " + err.Error())
	}

	os.Exit(m.Run())
}

func loadEnv(path string) {
	f, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(f), "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "export ")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				os.Setenv(key, val)
			}
		}
	}
}
