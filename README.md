# Merch-Market Project

### Описание
Merch Store — это веб-приложение для управления цифровыми товарами и монетами пользователей, поддерживающее регистрацию, авторизацию, покупку товаров, передачу монет другим пользователям и предоставляющее отчет о текущем состоянии аккаунта.

### Основные фичи:
- Регистрация и авторизация с использованием JWT.
- Покупка товаров из каталога с учетом баланса пользователя.
- Передача монет между пользователями.
- Просмотр информации о текущем состоянии аккаунта (баланс, инвентарь, история транзакций).

---

### Технологический стек:
- **Язык:** Go 1.24
- **Фреймворк:** Chi (REST API)
- **База данных:** PostgreSQL 15
- **Миграции:** Migrate
- **Тестирование:**
   - Юнит-тесты: testify, sqlmock
   - Интеграционные тесты: net/http/httptest
   - Нагрузочное тестирование: k6
- **Docker-деплой:** Docker Compose
- **Прочее:** bcrypt (для хэширования паролей), JWT (для аутентификации), Cleanenv (для работы с конфигурацией).

---

### Запуск проекта

1. **Убедитесь, что установлены необходимые зависимости:**
   - Docker и Docker Compose.
   - Make (для выполнения миграций через Makefile).

     ```bash
     brew install make
     ```

2. **Сконфигурируйте окружение:**

     .ENV файл был добавлен для вашего удобства 

3. **Запустите приложение с помощью Docker Compose:**
   ```bash
   docker-compose up -d
   ```

4. **Примените миграции для базы данных:**
   ```bash
   make migrate-up
   ```

5. **Проверьте состояние миграций (необязательно):**
   ```bash
   make migrate-status
   ```

6. **Остановите и откатите миграции (опционально):**
   ```bash
   make migrate-down
   ```

---

### Тестирование

#### 1. **Запуск юнит-тестов:**
Выполните все юнит-тесты с помощью `go test`:
```bash
go test ./... -coverprofile=coverage.out

go tool cover -func=coverage.out
```
![Image 15 02 25 at 20 02](https://github.com/user-attachments/assets/4296bdf5-4562-448b-a37d-cbac01e35bf1)

#### 2. **Запуск тестов на высокий уровень нагрузки:**
Для проверки производительности и поведения под нагрузкой используйте скрипт `load-test.js` с инструментом k6:
- Установите [k6](https://k6.io/).
- Выполните команду:
  ```bash
  k6 run load-test.js
  ```
- Результаты нагрузки будут содержать:
   - Достижение заданной скорости запросов.
   - Доля успешных и не успешных запросов.
   - Время ответа (95% запросов должны укладываться в 50 мс).
![photo_2025-02-15 18 49 21](https://github.com/user-attachments/assets/0781f462-3623-44f6-a316-233ebdb339d1)

---

### Использование API

#### 1. **Регистрация:**
- **Эндпоинт:** `POST /api/auth`
- **Тело запроса:**
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Пример ответа:**
  ```json
  {
    "token": "your_jwt_token"
  }
  ```

#### 2. **Покупка товара:**
- **Эндпоинт:** `GET /api/buy/{item_id}`
- **Требуется:** Заголовок `Authorization: Bearer <token>`
- **Пример ответа:**
  ```json
  {
    "message": "Item purchased successfully!"
  }
  ```

#### 3. **Передача монет:**
- **Эндпоинт:** `POST /api/sendCoin`
- **Тело запроса:**
  ```json
  {
    "toUser": "receiver_username",
    "amount": 100
  }
  ```
- **Пример ответа:**
  ```json
  {
    "message": "Coins sent successfully"
  }
  ```

#### 4. **Информация о пользователе:**
- **Эндпоинт:** `GET /api/info`
- **Требуется:** Заголовок `Authorization: Bearer <token>`
- **Пример ответа:**
  ```json
  {
    "coins": 500,
    "inventory": [
      {
        "type": "t-shirt",
        "quantity": 2
      }
    ],
    "coinHistory": {
      "received": [
        {
          "fromUser": "user1",
          "amount": 100
        }
      ],
      "sent": [
        {
          "toUser": "user2",
          "amount": 50
        }
      ]
    }
  }
  ```

---

### Результаты нагрузочного тестирования
- **Заявленные метрики:**
   - 1000 запросов в секунду с использованием 1000 виртуальных пользователей.
   - 95% запросов выполняются за время < 50 мс.
- **Результат:**
   - Успешно: ≥99,99% запросов.
   - Время ответа: Среднее < 45 мс.

---

### Информация для разработчиков

#### Установка линтера

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```
#### Запуск линтеров
Проект настроен для линтинга с помощью GolangCI-Lint:
```bash
golangci-lint run
```

#### Дополнительное тестирование эндпоинтов
Используйте Postman или аналогичные инструменты для ручного тестирования REST API.

---

### Автор
- [Alias1177](https://github.com/Alias1177) - разработчик проекта.
