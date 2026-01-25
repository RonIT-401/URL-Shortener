package handlers

import (
	"encoding/json"
	"fmt"
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
	Storage *storage.MemStorage
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

	h.Storage.Save(id, request.URL) // Сохраняем короткий ID и длинный ID

	// Записываем ответ в переменную
	response := ResponseJSON{
		Result: fmt.Sprintf("http://lacalhost:8080/%s", id),
	}

	// Передаем тип контента что бы Postman корректно обработал данные
	w.Header().Set("Content-Type", "application/json")

	// Устанавливаем HTTP статус 201 (Created)
	w.WriteHeader(http.StatusCreated)

	// Подключаемся на выход к пользователю и отправляем текст JSON в сеть
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	url, ok := h.Storage.Get(id)
	if !ok {
		http.Error(w, "Не найдено", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
