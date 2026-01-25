package main

import (
	"fmt"
	"net/http"
	"shortener/internal/handlers"
	"shortener/internal/storage"
)

func main() {
	store := storage.New()

	h := &handlers.Handler{Storage: store}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", h.CreateShortUrl)
	mux.HandleFunc("GET /{id}", h.Redirect)

	fmt.Println("Сервер запущен")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Ошибка запуска:", err)
	}
}
