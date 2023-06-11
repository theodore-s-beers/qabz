package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	_, err = dbpool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS qabz_bast (
			lang TEXT PRIMARY KEY,
			status TEXT NOT NULL
		);
		INSERT INTO qabz_bast (lang, status) VALUES ('en', 'qabż') ON CONFLICT (lang) DO NOTHING
	`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
		os.Exit(1)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	r.Get("/en", func(w http.ResponseWriter, r *http.Request) {
		var status string
		err := dbpool.QueryRow(context.Background(), "SELECT status FROM qabz_bast WHERE lang = 'en'").Scan(&status)
		if err != nil {
			http.Error(w, "Failed to retrieve status", 500)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write([]byte(status))
	})

	r.Get("/fa", func(w http.ResponseWriter, r *http.Request) {
		var status string
		err := dbpool.QueryRow(context.Background(), "SELECT status FROM qabz_bast WHERE lang = 'en'").Scan(&status)
		if err != nil {
			http.Error(w, "Failed to retrieve status", 500)
			return
		}

		if status == "qabż" {
			status = "قبض"
		} else {
			status = "بسط"
		}

		w.Write([]byte(status))
	})

	r.With(BasicAuth).Post("/poke", func(w http.ResponseWriter, r *http.Request) {
		_, err := dbpool.Exec(context.Background(), "UPDATE qabz_bast SET status = (CASE WHEN status = 'qabż' THEN 'basṭ' ELSE 'qabż' END) WHERE lang = 'en'")
		if err != nil {
			http.Error(w, "Failed to update status", 500)
			return
		}

		var status string
		err = dbpool.QueryRow(context.Background(), "SELECT status FROM qabz_bast WHERE lang = 'en'").Scan(&status)
		if err != nil {
			http.Error(w, "Failed to retrieve status", 500)
			return
		}

		w.Write([]byte(status))
	})

	http.ListenAndServe(":8080", r)
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
