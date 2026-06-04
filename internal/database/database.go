package database

import (
	"database/sql"
	"effective_mobile/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

func DbConnection() (*sql.DB, error) {

	configDb := config.GetDBConfig()

	connStr :=
		"host=" + configDb.DBHost +
			" port=" + configDb.DBPort +
			" user=" + configDb.DBUser +
			" password=" + configDb.DBPassword +
			" dbname=" + configDb.DBName +
			" sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	if err = db.Ping(); err != nil {
		log.Fatal("БД не отвечает:", err)
	}

	// Возвращаем подключение
	return db, nil
}

func RunMigrations(db *sql.DB, migrationsPath string) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("Не удалось создать driver для миграций: %v", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver,
	)
	if err != nil {
		log.Printf("Не удалось инициализировать миграции: %v", err)
		return
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Миграции не требуются — таблицы уже актуальны")
			return
		}
		log.Printf("Ошибка при выполнении миграций: %v", err)
		return
	}

	log.Println("Миграции успешно применены")
}
