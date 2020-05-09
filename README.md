# test-avito-antibot
Тестовое задание для стажировки в Авито

Добавлен docker для поднятия сервиса на 500 порту (и добавлен docker-compose)

Есть handler для сброса : URL "/reset" с методом PUT

В случае успеха - 201

В случе ошибки - 409 (неверный формат IP / пустые данные в Хедере)

Для регистрации IP - URL : "/" с методом POST

В случае пустого значения в Хедере или неправильных данных отсылается ошибка 409

В случае кулдауна - 429 и в Хедере Retry-After время в секундах для ожидания

В случе успеха - 200

Добавлены тесты - покрыт весь функционал, не считай main.go файла

Используемая база данных - PostgreSQL

Для изменения конфигурации - файл project.env

LIMIT - количество

MASK - число, без "/"

COOLDOWN - число кулдауна в минутах

PORT - порт на котором поднимается сервис (по умолчанию 5000)

