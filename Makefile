.PHONY: up down build seed test e2e clean mocks test-cover cover filter-cover e2e-positive e2e-negative

# Переменные
DOCKER_COMPOSE = docker-compose
APP_NAME = github.com/butorovv/meeting-room-booking
COVERFILE = coverage.out
FILTERED_COVERFILE = coverage.filtered.out
SRC_DIRS = ./...
TEST_FLAGS = -covermode=atomic -coverprofile=$(COVERFILE)

# Запуск проекта
up:
	$(DOCKER_COMPOSE) up --build

# Запуск в фоне
up-d:
	$(DOCKER_COMPOSE) up -d --build

# Остановка
down:
	$(DOCKER_COMPOSE) down

# Остановка с удалением данных
down-v:
	$(DOCKER_COMPOSE) down -v

# Пересборка
build:
	$(DOCKER_COMPOSE) build

# Логи
logs:
	$(DOCKER_COMPOSE) logs -f

# Заполнение БД тестовыми данными
seed:
	@echo "Seeding database..."
	docker exec -i room_booking_db psql -U postgres -d room_booking < ./scripts/seed.sql
	@echo "Seed completed"

# Генерация моков
mocks:
	@echo "======== Creating mocks... ========"
	@echo "======== REPOSITORY MOCKS (interfaces from usecase) ========"
	mockgen -source=internal/usecase/booking_usecase.go -destination=internal/usecase/mock/booking_repo_mock.go -package=mock
	mockgen -source=internal/usecase/room_usecase.go -destination=internal/usecase/mock/room_repo_mock.go -package=mock
	mockgen -source=internal/usecase/schedule_usecase.go -destination=internal/usecase/mock/schedule_repo_mock.go -package=mock
	mockgen -source=internal/usecase/interfaces.go -destination=internal/usecase/mock/pgxpool_mock.go -package=mock
	mockgen -source=internal/usecase/auth_usecase.go -destination=internal/usecase/mock/auth_repo_mock.go -package=mock
	mockgen -source=internal/usecase/tx_interface.go -destination=internal/usecase/mock/tx_mock.go -package=mock
	@echo "======== USECASE MOCKS (interfaces from delivery/http) ========"
	mockgen -source=internal/delivery/http/booking_handler.go -destination=internal/delivery/http/mock/booking_usecase_mock.go -package=mock
	mockgen -source=internal/delivery/http/room_handler.go -destination=internal/delivery/http/mock/room_usecase_mock.go -package=mock
	mockgen -source=internal/delivery/http/schedule_handler.go -destination=internal/delivery/http/mock/schedule_usecase_mock.go -package=mock
	mockgen -source=internal/delivery/http/slot_handler.go -destination=internal/delivery/http/mock/slot_usecase_mock.go -package=mock
	mockgen -source=internal/delivery/http/auth_handler.go -destination=internal/delivery/http/mock/auth_usecase_mock.go -package=mock
	@echo "======== Mocks created ========"

# Запуск всех тестов
test:
	go test -v -cover ./...

# Запуск тестов с покрытием
test-cover: mocks cover

# Фильтрация покрытия (убираем моки и тесты)
filter-cover:
	@grep -vE "mock_|_test.go" $(COVERFILE) > $(FILTERED_COVERFILE)

# Покрытие
cover: test filter-cover
	@echo ""
	@echo "======== Coverage ========"
	@go tool cover -func=$(FILTERED_COVERFILE) | grep total

# Запуск E2E тестов (основной скрипт)
e2e:
	@echo "======== Running all E2E tests ========"
	chmod +x ./scripts/e2e_test1.sh
	bash ./scripts/e2e_test1.sh

# Запуск только позитивных E2E тестов
e2e-positive:
	@echo "======== Running positive E2E tests ========"
	chmod +x ./scripts/e2e_test_positive.sh
	bash ./scripts/e2e_test_positive.sh

# Запуск только негативных E2E тестов
e2e-negative:
	@echo "======== Running negative E2E tests ========"
	chmod +x ./scripts/e2e_test_negative.sh
	bash ./scripts/e2e_test_negative.sh

# Очистка
clean:
	$(DOCKER_COMPOSE) down -v
	rm -rf ./tmp
	go clean -testcache
	rm -f $(COVERFILE) $(FILTERED_COVERFILE)

# Помощь
help:
	@echo "Available commands:"
	@echo "  make up           - Start project"
	@echo "  make down         - Stop project"
	@echo "  make seed         - Fill database with test data"
	@echo "  make mocks        - Generate mocks for tests"
	@echo "  make test         - Run unit tests"
	@echo "  make test-cover   - Run tests with coverage (>40%)"
	@echo "  make e2e          - Run all E2E tests"
	@echo "  make e2e-positive - Run positive E2E tests only"
	@echo "  make e2e-negative - Run negative E2E tests only"
	@echo "  make clean        - Clean everything"