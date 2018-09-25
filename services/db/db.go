package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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

// creates a connection to the database with the given configuration
func init() {
	dblogin := os.Getenv("DBLOGIN")
	if dblogin == "" {
		log.Fatal("$DBLOGIN must be set")
	}
	cfg := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     dblogin,
		Password: os.Getenv("DBPASSWORD"),
		DBName:   "mydevto",
	}
	if !cfg.valid() {
		log.Fatal("Invalid configuration to connect to database")
	}

	db, err := sql.Open("postgres", cfg.String())

	if err != nil {
		log.Fatalf("Errored opening connection to database: %s", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Errored pinging database: %s", err)
	}
}

// DB returns a pointer to the existing connection. Note that it might be nil
// if Open hasn't been called before hand
func DB() *sql.DB {
	return db
}
