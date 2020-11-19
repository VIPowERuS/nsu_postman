package apiserver

import (
	"net/http"
	"net/smtp"
	"strings"

	"github.com/VIPowERuS/nsu_postman/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

// Config ...
type Config struct {
	BindAddr     string `toml:"bind_addr"`
	LogLevel     string `toml:"log_level"`
	SMTPMail     string `toml:"smtp_mail"`
	SMTPPassword string `toml:"smtp_password"`
	Store        *store.Config
	SessionKey   string `toml:"session_key"`
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
	config       *Config
	logger       *logrus.Logger
	router       *mux.Router
	store        *store.Store
	sessionStore sessions.Store
	mail         smtp.Auth
}

// New ...
func New(config *Config) *APIServer {
	return &APIServer{
		config:       config,
		logger:       logrus.New(),
		router:       mux.NewRouter(),
		sessionStore: sessions.NewCookieStore([]byte(config.SessionKey)),
	}
}

// Start ...
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()
	s.mail = s.configureMailing()
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

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

func (s *APIServer) configureMailing() smtp.Auth {
	return smtp.PlainAuth("", s.config.SMTPMail, s.config.SMTPPassword, "smtp.gmail.com")
}

// MailData ...
type MailData struct {
	To      string
	Subject string
	Body    string
}

func (s *APIServer) sendMails(data MailData) error {
	msg := ("To: " + data.To + " " + "\r\n" +
		"Subject:" + data.Subject + " \r\n" +
		"\r\n" +
		data.Body)
	if err := smtp.SendMail("smtp.gmail.com:587", s.mail, s.config.SMTPMail, strings.Split(data.To, ", "), []byte(msg)); err != nil {
		return err
	}
	return nil
}
