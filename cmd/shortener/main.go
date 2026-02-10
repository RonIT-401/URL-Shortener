package main

import (
	"fmt"
	"net/http"
	"shortener/internal/handlers"
	"shortener/internal/storage/postgresql"
)

func main() {
	// Работа через память
	// store := memstorage.New()

	// Postgres
	dsn := "postgres://postgres:1111@localhost:5432/postgres"
	store, err := postgresql.New(dsn)
		if err != nil {
			panic(err)
		}

	h := &handlers.Handler{Storage: store}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", h.CreateShortUrl)
	mux.HandleFunc("GET /{id}", h.Redirect)

	fmt.Println("Сервер запущен")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Ошибка запуска:", err)
	}
}
