# Mattermost Vote Bot
## Бот для создания голосований внутри чатов мессенджера Mattermost


## Установка
1. Склонируйте репозиторий:
   ```bash
   git clone git@github.com:gfarida/mattermost-vote-bot.git
   ```

2. Создайте бота в Mattermost:

    a. Откройте http://localhost:8065.

    b. Перейдите: System Console → Integrations → Bot Accounts → Create.

    c. Скопируйте токен и вставьте его в "MATTERMOST_TOKEN" в docker-compose.yaml

2. Запустите систему:

``` bash
    docker-compose up --build -d
```

## Команды

1. Создать голосование:

```bash
/vote create <Вопрос?> <Вариант1,Вариант2>
```

2. Проголосовать:

```bash
/vote vote <ID> <номер_варианта>
```

3. Посмотреть результаты опроса:

```bash
/vote results <ID>
```

4. Завершить голосование:

```bash
/vote end <ID>
```

5. Удалить голосование:

```bash
/vote delete <ID>
```