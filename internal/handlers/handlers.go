package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"shortener/internal/storage"
	"shortener/internal/utils"
)

// Структура для входящих данных
type RequestJSON struct {
	URL string `json:"url"`
}

// Структура для ответа
type ResponseJSON struct {
	Result string `json:"result"`
}

// Структура обработчик
type Handler struct {
	Storage storage.Storage
}

// Функция преобразования короткой ссылки
func (h *Handler) CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	var request RequestJSON // Переменная для данных из запроса

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Некоректный URL", http.StatusBadRequest)
		return
	}

	// Проверяем прислали ли нам хоть какой то URL
	if request.URL == "" {
		http.Error(w, "Поле url пустое", http.StatusBadRequest)
	}

	id := utils.GenerateID(6) // Генерируем уникальный ID

	exists, err := h.Storage.CheckExistURL(request.URL)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		http.Error(w, "Ссылка уже была сокращена", http.StatusAlreadyReported)

		url, ok, err := h.Storage.GetShort(request.URL)

		if err != nil {
			http.Error(w, "Ошибка сервера при поиске", http.StatusBadRequest)
			return
		}

		if !ok {
			http.Error(w, "Не найдено", http.StatusBadRequest)
			return
		}

		// Записываем ответ в переменную
		response := ResponseJSON{
			Result: fmt.Sprintf("http://localhost:8080/%s", url),
		}

		// Передаем тип контента что бы Postman корректно обработал данные
		w.Header().Set("Content-Type", "application/json")

		// Устанавливаем HTTP статус 201 (Created)
		w.WriteHeader(http.StatusCreated)

		// Подключаемся на выход к пользователю и отправляем текст JSON в сеть
		json.NewEncoder(w).Encode(response)

		return
	}

	err = h.Storage.Save(id, request.URL) // Сохраняем короткий ID и длинный ID
	if err != nil {
		http.Error(w, "Ошибка сохранения в базу", http.StatusBadRequest)
		return
	}

	// Записываем ответ в переменную
	response := ResponseJSON{
		Result: fmt.Sprintf("http://localhost:8080/%s", id),
	}

	// Передаем тип контента что бы Postman корректно обработал данные
	w.Header().Set("Content-Type", "application/json")

	// Устанавливаем HTTP статус 201 (Created)
	w.WriteHeader(http.StatusCreated)

	// Подключаемся на выход к пользователю и отправляем текст JSON в сеть
	json.NewEncoder(w).Encode(response)
}

// Переход на сайт с новой(короткой) ссылкой
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")           // Извлекаем ссылку
	url, ok, err := h.Storage.Get(id) // Ищем ссылку в хранилище

	if err != nil {
		http.Error(w, "Ошибка сервера при поиске", http.StatusBadRequest)
		return
	}

	if !ok {
		http.Error(w, "Не найдено", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusFound) // Перенаправляем пользователя по ссылке
}
