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


### 1️⃣   Запуск контейнеров

Перед запуском сервера выполните миграции, чтобы создать необходимые таблицы:
```sh
  make migrations-up
  
  make seed
 
```

### 2️⃣ Применение миграций базы данных и сидинг продуктов

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
├── pkg/                 # Папка с утилитами 
├── test/                # Файлы с тестами
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

## 🧪 Unit тесты

1. Тесты сервиса продуктов находяться в internal/feature/product/mocks
2. Тесты сервиса пользователя находяться в internal/feature/user/mocks


## 📮 Тестирование API

Для удобного тестирования API в проект включены готовые конфигурации:

### 🐦 Insomnia
1. Скачайте и установите [Insomnia](https://insomnia.rest/)
2. Импортируйте конфиг:
    - **Settings** → **Data** → **Import Data** → **From File**
    - Выберите файл `insomnia.json` из корня проекта
3. Все запросы будут доступны в коллекции `Avito`

### 📮 Postman
1. Скачайте и установите [Postman](https://www.postman.com/)
2. Импортируйте конфиг:
    - Нажмите **Import** в левом верхнем углу
    - Выберите файл `postman.json` из корня проекта
3. Коллекция `Avito` появится в вашем рабочем пространстве

### ⚠️ Перед использованием:
1. Убедитесь, что сервер запущен (`make docker` или `make run`)
2. Проверьте, что адрес сервера в запросах указан как: "localhost:8080"
3. Для авторизованных запросов используйте JWT-токен

## 🛠️ Проблемы и их решения

### Проблема 1: Тестирование покупки товара

Для покупки товара нужно было получить айди и при вызове endpoint-а `/api/buy/{item}`  возникали ошибки, если сервер возвращал статус,
отличный от 200 OK. Изначально я сразу выполнял проверку с помощью 
`.Status(http.StatusOK)`, что приводило к немедленному завершению теста при ошибке. Для получения айди товара я реализовал отдельный эндпоит для получения всех товаров
 и из всех товаров я выбирал первый и брал его айди

**Решение:**  
Вместо немедленной проверки, я сохранил ответ в переменную и вручную проверял статус. Если статус был не 200, извлекал сообщение об ошибке из JSON-ответа и выводил его с помощью `t.Errorf()`. Пример:

```go
resp := e.GET("/api/buy/{item}", firstProductID).
    WithHeader("Authorization", "Bearer "+token).
    Expect()

statusCode := resp.Raw().StatusCode
if statusCode != http.StatusOK {
    errObj := resp.JSON().Object()
    errMsg := errObj.Value("error").String().Raw()
    t.Errorf("Ошибка при покупке продукта. Статус: %d, сообщение: %s", statusCode, errMsg)
} else {
    resp.Status(http.StatusOK)
    // Дополнительные проверки успешного ответа
}
```

### Проблема 2: Получение списка продуктов

Изначально я ожидал, что ответ от endpoint-а `/api/products` будет массивом JSON, но фактический ответ имел структуру:

```json
{
  "products": [ ... ],
  "status": "OK"
}
```

**Решение:**  
Был изменён тест так, чтобы сначала получался объект, а затем извлекался массив продуктов из поля `products`. Пример:

```go
obj := e.GET("/api/products").
    WithHeader("Authorization", "Bearer "+token).
    Expect().
    Status(http.StatusOK).
    JSON().
    Object()

products := obj.Value("products").Array()
firstProductID := products.Value(0).Object().Value("ID").String().Raw()
```

```go
obj := e.GET("/api/products").
Expect().
Status(http.StatusOK).
JSON().
Object().ContainsKey("products")

products := obj.Value("products").Array()

firstProductID := products.Value(0).Object().Value("ID").String().Raw()
```



