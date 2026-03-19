package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TI2022/next-go-supabase-base/app-service/internal/infrastructure/database"
	"github.com/TI2022/next-go-supabase-base/app-service/internal/interfaces/handler"
)

func main() {
	// DB 接続
	db, err := database.OpenFromEnv()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	// ハンドラ初期化
	mux := http.NewServeMux()
	mux.Handle("/health", handler.NewHealthHandler())
	mux.Handle("/auth/login", handler.NewLoginHandler(db))
	mux.Handle("/me", handler.NewMeHandler(db))

	addr := os.Getenv("APP_SERVICE_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("app-service listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

