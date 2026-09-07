package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// ============================================================
// Connect
// ============================================================
//
// Membuat koneksi ke database MySQL menggunakan konfigurasi
// dari file .env.
// ============================================================

func Connect() (*sql.DB, error) {

	// --------------------------------------------------------
	// Load .env
	// --------------------------------------------------------

	err := godotenv.Load(".env")

	if err != nil {
		return nil, fmt.Errorf(
			"failed to load .env: %w",
			err,
		)
	}

	// --------------------------------------------------------
	// Ambil konfigurasi database
	// --------------------------------------------------------

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	// --------------------------------------------------------
	// Buat DSN
	// --------------------------------------------------------

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user,
		password,
		host,
		port,
		name,
	)

	// --------------------------------------------------------
	// Buka koneksi
	// --------------------------------------------------------

	db, err := sql.Open(
		"mysql",
		dsn,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to open database: %w",
			err,
		)
	}

	// --------------------------------------------------------
	// Test koneksi
	// --------------------------------------------------------

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf(
			"failed to connect to database: %w",
			err,
		)
	}

	return db, nil
}
