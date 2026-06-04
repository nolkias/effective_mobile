package main

import (
	"effective_mobile/internal/database"
	"effective_mobile/internal/handlers"
	"effective_mobile/internal/repository"
	"effective_mobile/internal/router"
	"effective_mobile/internal/service"
	"log"
)

func main() {
	// Создаем Подключение к бд
	db, err := database.DbConnection()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()
	// Миграции
	database.RunMigrations(db, "migrations")

	// Можно было все в одном main, но рука не поднимается не делить на слои
	// Пробрасываем репозиторий
	subscriptionrepo := repository.NewSubscriptionRepository(db)
	// Пробрасываем сервис
	subscriptionService := service.NewSubscriptionService(subscriptionrepo)
	// Пробрасываем хэндлер
	subscriptionHandlers := handlers.NewSubscriptionHandler(subscriptionService)

	// Роуты
	r := router.SetupRouter(subscriptionHandlers)

	log.Println("Сервер запущен на порту 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
