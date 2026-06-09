package main

import (
	"context"
	"effective_mobile/internal/database"
	"effective_mobile/internal/handlers"
	"effective_mobile/internal/repository"
	"effective_mobile/internal/router"
	"effective_mobile/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Subscription Service API
// @version 1.0
// @description API для управления подписками пользователей
// @host localhost:8080
// @BasePath /api/v1
func main() {
	db, err := database.DbConnection()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	database.RunMigrations(db, "migrations")

	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionHandlers := handlers.NewSubscriptionHandler(subscriptionService)

	r := router.SetupRouter(subscriptionHandlers)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Запускаем сервер в горутине, чтобы не блокировать обработку сигналов
	go func() {
		log.Println("Сервер запущен на порту 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Ошибка запуска сервера:", err)
		}
	}()

	// Ждём сигнала завершения (Ctrl+C или SIGTERM от Docker/k8s)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Завершение работы сервера...")

	// Даём активным запросам 5 секунд на завершение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Ошибка при завершении сервера:", err)
	}

	log.Println("Сервер остановлен")
}
