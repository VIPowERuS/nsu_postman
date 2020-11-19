package apiserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/VIPowERuS/nsu_postman/internal/app/model"
)

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/", s.indexHandler())
	s.router.HandleFunc("/writeAnnouncement", s.writeAnnouncementHandler())
	s.router.HandleFunc("/SaveAnnouncement", s.saveAnnouncementHandler())
	s.router.HandleFunc("/login", s.loginCheck()).Methods("POST")
	s.router.HandleFunc("/login", s.loginUser()).Methods("GET")
	s.router.HandleFunc("/writeMail", s.writeMail()).Methods("GET")
	s.router.HandleFunc("/sendMail", s.sendMail()).Methods("POST")
}

func (s *APIServer) indexHandler() http.HandlerFunc {
	s.logger.Info("Index handler was called")
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
	s.logger.Info("Write announcement was called")
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
			fmt.Fprintf(w, err.Error())
			return
		}
		t.ExecuteTemplate(w, "login", nil)
	}
}

func (s *APIServer) loginCheck() http.HandlerFunc { // "POST" method
	s.logger.Info("Login check was called")
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/templates/login.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		mail := r.FormValue("mail")
		password := r.FormValue("password")
		usr, err := s.store.User().FindByEmail(mail)
		if tempPassword := model.ToHash(password); tempPassword == usr.EncryptedPassword {
			session, err := s.sessionStore.Get(r, "mail")
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			session.Values["user_id"] = usr.ID
			if err := s.sessionStore.Save(r, w, session); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			http.Redirect(w, r, "/", 302)
			return
		}
		usr = nil
		//fmt.Println(usr)

		t.ExecuteTemplate(w, "login", nil)
	}
}

func (s *APIServer) writeMail() http.HandlerFunc {
	s.logger.Info("Write Mail was called")
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/templates/mail.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		t.ExecuteTemplate(w, "mail", nil)
	}
}

func (s *APIServer) sendMail() http.HandlerFunc {
	s.logger.Info("Send Mail was called")
	return func(w http.ResponseWriter, r *http.Request) {
		var data = MailData{r.FormValue("receiver"), r.FormValue("subject"), r.FormValue("content")}
		if err := s.sendMails(data); err != nil {
			s.logger.Error(err)
			s.respond(w, r, http.StatusInternalServerError, err)
		}
		s.logger.Info("mail was sended")
		t, err := template.ParseFiles("internal/templates/mail.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		t.ExecuteTemplate(w, "mail", nil)
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
