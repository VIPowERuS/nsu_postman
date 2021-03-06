package store

import (
	"database/sql"

	_ "github.com/lib/pq" // ...
)

// Config ...
type Config struct {
	DatabaseURL string `toml:"database_url"`
}

// NewConfig return new initialized struct
func NewConfig() *Config {
	return &Config{}
}

// Store ...
type Store struct {
	config         *Config
	db             *sql.DB
	userRepository *UserRepository
}

// New ...
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

// Open database with config data
func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	s.db = db
	return nil
}

// Close ...
func (s *Store) Close() {
	s.db.Close()
}

// User ...
func (s *Store) User() *UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}
	return s.userRepository
}
