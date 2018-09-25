package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// a dodgy escape for the connection string
func escape(s string) string {
	return fmt.Sprintf("'%s'", strings.Replace(s, "'", "\\'", -1))
}

// Conn holds the connection to the database
type Conn struct {
	cfg Config
	DB  *sql.DB
}

// Close releases resources that 'sql/database' used
func (conn *Conn) Close() error {
	if conn.DB == nil {
		return errors.New("unexisting connection to database")
	}
	if err := conn.DB.Close(); err != nil {
		return errors.Wrapf(err,
			"errored closing the connection")
	}
	return nil
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

// Open creates a connection to the database with the given connection
func Open(cfg Config) (Conn, error) {
	var conn Conn
	if !cfg.valid() {
		return conn, fmt.Errorf("invalid configuration %s", cfg)
	}
	conn.cfg = cfg

	db, err := sql.Open("postgres", cfg.String())

	if err != nil {
		return conn, errors.Wrapf(err, "errored opening connection to database %s", cfg)
	}

	if err := db.Ping(); err != nil {
		return conn, errors.Wrapf(err, "errored pinging database %s", cfg)
	}

	conn.DB = db

	return conn, nil
}
