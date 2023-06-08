package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	status := &SafeStatus{en: "qabż", fa: "قبض"}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	r.Get("/fa", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(status.Get("fa")))
	})

	r.Get("/en", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(status.Get("en")))
	})

	r.With(BasicAuth).Post("/poke", func(w http.ResponseWriter, r *http.Request) {
		status.Toggle()
		w.Write([]byte(status.Get("en")))
	})

	http.ListenAndServe(":8080", r)
}

//
// Status type and methods
//

type SafeStatus struct {
	mux sync.RWMutex
	en  string
	fa  string
}

func (s *SafeStatus) Get(lang string) string {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if lang == "fa" {
		return s.fa
	} else {
		return s.en
	}
}

func (s *SafeStatus) Toggle() {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.en == "qabż" {
		s.en = "basṭ"
		s.fa = "بسط"
	} else {
		s.en = "qabż"
		s.fa = "قبض"
	}
}

//
// Basic Auth middleware
//

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || user != "soroush" || pass != os.Getenv("SUPER_SECRET_KEY") {
			http.Error(w, "Begone, infidel", 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}
