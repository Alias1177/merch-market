## Запуск проекта
1. Установите Docker и Docker Compose
2. Выполните команды:
   ```bash
   # Запуск базы данных
   docker-compose up -d postgres
   
   # Применение миграций
   make migrate-up
   ```
3. Убедитесь, что порт 6000 свободен
(если нет, измените docker-compose)

