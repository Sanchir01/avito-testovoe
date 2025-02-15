# 🚀 Avito testovoe

Этот проект представляет собой серверное приложение на **Go**, использующее **PostgreSQL** в качестве базы данных и **Chi Router** для маршрутизации API-запросов.

## 📋 Предварительные требования

Перед началом работы убедитесь, что у вас установлены:
- **Docker** и **Docker Compose**
- **Make** (для удобного запуска команд)
- **Goose** (для выполнения миграций)
- **Golangci-lint** (если планируете запускать локально без контейнеров)
- **Go** (если планируете запускать локально без контейнеров)

## 🔧 Установка и запуск


### 1️⃣  Применение миграций базы данных

Перед запуском сервера выполните миграции, чтобы создать необходимые таблицы:
```sh
  make migrations-up
```

### 2️⃣ Запуск контейнеров

Можно использовать контейнеры (PostgreSQL и приложение):
```sh
  make docker
```

### 3️ Запуск приложения

После миграций можно запустить сервер локально:
```sh
  make run
```
### 4️⃣ Запуск тестов

После миграций можно запустить тесты:
```sh
  make test
```
Сервер будет доступен по адресу: [http://localhost:8080](http://localhost:8080)

## 📂 Структура проекта
```
.
├── cmd/                 # Основной код сервера
│   ├── main/            # Точка входа в приложение
│   ├── seed/            # Код для начального заполнения продуктами
├── internal/            # Внутренние пакеты
│   ├── app/             # Обработчики HTTP-запросов
│   ├── config/          # Описание конфига и его чтение
│   ├── context/         # Контексты для получения значения из контекста
│   ├── feature/         # Repository, service и handlers для всех разделов
│   ├── servers/         # Логика создания HTTP-сервера и все endpoints, middleware
├── migrations/          # Файлы миграций базы данных
├── .env                 # Файл переменных окружения (скрыт в .gitignore)
├── docker-compose.yaml  # Конфигурация Docker Compose
├── golangci.yaml        # Конфигурация Linter
├── Dockerfile           # Конфигурация Docker
├── Makefile             # Упрощённые команды для запуска
├── README.md            # Этот файл 🙂
├── insomnia.json        # Файл конфигурации для Insomnia
└── postman.json         # Файл конфигурации для Postman
```


## ⚙️ Переменные окружения

Перед запуском убедитесь, что у вас есть `.env` файл с необходимыми параметрами:

```env
POSTGRES_PASSWORD="avito"
CONFIG_PATH="config/config.yaml"
JWT_SECRET="asdadsada&^&(as"
```

## 🛠️ Полезные команды

### 🔍 Проверка логов
Если нужно посмотреть логи бд:
```sh
  docker logs -f avito-testovoe-db-1
```

Если нужно посмотреть логи app:
```sh
  docker logs -f avito-testovoe-app-1
```
