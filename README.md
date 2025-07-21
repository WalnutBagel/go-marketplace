# Go Marketplace API

Простое REST API на Go для маркетплейса с авторизацией пользователей, размещением и просмотром объявлений. Проект упакован в Docker и использует PostgreSQL как базу данных.

## Функциональность

* Регистрация и авторизация пользователей
* JWT-токены для доступа к защищённым ресурсам
* CRUD для объявлений
* Пагинация, сортировка и фильтрация ленты
* Признак принадлежности объявления текущему юзеру

## Стек

* Go 1.24
* PostgreSQL 15
* GORM
* Docker / Docker Compose

## Запуск

1. Склонируйте репозиторий:

   ```bash
   git clone https://github.com/WalnutBagel/go-marketplace.git
   cd go-marketplace
   ```
2. Соберите и запустите:

   ```bash
   docker-compose up --build
   ```

   После этого API будет доступен на [http://localhost:8080](http://localhost:8080)

## Описание ендпоинтов

### 🔑 Регистрация

* `POST /register`
* Входной JSON:

  ```json
  {
    "username": "user",
    "password": "securepass"
  }
  ```
* Ответ: 201 Created

### 🌐 Авторизация

* `POST /login`
* Входной JSON:

  ```json
  {
    "username": "user",
    "password": "securepass"
  }
  ```
* Успех: 200 OK

  ```json
  {
    "token": "JWT"
  }
  ```

### 📅 Лента объявлений

* `GET /ads`
* Headers: `Authorization: Bearer <token>`
* Query-параметры:

  * `page` — страница
  * `sort_by=date|price`
  * `order=asc|desc`
  * `min_price`, `max_price`
* Респонс:

  ```json
  [
    {
      "id": 1,
      "title": "Товар",
      "description": "Описание",
      "image_url": "http://...",
      "price": 100,
      "owner": "user1",
      "is_owner": true
    }
  ]
  ```

### ✉️ Создание объявления

* `POST /ads`
* Headers: `Authorization: Bearer <token>`
* Входной JSON:

  ```json
  {
    "title": "Товар",
    "description": "Описание",
    "image_url": "http://...",
    "price": 500
  }
  ```

### ✏️ Редактирование

* `PUT /ads/{id}`
* Headers: `Authorization: Bearer <token>`
* Вход: JSON с полями для обновления

### ❌ Удаление

* `DELETE /ads/{id}`
* Headers: `Authorization: Bearer <token>`

## Тестирование

Есть файл `test_request.http` с полным набором запросов для VS Code REST Client (GET, POST, PUT, DELETE с token-ом и без).

## Заметки

* Все поля валидируются: логин, пароль, цена, длина описания
* Пароли хешируются
* Миграция в базе происходит автоматически

---

Сделано с нуля для тестового задания на Go.
