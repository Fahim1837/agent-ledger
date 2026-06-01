package config

import (
	"database/sql"
	"fmt"
	"log"

	"time"

	"os"

	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

func GetDBConnection() *sql.DB {
	log.Print("DB_DRIVER=", os.Getenv("DB_DRIVER"))
	d := Database{
		host:     os.Getenv("DB_HOST"),
		username: os.Getenv("DB_USERNAME"),
		password: os.Getenv("DB_PASSWORD"),
		driver:   os.Getenv("DB_DRIVER"),
		port:     os.Getenv("DB_PORT"),
		dbname:   os.Getenv("DB_NAME"),
	}

	db, err := d.connectDB()
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	return db
}

type Database struct {
	host     string
	port     string
	username string
	password string
	driver   string
	dbname   string
}

func (d Database) connectPostgreSQL() (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		d.username,
		d.password,
		d.host,
		d.port,
		d.dbname,
	)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open Postgres DB: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {

		if err := db.Close(); err != nil {
			return nil, fmt.Errorf("failed to close the postgres db: %w", err)
		}

		return nil, fmt.Errorf("failed to ping postgres db: %w", err)
	}

	log.Println("Connected to Postgres successfully!")
	return db, nil
}

func (d Database) connectSQLite() (*sql.DB, error) {
	dsn := fmt.Sprintf("./data/%s", d.dbname)

	path := "./data/"
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create sqlite directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	if err := db.Ping(); err != nil {

		if err := db.Close(); err != nil {
			return nil, fmt.Errorf("failed to close the sqlite db: %w", err)
		}

		return nil, fmt.Errorf("failed to ping sqlite db: %w", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {

		if err := db.Close(); err != nil {
			return nil, fmt.Errorf("failed to close the sqlite db: %w", err)
		}

		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	log.Println("Connected to SQLite successfully!")
	return db, nil
}

func (d Database) connectDB() (*sql.DB, error) {

	switch d.driver {
	case "postgres":
		return d.connectPostgreSQL()
	case "sqlite":
		return d.connectSQLite()
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", d.driver)
	}
}
