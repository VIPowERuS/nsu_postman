package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/VIPowERuS/nsu_postman/internal/app/model"
	"github.com/VIPowERuS/nsu_postman/internal/app/store"
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
	s.router.HandleFunc("/writePost", s.writePostHandler())
	s.router.HandleFunc("/savePost", s.savePostHandler())
	s.router.HandleFunc("/editPost", s.editPostHandler())
	s.router.HandleFunc("/login", s.loginCheckHandler()).Methods("POST")
	s.router.HandleFunc("/login", s.loginUserHandler()).Methods("GET")
	s.router.HandleFunc("/logout", s.logoutHandler())
	s.router.HandleFunc("/writeMail", s.writeMailHandler()).Methods("GET")
	s.router.HandleFunc("/sendMail", s.sendMailHandler()).Methods("POST")

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
		posts, err := s.store.User().FindAllPosts()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("db error (finding posts)")
			return
		}
		t.ExecuteTemplate(w, "header", r.Context().Value(ctxKeyUser).(*model.User))
		t.ExecuteTemplate(w, "index", posts)
	}
}

func (s *APIServer) writePostHandler() http.HandlerFunc {
	s.logger.Info("Write post was called")
	return func(w http.ResponseWriter, r *http.Request) {
		if cond := r.Context().Value(ctxKeyUser).(*model.User).ID == 0; cond {
			http.Redirect(w, r, "/", 302)
		}
		t, err := template.ParseFiles("internal/templates/write.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("templates error")
			return
		}
		t.ExecuteTemplate(w, "header", r.Context().Value(ctxKeyUser).(*model.User))
		t.ExecuteTemplate(w, "write", nil)
	}
}

func (s *APIServer) savePostHandler() http.HandlerFunc { // only "POST" method
	s.logger.Info("Save post was called")
	return func(w http.ResponseWriter, r *http.Request) {
		if cond := r.Context().Value(ctxKeyUser).(*model.User).ID == 0; cond {
			http.Redirect(w, r, "/", 302)
		}
		stringID := r.FormValue("id")
		post := store.Post{ID: 0, Header: r.FormValue("header"),
			Department: r.FormValue("department"), Content: r.FormValue("content"), Date: ""}
		if stringID != "" { // need to change data
			intID, err := strconv.Atoi(stringID)
			fmt.Print("int id == ", intID)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				s.logger.Error("db error (edit post 1)")
				return
			}
			post.ID = intID
			err = s.store.User().ChangePost(post)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				s.logger.Error("db error (adding post)")
				return
			}
			s.logger.Infof("Post %d was changed", post.ID)
			http.Redirect(w, r, "/", 301)
			return
		} else {
			postID, err := s.store.User().AddPost(post)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				s.logger.Error("db error (adding post)")
				return
			}
			s.logger.Infof("Post %d was added", postID)
			http.Redirect(w, r, "/", 301)
		}
	}
}

func (s *APIServer) editPostHandler() http.HandlerFunc { // only "POST" method
	s.logger.Info("Edit post was called")
	return func(w http.ResponseWriter, r *http.Request) {
		if cond := r.Context().Value(ctxKeyUser).(*model.User).ID == 0; cond {
			http.Redirect(w, r, "/", 302)
		}
		stringID := r.FormValue("id")
		post, err := s.store.User().FindPost(stringID)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			s.logger.Error("db error (edit post 2)")
			return
		}
		ID, err := strconv.Atoi(r.FormValue("id"))
		post.ID = ID
		s.logger.Infof("Post %d was choosen to be changed", post.ID)
		t, err := template.ParseFiles("internal/templates/write.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("templates error")
			return
		}
		t.ExecuteTemplate(w, "header", r.Context().Value(ctxKeyUser).(*model.User))
		t.ExecuteTemplate(w, "write", post)
	}
}

func (s *APIServer) loginUserHandler() http.HandlerFunc { // "GET" method
	s.logger.Info("Login user was called")
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("internal/templates/login.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("templates error")
			return
		}
		t.ExecuteTemplate(w, "header", r.Context().Value(ctxKeyUser).(*model.User))
		t.ExecuteTemplate(w, "login", r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *APIServer) loginCheckHandler() http.HandlerFunc { // "POST" method
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
		t.ExecuteTemplate(w, "header", r.Context().Value(ctxKeyUser).(*model.User))
		t.ExecuteTemplate(w, "login", nil)
	}
}

func (s *APIServer) logoutHandler() http.HandlerFunc {
	s.logger.Info("Logout was called")
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("cookies GET error")
			return
		}
		session.Options.MaxAge = -1
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("cookies DELETE error")
			return
		}
		http.Redirect(w, r, "/", 302)
		return
	}
}

func (s *APIServer) writeMailHandler() http.HandlerFunc {
	s.logger.Info("Write Mail was called")
	return func(w http.ResponseWriter, r *http.Request) {
		if cond := r.Context().Value(ctxKeyUser).(*model.User).ID == 0; cond {
			http.Redirect(w, r, "/", 302)
		}
		t, err := template.ParseFiles("internal/templates/mail.html", "internal/templates/header.html", "internal/templates/footer.html")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			s.logger.Error("templates error")
			return
		}
		t.ExecuteTemplate(w, "header", r.Context().Value(ctxKeyUser).(*model.User))
		t.ExecuteTemplate(w, "mail", nil)
	}
}

func (s *APIServer) sendMailHandler() http.HandlerFunc {
	s.logger.Info("Send Mail was called")
	return func(w http.ResponseWriter, r *http.Request) {
		if cond := r.Context().Value(ctxKeyUser).(*model.User).ID == 0; cond {
			http.Redirect(w, r, "/", 302)
		}
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
		t.ExecuteTemplate(w, "header", r.Context().Value(ctxKeyUser).(*model.User))
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
