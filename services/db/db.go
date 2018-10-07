package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// initiate the drivers for postgresql
	_ "github.com/lib/pq"
	"github.com/math2001/mydevto/services/uli"
)

var db *sql.DB

// a dodgy escape for the connection string
func escape(s string) string {
	return fmt.Sprintf("'%s'", strings.Replace(s, "'", "\\'", -1))
}

// Config holds the configuration to connect to the database
type Config struct {
	// Host to connect to
	Host string
	// Port to connect to
	Port string
	// User who will connect
	User string
	// Password to use
	Password string
	// DBName name to connect to (it has to have been created before hand)
	DBName string
}

func (cfg Config) String() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		escape(cfg.Host), escape(cfg.Port), escape(cfg.DBName), escape(cfg.User), escape(cfg.Password))
}

func (cfg Config) valid() bool {
	return cfg.Host != "" && cfg.Port != "" && cfg.User != "" && cfg.DBName != ""
}

// Init creates a connection to the database with the given configuration
// It should be called only once
func init() {
	dblogin := os.Getenv("DBLOGIN")
	if dblogin == "" {
		log.Fatal("$DBLOGIN must be set")
	}
	dbname := os.Getenv("DBNAME")
	if dbname == "" {
		log.Fatal("$DBNAME must be set")
	}
	cfg := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     dblogin,
		Password: os.Getenv("DBPASSWORD"), // password can be empty
		DBName:   dbname,
	}
	if !cfg.valid() {
		log.Fatal("Invalid configuration to connect to database")
	}

	var err error
	db, err = sql.Open("postgres", cfg.String())

	if err != nil {
		log.Fatalf("Errored opening connection to database: %s", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Errored pinging database: %s", err)
	}
}

// QueryContext executes a query, logging and timing it
func QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	uli.Printf(ctx, "DB: query %q %v", query, args)
	defer uli.Printf(ctx, "DB: query %q %v done after %.2f", query, args,
		time.Now().Sub(start).Seconds())
	return db.QueryContext(ctx, query, args...)
}

// QueryRowContext executes a query, logging and timing it
func QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	uli.Printf(ctx, "DB: query row %q %v", query, args)
	defer uli.Printf(ctx, "DB: query row %q %v done after %.2f", query, args,
		time.Now().Sub(start).Seconds())
	return db.QueryRowContext(ctx, query, args...)
}

// ExecContext executes a query, logging and timing it
func ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	uli.Printf(ctx, "DB: exec %q %v", query, args)
	defer uli.Printf(ctx, "DB: exec %q %v done after %.2f", query, args,
		time.Now().Sub(start).Seconds())
	return db.ExecContext(ctx, query, args...)
}
