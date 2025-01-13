package main

import (
	"log"
	_ "songs/docs"
	"songs/internal/app"
)

// @title Songs API
// @version 1.0
// @description API for searsh songs.
// @contact.name API Support
// @contact.email support@currencywallet.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("Ошибка инициализации приложения: %v", err)
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}
}
