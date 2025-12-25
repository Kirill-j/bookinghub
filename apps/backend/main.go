package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type App struct {
	DB *sqlx.DB
}

func main() {
	// Конфиг из env (пока самый простой вариант)
	port := getEnv("PORT", "8080")

	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "bookinghub")

	// DSN для MySQL
	// parseTime=true — чтобы DATETIME нормально читался в time.Time
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("db open error: %v", err)
	}

	// Настройки пула (базово)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	app := &App{DB: db}

	r := chi.NewRouter()

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Проверка связи с БД
	r.Get("/db/ping", app.handleDBPing)

	log.Printf("Backend started on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func (a *App) handleDBPing(w http.ResponseWriter, r *http.Request) {
	if err := a.DB.Ping(); err != nil {
		http.Error(w, "db not reachable: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("db ok"))
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
