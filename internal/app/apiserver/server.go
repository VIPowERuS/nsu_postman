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
	s.router.HandleFunc("/login", s.loginUser())
	s.router.HandleFunc("/afterLogin", s.loginCheck())
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

func (s *APIServer) loginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/templates/login.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		t.ExecuteTemplate(w, "login", nil)
	}
}

func (s *APIServer) loginCheck() http.HandlerFunc {
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
			http.RedirectHandler("/", 300)
			return
		} else {
			usr = nil
		}
		//fmt.Println(usr)

		t.ExecuteTemplate(w, "login", nil)
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
