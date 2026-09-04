package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	logger "github.com/ElJhonnypro/url-shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type URLRepository struct {
	log           *logger.Logger
	isDBConnected bool
	database      *sql.DB
}

func NewURLRepository(log *logger.Logger) *URLRepository {
	return &URLRepository{
		log: log,
	}
}

func (r *URLRepository) Connect() (bool, error) {
	postgresURL := os.Getenv("POSTGRES_URL")

	if postgresURL == "" {
		err := errors.New("POSTGRES_URL is not set in the environment variables")

		r.log.Error("Database connection failed: " + err.Error())

		return false, err
	}

	r.log.Info("Attempting to connect to PostgreSQL...")

	db, err := sql.Open("pgx", postgresURL)
	if err != nil {
		r.log.Error("Failed to initialize database driver: " + err.Error())

		return false, fmt.Errorf("failed to open database handle: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(15 * time.Minute)

	// Test the actual database connection.
	// sql.Open() does not actually connect to PostgreSQL.
	if err := db.Ping(); err != nil {
		r.log.Error("Failed to connect to PostgreSQL: " + err.Error())

		db.Close()

		return false, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	r.log.Info("Successfully connected to PostgreSQL.")

	query := `CREATE TABLE IF NOT EXISTS urls (
		code VARCHAR(10) PRIMARY KEY,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
		r.log.Error("Failed to create urls table: " + err.Error())

		db.Close()

		return false, fmt.Errorf("failed to create urls table: %w", err)
	}

	r.database = db
	r.isDBConnected = true

	r.log.Info("Database initialized successfully.")

	return true, nil
}

func (r *URLRepository) Close() error {
	if r.database != nil {
		return r.database.Close()
	}
	return nil
}

// Create inserta el código y la URL original en PostgreSQL

func (r *URLRepository) Create(url string, code string) error {
	if !r.isDBConnected || r.database == nil {
		err := errors.New("database connection is not active")

		r.log.Error("Failed to create URL: " + err.Error())

		return err
	}

	query := `INSERT INTO urls (code, original_url) VALUES ($1, $2)`

	_, err := r.database.Exec(query, code, url)
	if err != nil {
		r.log.Error(fmt.Sprintf(
			"Failed to insert code %s: %v",
			code,
			err,
		))

		return fmt.Errorf("failed to create url mapping: %w", err)
	}

	r.log.Info(fmt.Sprintf(
		"Successfully saved code: %s for url: %s",
		code,
		url,
	))

	return nil
}

func (r *URLRepository) Get(code string) (string, error) {
	if !r.isDBConnected || r.database == nil {
		return "", fmt.Errorf("database connection is not active")
	}

	query := `SELECT original_url FROM urls WHERE code = $1`

	var originalURL string
	err := r.database.QueryRow(query, code).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Warn("CODE NOT FOUND. Information: " + code)
			return "", fmt.Errorf("code not found")
		}
		r.log.Error(fmt.Sprintf("Failed to query code %s: %v", code, err))
		return "", fmt.Errorf("database query error: %w", err)
	}

	r.log.Info(fmt.Sprintf("CODE FOUND. Information: %s -> %s", code, originalURL))
	return originalURL, nil
}
