# Simple Web Games Service


## Describe
Сервис предоставляет возможность создать сессию игры крестики нолики с указанной размерностью поля (по умолчанию 3x3, максимум 10x10).
При создании игры пользователю отправляется Cookie, в котором содержится `game_token`, необходимый для аутентификации пользователя на сервере.
При подключении пользователя к игре, ему так же возвращается Cookie, в котором содержится `game_token`.
Далее игроки совершают ходы и по окончании игры, игровое поле удаляется и их сессия завершается.

## Routers
* `POST /create body:{size[optional]=3}` - создание новой сессии с указанием размера поля (3x3, 4x4 ...), возвращает cookie с параметром game_token
* `GET /list` - возвращает список доступных к подключению сессий и название игры
* `GET /join/{session_id}` - подключение к игре по токену сессии, возвращает game_id (url)
* `POST /game/{game_id} body:{row=1, col=1}` - процесс игры передаеётся в body запроса и каждый раз возвращается ответ со статусом игры и состоянием игрового поля
* `GET /swagger.html` - документация swagger

## Commands
* `make env` - создание конфигурационного файла .env из данных в .env.sample
* `make build` - сборка сервиса
* `make run` - запусе сервиса
* `make unit-test` - запуск юнит-тестов
* `make swagger` - создание swagger документации
* `make df` - создание Dockerfile с сервисом
* `make service_up` - запуск Docker-compose с заранее собранным в Dockerfile сервисом
