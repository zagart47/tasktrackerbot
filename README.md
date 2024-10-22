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
- [To do](#to-do)
- [Команда проекта](#команда-проекта)

## Технологии
- [Golang](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/)
- [Redis](https://redis.io/)

## Использование
В файле config/config.yaml необходимо указать токен вашего телеграм-бота
```yaml
bot:
  token: <ваш токен>
```
Далее собрать контейнеры с помощью docker compose:
```powershell
docker compose up -d
```


## Разработка

### Требования
Для установки и запуска проекта необходимы golang, docker и прямые руки.

## Contributing
Если у вас есть предложения или идеи по дополнению проекта или вы нашли ошибку, то пишите мне в tg: @zagart47

## FAQ
### Зачем вы разработали этот проект?
Это проект в рамках Осеннего Мегахакатона 2024 (by SF).

## Команда проекта
- [Артур Загиров](https://t.me/zagart47) — Golang Developer

