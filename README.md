![Static Badge](https://img.shields.io/badge/%D1%81%D1%82%D0%B0%D1%82%D1%83%D1%81-%D0%B3%D0%BE%D1%82%D0%BE%D0%B2-blue)
![Static Badge](https://img.shields.io/badge/GO-1.23-blue)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/zagart47/tasktrackerbot)
![GitHub last commit (by committer)](https://img.shields.io/github/last-commit/zagart47/tasktrackerbot)
![GitHub forks](https://img.shields.io/github/forks/zagart47/tasktrackerbot)

# TaskTrackerBot 
Бот для телеграм который напоминает о задачах

## Содержание
- [Технологии](#технологии)
- [Использование](#использование)
- [Разработка](#разработка)
- [Contributing](#contributing)
- [FAQ](#faq)
- [Команда проекта](#команда-проекта)

## Технологии
- [Golang](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)
- [Redis](https://redis.io/)

## Использование
В файле config/config.yaml необходимо указать токен вашего телеграм-бота:
```yaml
bot:
  token: ваш токен
```
Далее собрать контейнеры с помощью docker compose:
```powershell
docker compose up -d
```

Чтобы бот напоминал о задачах, а делает он это только в личку, необходимо перейти на бота и нажать `Запустить`.

Бота необходимо добавить в вашу группу. Он запоминает последние сообщения от пользователей в этой группе.
После того, как вы отправили напоминание в чат, необходимо тегнуть бота в чате и написать ему команду:
Например ```@имя_бота ctrl 5h```, где `имя_бота` это имя бота, которое вы дали в BotFather, `ctrl` - это команда взятия 
задачи на контроль, `5h` - это значит, что бот должен напомнить о задаче через 5 часов. 

После отправки команды, бот запомнит ваше последнее сообщение
(то, которое вы отправили в группу/чат до команды).

Бот поддерживает интервалы типа `h` - часы, `d` - дни, `w` - недели и `m` - месяцы.

Для дебага также есть поддержка секунд - `s`.

## Разработка

### Требования
Для установки и запуска проекта необходимы golang и docker.

## Contributing
Если у вас есть предложения или идеи по дополнению проекта или вы нашли ошибку, то пишите мне в tg: @zagart47

## FAQ
### Зачем вы разработали этот проект?
Это проект в рамках Осеннего Мегахакатона 2024 (by SF).

## Команда проекта
- [Артур Загиров](https://t.me/zagart47) — Golang Developer

