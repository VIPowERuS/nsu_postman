package apiserver

import (
	"net/http"
	"net/smtp"

	"github.com/VIPowERuS/nsu_postman/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Config ...
type Config struct {
	BindAddr     string `toml:"bind_addr"`
	LogLevel     string `toml:"log_level"`
	SMTPMail     string `toml:"smtp_mail"`
	SMTPPassword string `toml:"smtp_password"`
	Store        *store.Config
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr:     ":8000",
		LogLevel:     "debug",
		SMTPMail:     "nsu.postman@gmail.com",
		SMTPPassword: "progasosat'",
		Store:        store.NewConfig(),
	}
}

// APIServer ...
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

// New ...
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start ...
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()
	_ = s.configureMailing()
	if err := s.configureStore(); err != nil {
		return err
	}
	s.logger.Info("api server is started")
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureMailing() smtp.Auth {
	return smtp.PlainAuth("", s.config.SMTPMail, s.config.SMTPPassword, "smtp.gmail.com")
}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}
