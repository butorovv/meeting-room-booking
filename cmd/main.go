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
	"log"
	"net/http"

	"github.com/butorovv/meeting-room-booking/config"
	_ "github.com/butorovv/meeting-room-booking/docs"
	deliveryHttp "github.com/butorovv/meeting-room-booking/internal/delivery/http"
	middleware "github.com/butorovv/meeting-room-booking/internal/delivery/middlewares"
	"github.com/butorovv/meeting-room-booking/internal/repository"
	"github.com/butorovv/meeting-room-booking/internal/usecase"
	"github.com/butorovv/meeting-room-booking/pkg/logger"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg := config.Load()
	logger.Global().Info("Starting server", "port", cfg.AppPort)

	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	slotRepo := repository.NewSlotRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// Инициализация usecase
	authUC := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret)
	roomUC := usecase.NewRoomUseCase(roomRepo)
	scheduleUC := usecase.NewScheduleUseCase(scheduleRepo, slotRepo)
	slotUC := usecase.NewSlotUseCase(slotRepo)
	bookingUC := usecase.NewBookingUseCase(bookingRepo, slotRepo)

	// Инициализация handlers
	authHandler := deliveryHttp.NewAuthHandler(authUC)
	roomHandler := deliveryHttp.NewRoomHandler(roomUC)
	scheduleHandler := deliveryHttp.NewScheduleHandler(scheduleUC)
	slotHandler := deliveryHttp.NewSlotHandler(slotUC)
	bookingHandler := deliveryHttp.NewBookingHandler(bookingUC)

	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /dummyLogin", authHandler.DummyLogin)
	mux.HandleFunc("GET /_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	mux.Handle("GET /rooms/list", authMiddleware(http.HandlerFunc(roomHandler.ListRooms)))
	mux.Handle("POST /rooms/create", authMiddleware(middleware.AdminOnly(http.HandlerFunc(roomHandler.CreateRoom))))
	mux.Handle("POST /rooms/{roomId}/schedule/create", authMiddleware(middleware.AdminOnly(http.HandlerFunc(scheduleHandler.CreateSchedule))))
	mux.Handle("GET /rooms/{roomId}/slots/list", authMiddleware(http.HandlerFunc(slotHandler.GetAvailableSlots)))
	mux.Handle("POST /bookings/create", authMiddleware(http.HandlerFunc(bookingHandler.CreateBooking)))
	mux.Handle("GET /bookings/my", authMiddleware(http.HandlerFunc(bookingHandler.MyBookings)))
	mux.Handle("POST /bookings/{bookingId}/cancel", authMiddleware(http.HandlerFunc(bookingHandler.CancelBooking)))
	mux.Handle("GET /bookings/list", authMiddleware(middleware.AdminOnly(http.HandlerFunc(bookingHandler.ListBookings))))

	logger.Global().Info("Server started", "addr", ":"+cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, mux))
}
