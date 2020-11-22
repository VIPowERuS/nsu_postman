package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"

	"github.com/VIPowERuS/nsu_postman/internal/app/model"
)

const (
	sessionName        = "authCookie"
	ctxKeyUser  ctxKey = iota
)

var (
	errNotAuthenticated = errors.New("Not authenticated")
)

type ctxKey int16

func (s *APIServer) configureRouter() {
	s.router.Use(s.authenticateUser)
	s.router.HandleFunc("/", s.indexHandler())
	s.router.HandleFunc("/writeAnnouncement", s.writeAnnouncementHandler())
	s.router.HandleFunc("/SaveAnnouncement", s.saveAnnouncementHandler())
	s.router.HandleFunc("/login", s.loginCheck()).Methods("POST")
	s.router.HandleFunc("/login", s.loginUser()).Methods("GET")
	s.router.HandleFunc("/writeMail", s.writeMail()).Methods("GET")
	s.router.HandleFunc("/sendMail", s.sendMail()).Methods("POST")

}

func (s *APIServer) getCookie(w http.ResponseWriter, r *http.Request) *model.User {
	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		s.logger.Error("cookies GET error")
		return nil
	}
	id, ok := session.Values["user_id"]
	if !ok {
		s.logger.Info("Not Authenticated user")
		return &model.User{ID: 0, Email: "", Access: 0}
	}
	mail, _ := session.Values["mail"]
	access, _ := session.Values["access"]
	return &model.User{ID: id.(int), Email: mail.(string), Access: access.(int)}
}

func (s *APIServer) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := s.getCookie(w, r)
		s.logger.Info(user)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	})
}

func (s *APIServer) indexHandler() http.HandlerFunc {
	s.logger.Info("Index handler was called")
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/templates/index.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("templates error")
			return
		}
		t.ExecuteTemplate(w, "index", r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *APIServer) writeAnnouncementHandler() http.HandlerFunc {
	s.logger.Info("Write announcement was called")
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/templates/write.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("templates error")
			return
		}

		t.ExecuteTemplate(w, "write", r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *APIServer) saveAnnouncementHandler() http.HandlerFunc {
	s.logger.Info("Save announcement was called")
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		title := r.FormValue("title")
		content := r.FormValue("content")
		s.logger.Infoln(id, title, content)
		http.Redirect(w, r, "/", 302)
	}
}

func (s *APIServer) loginUser() http.HandlerFunc { // "GET" method
	s.logger.Info("Login user was called")
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/templates/login.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("templates error")
			return
		}
		t.ExecuteTemplate(w, "login", r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *APIServer) loginCheck() http.HandlerFunc { // "POST" method
	s.logger.Info("Login check was called")
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/templates/login.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("templates error")
			return
		}
		usr, err := s.store.User().FindByEmail(r.FormValue("mail"))
		if tempPassword := model.ToHash(r.FormValue("password")); tempPassword == usr.EncryptedPassword {
			session, err := s.sessionStore.Get(r, sessionName)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				s.logger.Error("cookies GET error")
				return
			}
			session.Values["user_id"], session.Values["mail"], session.Values["access"] = usr.ID, usr.Email, usr.Access
			if err := s.sessionStore.Save(r, w, session); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				s.logger.Error("cookies SAVE error")
				return
			}
			http.Redirect(w, r, "/", 302)
			return
		}
		usr = nil
		t.ExecuteTemplate(w, "login", r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *APIServer) writeMail() http.HandlerFunc {
	s.logger.Info("Write Mail was called")
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/templates/mail.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("templates error")
			return
		}
		t.ExecuteTemplate(w, "mail", r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *APIServer) sendMail() http.HandlerFunc {
	s.logger.Info("Send Mail was called")
	return func(w http.ResponseWriter, r *http.Request) {
		var data = MailData{r.FormValue("receiver"), r.FormValue("subject"), r.FormValue("content")}
		if err := s.sendMails(data); err != nil {
			s.logger.Error("send mail error:", err)
			s.respond(w, r, http.StatusInternalServerError, err)
		}
		s.logger.Info("mail was sended")
		t, err := template.ParseFiles("internal/templates/mail.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("templates error")
			return
		}
		t.ExecuteTemplate(w, "mail", r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *APIServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}
func (s *APIServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}

}
