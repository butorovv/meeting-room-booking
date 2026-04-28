// @title Room Booking Service API
// @version 1.0
// @description Сервис бронирования переговорок
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/butorovv/meeting-room-booking/config"
	_ "github.com/butorovv/meeting-room-booking/docs"
	deliveryHttp "github.com/butorovv/meeting-room-booking/internal/delivery/http"
	middleware "github.com/butorovv/meeting-room-booking/internal/delivery/middlewares"
	"github.com/butorovv/meeting-room-booking/internal/repository"
	"github.com/butorovv/meeting-room-booking/internal/usecase"
	"github.com/butorovv/meeting-room-booking/pkg/cache"
	"github.com/butorovv/meeting-room-booking/pkg/logger"
	"github.com/butorovv/meeting-room-booking/pkg/worker"
	"github.com/jackc/pgx/v5/pgxpool"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg := config.Load()
	logger.Global().Info("Starting server", "port", cfg.AppPort)

	dbPool, err := pgxpool.New(context.Background(), cfg.DBPath())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	userRepo := repository.NewUserRepository(dbPool)
	roomRepo := repository.NewRoomRepository(dbPool)
	scheduleRepo := repository.NewScheduleRepository(dbPool)
	bookingRepo := repository.NewBookingRepository(dbPool)

	var slotsCache cache.Cache
	if cfg.RedisAddr != "" {
		slotsCache = cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)

		pingCtx, cancel := context.WithTimeout(context.Background(), time.Second)
		if err := slotsCache.Ping(pingCtx); err != nil {
			logger.Global().Warn("Redis unavailable", "addr", cfg.RedisAddr, "err", err)
		}
		cancel()
	}

	notificationWorkers := worker.NewWorkerPool(2, 100)

	// Инициализация usecase
	authUC := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret)
	roomUC := usecase.NewRoomUseCase(roomRepo)
	scheduleUC := usecase.NewScheduleUseCase(scheduleRepo)
	slotUC := usecase.NewSlotUseCase(scheduleRepo, bookingRepo, slotsCache)
	bookingUC := usecase.NewBookingUseCase(
		dbPool,
		bookingRepo,
		roomRepo,
		scheduleRepo,
		slotsCache,
		notificationWorkers,
	)

	// Инициализация handlers
	authHandler := deliveryHttp.NewAuthHandler(authUC)
	roomHandler := deliveryHttp.NewRoomHandler(roomUC)
	scheduleHandler := deliveryHttp.NewScheduleHandler(scheduleUC)
	slotHandler := deliveryHttp.NewSlotHandler(slotUC)
	bookingHandler := deliveryHttp.NewBookingHandler(bookingUC)
	healthHandler := deliveryHttp.NewHealthHandler()

	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /dummyLogin", authHandler.DummyLogin)
	mux.HandleFunc("GET /_info", healthHandler.Info)
	mux.HandleFunc("GET /health", healthHandler.Health)

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	mux.Handle("GET /rooms/list", authMiddleware(http.HandlerFunc(roomHandler.ListRooms)))
	mux.Handle("POST /rooms/create", authMiddleware(middleware.AdminOnly(http.HandlerFunc(roomHandler.CreateRoom))))
	mux.Handle("POST /rooms/{roomId}/schedule/create", authMiddleware(middleware.AdminOnly(http.HandlerFunc(scheduleHandler.CreateSchedule))))
	mux.Handle("GET /rooms/{roomId}/slots/list", authMiddleware(http.HandlerFunc(slotHandler.GetAvailableSlots)))
	mux.Handle("POST /bookings/create", authMiddleware(middleware.UserOnly(http.HandlerFunc(bookingHandler.CreateBooking))))
	mux.Handle("POST /bookings/{bookingId}/cancel", authMiddleware(middleware.UserOnly(http.HandlerFunc(bookingHandler.CancelBooking))))
	mux.Handle("GET /bookings/my", authMiddleware(middleware.UserOnly(http.HandlerFunc(bookingHandler.MyBookings))))
	mux.Handle("GET /bookings/list", authMiddleware(middleware.AdminOnly(http.HandlerFunc(bookingHandler.ListBookings))))

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: mux,
	}

	go func() {
		logger.Global().Info("Server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Global().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Global().Error("Server shutdown failed", "err", err)
	}

	logger.Global().Info("Shutting down worker pool...")
	notificationWorkers.Shutdown()

	if slotsCache != nil {
		logger.Global().Info("Closing redis cache...")
		if err := slotsCache.Close(); err != nil {
			logger.Global().Error("Redis cache close failed", "err", err)
		}
	}

	logger.Global().Info("Closing database connection...")
	dbPool.Close()

	logger.Global().Info("Server gracefully stopped")
}
