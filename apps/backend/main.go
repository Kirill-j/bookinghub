package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"bookinghub-backend/internal/db"
	"bookinghub-backend/internal/domain"
	"bookinghub-backend/internal/handler"
	"bookinghub-backend/internal/repo"
	"bookinghub-backend/internal/service"
)

type App struct {
	DB *sqlx.DB
}

func main() {
	_ = godotenv.Load()
	jwtSecret := getEnv("JWT_SECRET", "dev-secret")
	ttlStr := getEnv("JWT_ACCESS_TTL_MIN", "15")
	ttlMin, _ := strconv.Atoi(ttlStr)
	authSvc := service.NewAuthService(jwtSecret, ttlMin)

	port := getEnv("PORT", "8080")

	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "bookinghub")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	dbx, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("db open error: %v", err)
	}
	dbx.SetMaxOpenConns(10)
	dbx.SetMaxIdleConns(5)
	dbx.SetConnMaxLifetime(30 * time.Minute)

	if err := db.ApplyMigrations(dbx.DB, "./migrations"); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	app := &App{DB: dbx}

	resourceRepo := repo.NewResourceRepo(dbx)
	resourceHandler := handler.NewResourceHandler(resourceRepo)
	categoryRepo := repo.NewCategoryRepo(dbx)
	categoryHandler := handler.NewCategoryHandler(categoryRepo)
	userRepo := repo.NewUserRepo(dbx)
	authHandler := handler.NewAuthHandler(userRepo, authSvc)
	bookingRepo := repo.NewBookingRepo(dbx)
	bookingSvc := service.NewBookingService(bookingRepo)
	bookingHandler := handler.NewBookingHandler(bookingRepo, bookingSvc)
	resourceBookingsHandler := handler.NewResourceBookingsHandler(bookingRepo)

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

		// Создание ресурса — только для авторизованных
		r.With(handler.AuthMiddleware(authSvc)).Post("/resources", resourceHandler.Create)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)

			// защищённый роут
			r.With(handler.AuthMiddleware(authSvc)).Get("/me", authHandler.Me)
		})

		// Бронирования: только авторизованные
		r.With(handler.AuthMiddleware(authSvc)).Get("/bookings/my", bookingHandler.My)
		r.With(handler.AuthMiddleware(authSvc)).Post("/bookings", bookingHandler.Create)

		// Менеджер: смотреть ожидающие и менять статус
		r.With(handler.AuthMiddleware(authSvc)).Get("/bookings/pending", bookingHandler.Pending)

		r.With(handler.AuthMiddleware(authSvc)).Patch("/bookings/{id}/status", bookingHandler.UpdateStatus)

		r.With(handler.AuthMiddleware(authSvc)).Post("/bookings/{id}/cancel", bookingHandler.Cancel)

		r.Get("/resources/{id}/bookings", resourceBookingsHandler.List)

		r.With(handler.AuthMiddleware(authSvc)).Get("/resources/my", resourceHandler.My)

		r.With(
			handler.AuthMiddleware(authSvc),
			handler.RequireRoles(domain.RoleAdmin),
		).Post("/categories", categoryHandler.Create)
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
