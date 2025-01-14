package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"songs/internal/config"
	"songs/internal/hanlers"
	"songs/internal/routes"
	"songs/internal/storages/postgres"
	"songs/pkg/logger"
)

type App struct {
	logger *logrus.Logger // Логгер для логирования событий
	router *gin.Engine    // Gin маршрутизатор для обработки HTTP-запросов // Клиент для работы с Redis
}

// New создаёт новый экземпляр приложения, инициализирует все необходимые сервисы
func New() (*App, error) {
	// Инициализация логгера с помощью вспомогательной функции
	log := logger.InitLogger()

	// Загрузка конфигурации из переменных окружения или файла
	cfg, err := config.New()
	if err != nil {
		// Логирование ошибки и завершение работы, если конфигурация не загружена
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
		return nil, err
	}

	// Инициализация подключения к базе данных PostgreSQL с параметрами из конфигурации
	db, err := postgres.NewPostgresConnection(postgres.ConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		// Логирование ошибки подключения и завершение работы, если соединение не удалось
		log.Fatal(err)
		return nil, err
	}

	// Устанавливаем подключение к базе данных и запускаем миграции
	postgres.SetDB(db)
	postgres.RunMigrations()

	// Создание хранилища данных для работы с PostgreSQL
	storage := postgres.NewPostgresStorage(db)

	err = storage.CreateIndexes()
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании индексов: %v", err)
	}

	// Создание обработчиков для аутентификации и обмена валютами
	Handler := hanlers.NewHandler(storage, log, cfg)

	// Настройка маршрутов для HTTP-сервера с использованием Gin
	router := routes.SetupRouter(Handler)

	// Возвращаем структуру приложения с логгером и маршрутизатором
	return &App{
		logger: log,
		router: router,
	}, nil
}

// Run запускает сервер приложения
// Он запускает HTTP-сервер на порту, указанном в конфигурации
func (a *App) Run() error {
	// Загрузка конфигурации для получения порта сервера
	cfg, err := config.New()
	if err != nil {
		// Логирование ошибки, если не удалось загрузить конфигурацию
		a.logger.Fatalf("Ошибка загрузки конфигурации: %v", err)
		return err
	}

	// Порт, на котором будет работать сервер
	port := cfg.Server.Port
	a.logger.Printf("Запуск сервера на порту: %d", port)

	// Запускаем сервер с использованием маршрутизатора и указываем порт для прослушивания
	err = a.router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		// Логирование ошибки, если сервер не смог запуститься
		a.logger.Fatalf("Ошибка запуска сервера: %v", err)
		return err
	}

	// Если сервер успешно запустился, возвращаем nil
	return nil
}
