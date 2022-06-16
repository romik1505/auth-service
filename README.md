# Auth Service

Часть сервиса аутентификации

**Запуск сервера**

`make run`


**REST  маршруты**

1. Выдача пары токенов 

Метод *Post* localhost:8080/login 

*Формат входных данных:* 
`{"user_id":"user1"}`

2. Обновление пары токенов

Метод *Post* localhost:8080/refresh-token

*Формат входных данных:*
`{"access_token":"token1","refresh_token":"token2"}`