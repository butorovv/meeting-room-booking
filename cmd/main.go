package main

import (
	"log"
	"net/http"

	"github.com/butorovv/meeting-room-booking/config"
	deliveryHttp "github.com/butorovv/meeting-room-booking/internal/delivery/http"
	"github.com/butorovv/meeting-room-booking/internal/repository"
	"github.com/butorovv/meeting-room-booking/internal/usecase"
	"github.com/butorovv/meeting-room-booking/pkg/logger"
)

func main() {
	cfg := config.Load()

	// Подключение к БД
	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	logger.Global().Info("Connected to database")

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)

	// Инициализация usecase
	authUC := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret)

	// Инициализация handler
	authHandler := deliveryHttp.NewAuthHandler(authUC)

	// Регистрация роутов
	mux := http.NewServeMux()
	mux.HandleFunc("POST /dummyLogin", authHandler.DummyLogin)
	mux.HandleFunc("GET /_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Запуск сервера
	addr := ":" + cfg.AppPort
	logger.Global().Info("Server started", "addr", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
