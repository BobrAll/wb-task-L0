# WB Task L0 - Сервис обработки заказов

Микросервис для обработки и отображения данных о заказах с использованием Kafka, PostgreSQL и in-memory кэша.

## Функциональность

- **Получение заказов** из Kafka в реальном времени
- **Сохранение данных** в PostgreSQL с транзакционной безопасностью
- **In-memory кэширование** для быстрого доступа к данным
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
go run cmd/main/main.go
```

### Скрипты

В проекте доступны удобные скрипты:

- **`scripts/start_project.sh`** - запуск всего проекта (Docker + приложение)
- **`scripts/run_tests.sh`** - запуск всех тестов (unit + integration)
- **`scripts/kafka-producer.sh`** - скрипт для отправки тестовых заказов в Kafka

## Конфигурация

Сервис настраивается через переменные окружения в `configs/.env`:

```env
POSTGRES_DB: wb-db
POSTGRES_USER: wb-user
POSTGRES_PASSWORD: wb-password
POSTGRES_PORT: 5432
POSTGRES_HOST: localhost
POSTGRES_SSL_MODE: disable
```
Из этого файла конфиг подтягивается как в приложение, так и в докер.

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