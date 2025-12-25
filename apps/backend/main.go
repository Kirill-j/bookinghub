package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"bookinghub-backend/internal/handler"
	"bookinghub-backend/internal/repo"
)

type App struct {
	DB *sqlx.DB
}

func main() {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")

	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "bookinghub")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("db open error: %v", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	app := &App{DB: db}

	resourceRepo := repo.NewResourceRepo(db)
	resourceHandler := handler.NewResourceHandler(resourceRepo)
	categoryRepo := repo.NewCategoryRepo(db)
	categoryHandler := handler.NewCategoryHandler(categoryRepo)

	r := chi.NewRouter()

	// Логи + базовая защита от паники
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})

		r.Get("/categories", categoryHandler.List)

		r.Get("/resources", resourceHandler.List)
		r.Post("/resources", resourceHandler.Create)
	})

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
