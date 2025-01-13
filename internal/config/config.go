package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

// Config структура для хранения всей конфигурации приложения
// Она содержит информацию о базе данных, сервере, обменном сервисе и Redis.
type Config struct {
	// Структура, которая хранит настройки подключения к базе данных PostgreSQL
	DB Postgres
	// Структура для параметров сервера
	Server struct {
		// Порт, на котором будет работать сервер
		// Использует значение по умолчанию 8080, если переменная окружения SERVER_PORT не задана
		Port int `envconfig:"SERVER_PORT" default:"8080"`
		// Секретный ключ для JWT (обязателен)
		JWTSecret string `envconfig:"JWT_SECRET" required:"true"`
	}

	// Структура для конфигурации обменного сервиса
	// Адрес обменного сервиса (обязателен)
	ExchangeService struct {
		Address string `envconfig:"EXCHANGE_SERVICE_ADDRESS" required:"true"`
	}
}

// Структура для хранения параметров подключения к базе данных PostgreSQL
type Postgres struct {
	// Хост PostgreSQL сервера (обязателен)
	Host string `envconfig:"DB_HOST" required:"true"`
	// Порт PostgreSQL сервера (обязателен)
	Port int `envconfig:"DB_PORT" required:"true"`
	// Имя пользователя для подключения к базе данных (обязателен)
	Username string `envconfig:"DB_USERNAME" required:"true"`
	// Имя базы данных (обязателен)
	Name string `envconfig:"DB_NAME" required:"true"`
	// Режим SSL для подключения к базе данных (по умолчанию "disable")
	SSLMode string `envconfig:"DB_SSLMODE" default:"disable"`
	// Пароль для подключения к базе данных (обязателен)
	Password string `envconfig:"DB_PASSWORD" required:"true"`
}

// Функция New загружает конфигурацию из переменных окружения и возвращает структуру Config
func New() (*Config, error) {
	// Загружаем файл .env в случае его наличия
	// Это нужно для использования локальных значений в разработке
	if err := godotenv.Load("config.env"); err != nil {
		// Если файл .env не найден, выводим сообщение, но продолжаем работу с системными переменными окружения
		log.Println("No config.env file found, using system environment variables")
	}

	// Создаем новый экземпляр конфигурационной структуры
	cfg := new(Config)

	// Заполняем структуру cfg значениями из переменных окружения
	// Если переменная окружения не задана, будет использовано значение по умолчанию, если оно указано
	if err := envconfig.Process("", cfg); err != nil {
		// Возвращаем ошибку, если не удалось обработать переменные окружения
		return nil, err
	}

	// Возвращаем структуру конфигурации
	return cfg, nil
}
