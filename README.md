## Запуск проекта
1. Установите Docker и Docker Compose
2. Выполните команды:
   ```bash
   # Запуск Приложения
   docker-compose up -d
   
   # Применение миграций
   make migrate-up
   make migrate-status
   make migrate-down
   ```
3. Убедитесь, что порт 6000 свободен
(если нет, измените docker-compose)

4. Результат высоконагрузочного тестирования
