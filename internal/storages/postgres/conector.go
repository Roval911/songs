package postgres

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
)

var db *sql.DB

type ConnectionInfo struct {
	Host     string
	Port     int
	Username string
	DBName   string
	SSLMode  string
	Password string
}

type PostgresStorage struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

// NewPostgresConnection создаёт новое подключение к базе данных PostgreSQL.
// Функция принимает параметры подключения через структуру ConnectionInfo и возвращает объект *sql.DB,
// который можно использовать для взаимодействия с базой данных.
func NewPostgresConnection(info ConnectionInfo) (*sql.DB, error) {
	// Формируем строку подключения для PostgreSQL с использованием параметров из ConnectionInfo
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		info.Host, info.Port, info.Username, info.DBName, info.SSLMode, info.Password))
	if err != nil {
		// Возвращаем ошибку, если не удалось создать подключение
		return nil, err
	}

	// Пингуем базу данных для проверки доступности
	if err := db.Ping(); err != nil {
		// Возвращаем ошибку, если не удаётся установить соединение с базой данных
		return nil, err
	}

	// Возвращаем объект базы данных, если подключение прошло успешно
	return db, nil
}

// Закрытие соединения с базой данных
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// RunMigrations выполняет миграции базы данных, используя библиотеку migrate.
// Функция инициализирует драйвер для миграций, загружает и применяет миграции из указанной папки.
func RunMigrations() {
	// Создаём миграционный драйвер для подключения к базе данных PostgreSQL.
	// С помощью WithInstance и передаём текущую базу данных и её конфигурацию.
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		// Если не удалось создать миграционный драйвер, логируем ошибку и завершаем приложение.
		log.Fatalf("Error creating migration driver: %v", err)
	}

	// Создаём мигратор, который будет управлять миграциями.
	// Указываем путь к папке с миграциями ("file://migration") и тип базы данных ("postgres").
	m, err := migrate.NewWithDatabaseInstance(
		"file://migration", // Путь к папке с миграциями
		"postgres",         // Имя базы данных
		driver,
	)
	if err != nil {
		// Если не удалось создать мигратор, логируем ошибку и завершаем приложение.
		log.Fatalf("Error initializing migration: %v", err)
	}

	// Применяем все миграции. Если ошибок в процессе нет, применяются все миграции.
	// Если изменений нет, игнорируем это.
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		// Если при применении миграции произошла ошибка (кроме ситуации, когда изменений нет),
		// логируем ошибку и завершаем приложение.
		log.Fatalf("Error applying migration: %v", err)
	}

	// Если миграции успешно применены, выводим сообщение в лог.
	log.Println("Migrations applied successfully")
}

// RollbackLastMigration откатывает последнюю применённую миграцию из базы данных.
// Функция использует библиотеку migrate для выполнения отката миграции.
func RollbackLastMigration() {
	// Создаём миграционный драйвер для подключения к базе данных PostgreSQL.
	// С помощью WithInstance и передаём текущую базу данных и её конфигурацию.
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		// Если не удалось создать миграционный драйвер, логируем ошибку и завершаем приложение.
		log.Fatalf("Error creating migration driver: %v", err)
	}

	// Создаём мигратор, который будет управлять миграциями.
	// Указываем путь к папке с миграциями ("file://migration") и тип базы данных ("postgres").
	m, err := migrate.NewWithDatabaseInstance(
		"file://migration", // Путь к папке с миграциями
		"postgres",         // Имя базы данных
		driver,
	)
	if err != nil {
		// Если не удалось создать мигратор, логируем ошибку и завершаем приложение.
		log.Fatalf("Error initializing migration: %v", err)
	}

	// Откатываем последнюю миграцию с помощью m.Steps(-1), где -1 означает откат на 1 шаг.
	// Если не удаётся откатить миграцию, логируем ошибку и завершаем приложение.
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error rolling back migration: %v", err)
	}

	// Если откат прошел успешно, выводим сообщение в лог.
	log.Println("Last migration rolled back successfully")
}

// SetDB используется для установки подключения к базе данных.
// Эта функция полезна для тестов, когда мы можем использовать mock-объект базы данных.
func SetDB(mockDB *sql.DB) {
	// Присваиваем переданное подключение к mockDB для использования в других частях кода.
	db = mockDB
}
