# WB Task L0 - Сервис обработки заказов

Микросервис для обработки и отображения данных о заказах с использованием Kafka, PostgreSQL, Nginx и Redis.

## Функциональность

- **Балансировка нагрузки** с помощью Nginx
- **Получение заказов** из Kafka в реальном времени
- **Сохранение данных** в PostgreSQL с транзакционной безопасностью
- **Redis-кэширование** для быстрого доступа к данным в распределенной среде
- **REST API** для получения информации о заказах
- **Веб-интерфейс** для поиска заказов по ID
- **Автовосстановление кэша** при перезапуске сервиса

## Установка и запуск

### Быстрый старт

```bash
# Клонирование репозитория
git clone https://github.com/BobrAll/wb-task-L0
cd wb-task-L0

# Запуск микросервиса
./scripts/start_project.sh

# Или вручную:
docker-compose -f deployments/docker-compose.yaml up -d
PORT=8081 go run cmd/server/main.go &
PORT=8082 go run cmd/server/main.go
```

### Скрипты

В проекте доступны удобные скрипты:

- **`scripts/start_project.sh`** - запуск всего проекта (Docker + микросервис в нескольких экземплярах)
- **`scripts/run_tests.sh`** - запуск всех тестов (unit + integration)
- **`scripts/kafka-producer.sh`** - скрипт для отправки тестовых заказов в Kafka

## Конфигурация

Сервис настраивается через переменные окружения в `./configs/.env`.  
Из этого файла конфиг подтягивается как в приложение, так и в докер. Пример есть в .env-example

## API Endpoints

### Получить заказ по ID

```http
GET /order/{order_uid}
```

### Получить список ID заказов

```http
GET /orders?search=test&page=0&size=10
```

Подробная документация API доступна в /api/swagger.yaml

## Тестирование

```bash
# Все тесты
./scripts/run_tests.sh

# Только unit-тесты
go test ./test/unit/ -v

# Только интеграционные тесты
go test ./test/integration/ -v -timeout 30s
```

## Структура БД

Сервис использует 4 основные таблицы:
- `orders` - основная информация о заказах
- `deliveries` - данные доставки
- `payments` - информация об оплате
- `items` - товары в заказе

Файлы миграции находятся в директории `migrations/`.