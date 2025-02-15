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



5. Результат высоконагрузочного тестирования
![photo_2025-02-15 18 49 21](https://github.com/user-attachments/assets/a482c0de-042e-4954-8117-2b66256a98c1)
