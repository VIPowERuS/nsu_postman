package apiserver

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/VIPowERuS/nsu_postman/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

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

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/", s.indexHandler())
	s.router.HandleFunc("/writeAnnouncement", s.writeAnnouncementHandler())
	s.router.HandleFunc("/SaveAnnouncement", s.saveAnnouncementHandler())
}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st
	return nil
}

func (s *APIServer) indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/templates/index.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		t.ExecuteTemplate(w, "index", nil)
	}
}

func (s *APIServer) writeAnnouncementHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/templates/write.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		t.ExecuteTemplate(w, "write", nil)
	}
}

func (s *APIServer) saveAnnouncementHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		title := r.FormValue("title")
		content := r.FormValue("content")
		s.logger.Infoln(id, title, content)
		http.Redirect(w, r, "/", 302)
	}
}
