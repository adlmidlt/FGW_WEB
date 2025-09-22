# 1. Команды для докера
<code>docker-compose up -d --build</code> - пересоборка всех сервисов.<br>
<code>docker-compose up -d --build <service_name></code> - пересборка конкретного сервиса.<br>
<code>docker-compose build</code> - пересборка без запуска.<br>
<code>docker-compose up -d --build <your_app_service_name></code> - пересборка только Go приложения.<br>
<code>docker-compose up -d --build app</code> - пересобрать только приложение.<br>
<code>docker builder prune</code> - удалить слои, которые не используются образами.<br>