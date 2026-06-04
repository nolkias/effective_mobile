# Инутрукция для работы с проектом  


## Подготовка проекта к работе 

- ### Создаем переменные окруджения 
```text
.env.example -> .env
```
- ### Сбираем и запускаем контейнеры 
```text
docker compose up --build -d
```

## Документация для работы с API 


## API Endpoints

Базовый URL: `http://localhost:8080/api/v1`

### 1. Создание подписки (Create)

**Endpoint:** `POST /subscriptions`

**Request body:**
```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": null
}
Поле	Тип	Обязательное	Описание
service_name	string	Да	Название сервиса
price	integer	Да	Стоимость в рублях (целое число)
user_id	string (UUID)	Да	ID пользователя в формате UUID
start_date	string	Да	Дата начала в формате MM-YYYY
end_date	string	Нет	Дата окончания в формате MM-YYYY
Response (201 Created):

json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": null,
  "created_at": "2025-06-04T10:00:00Z",
  "updated_at": "2025-06-04T10:00:00Z"
}
2. Получение подписки по ID (Read)
Endpoint: GET /subscriptions/{id}

Пример: GET /subscriptions/550e8400-e29b-41d4-a716-446655440000

Response (200 OK):

json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": null,
  "created_at": "2025-06-04T10:00:00Z",
  "updated_at": "2025-06-04T10:00:00Z"
}
Ошибки:

404 Not Found — подписка с таким ID не найдена

400 Bad Request — некорректный формат UUID

3. Обновление подписки (Update)
Endpoint: PUT /subscriptions/{id}

Пример: PUT /subscriptions/550e8400-e29b-41d4-a716-446655440000

Request body (обновляются только переданные поля):

json
{
  "price": 450,
  "end_date": "12-2026"
}
Response (200 OK):

json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "service_name": "Yandex Plus",
  "price": 450,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2026",
  "created_at": "2025-06-04T10:00:00Z",
  "updated_at": "2025-06-04T10:30:00Z"
}
Ошибки:

404 Not Found — подписка не найдена

400 Bad Request — некорректные данные

4. Удаление подписки (Delete)
Endpoint: DELETE /subscriptions/{id}

Пример: DELETE /subscriptions/550e8400-e29b-41d4-a716-446655440000

Response (204 No Content) — тело ответа пустое

Ошибки:

404 Not Found — подписка не найдена

5. Список всех подписок (List)
Endpoint: GET /subscriptions

Параметры пагинации (опционально):

Параметр	Тип	По умолчанию	Описание
limit	integer	20	Количество записей на странице
offset	integer	0	Смещение для пагинации
Примеры запросов:

text
GET /subscriptions
GET /subscriptions?limit=10&offset=0
GET /subscriptions?limit=50&offset=100
Response (200 OK):

json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "service_name": "Yandex Plus",
      "price": 400,
      "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
      "start_date": "07-2025",
      "end_date": null
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "service_name": "IVI",
      "price": 299,
      "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
      "start_date": "01-2025",
      "end_date": "06-2025"
    }
  ],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "total": 2
  }
}
6. Подсчёт стоимости подписок за период (Дополнительная ручка)
Endpoint: GET /subscriptions/total-cost

Параметры запроса:

Параметр	Тип	Обязательный	Описание
user_id	string (UUID)	Нет	Фильтр по ID пользователя
service_name	string	Нет	Фильтр по названию сервиса
start_date	string	Да	Начало периода (MM-YYYY)
end_date	string	Да	Конец периода (MM-YYYY)
Важно: Учитываются подписки, у которых период действия пересекается с запрошенным интервалом. Для бессрочных подписок (end_date = null) считается, что они действуют бесконечно.

Примеры запросов:

text
# Все подписки за период
GET /subscriptions/total-cost?start_date=01-2025&end_date=12-2025

# Подписки конкретного пользователя
GET /subscriptions/total-cost?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&start_date=01-2025&end_date=12-2025

# Подписки по названию сервиса
GET /subscriptions/total-cost?service_name=Yandex Plus&start_date=01-2025&end_date=12-2025

# Комбинированная фильтрация
GET /subscriptions/total-cost?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex Plus&start_date=01-2025&end_date=12-2025
Response (200 OK):

json
{
  "total_cost": 12500,
  "period": {
    "start_date": "01-2025",
    "end_date": "12-2025"
  },
  "filters": {
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "service_name": "Yandex Plus"
  },
  "subscriptions_count": 3
}
Поле	Описание
total_cost	Суммарная стоимость всех подписок за указанный период
period	Запрошенный период
filters	Применённые фильтры (если указаны)
subscriptions_count	Количество подписок, участвовавших в расчёте
Примеры запросов с cURL
Создание подписки
bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
Получение всех подписок
bash
curl -X GET http://localhost:8080/api/v1/subscriptions?limit=10&offset=0
Получение подписки по ID
bash
curl -X GET http://localhost:8080/api/v1/subscriptions/550e8400-e29b-41d4-a716-446655440000
Обновление подписки
bash
curl -X PUT http://localhost:8080/api/v1/subscriptions/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"price": 450}'
Удаление подписки
bash
curl -X DELETE http://localhost:8080/api/v1/subscriptions/550e8400-e29b-41d4-a716-446655440000
Подсчёт стоимости за период
bash
curl -X GET "http://localhost:8080/api/v1/subscriptions/total-cost?start_date=01-2025&end_date=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
HTTP статус-коды
Код	Описание
200	Успешный запрос (GET, PUT)
201	Успешное создание (POST)
204	Успешное удаление (DELETE)
400	Некорректный запрос (невалидные данные)
404	Ресурс не найден
405	Метод не поддерживается
500	Внутренняя ошибка сервера
text
