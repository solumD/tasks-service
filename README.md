# tasks-service
Сервис для хранения и отслеживания задач.

Версия Go - 1.25.5

В проекте использовалась только стандартная библиотека.

## Установка проекта
```bash
 git clone github.com/solumD/tasks-service
 cd tasks-service/
 go mod tidy
```

## Запуск
Поменять значения в .env-файле, если необходимо.
```dotenv
  HTTP_SERVER_HOST=0.0.0.0
  HTTP_SERVER_PORT=8080

  #debug, info, warn, error
  LOGGER_LEVEL=info 
```

Для запуска локально выполнить в терминале команду. При вводе команды проект собирается и запускатеся локально.
```bash
  make run-locally
```

Для запуска в контейнере (необходим Docker) выполнить в терминале команду. При вводе команды собирается docker-образ и поднимается docker-контейнере. По умолчанию пробрасывается порт 8080. 
```bash
  make run
```

## Тестирование
Для запуска unit-тестов выполнить в терминале команду.
```bash
  make test
```

## Эндпоинты
### POST /todos - создание задачи

Тело запроса:
```
{
  "title": "string",
  "description": "string",
  "done": false
}
```
Тело успешного ответа:
```
{
  "id": 1
}
```

### GET /todos - получение списка всех задач

Тело запроса: отсутствует

Тело успешного ответа:
```
{
    "todos": [
        {
            "id": 1,
            "title": "string",
            "description": "string",
            "done": false
        },
        {
            "id": 2,
            "title": "string",
            "description": "string",
            "done": false
        }
    ]
}
```

### GET /todos/{id} - получение задачи по id

Тело запроса: отсутствует

Тело успешного ответа:
```
{
  "id": 1,
  "title": "string",
  "description": "string",
  "done": false
}
```

### PUT /todos/{id} - обновление информации о задаче по id (полностью меняет информацию, потому что это не PATCH)

Тело запроса:
```
{
  "title": "string",
  "description": "string",
  "done": false
}
```
Тело успешного ответа: отсутствует

### DELETE /todos/{id} - удаление задачи по id

Тело запроса: отсутствует

Тело успешного ответа: отсутствует



