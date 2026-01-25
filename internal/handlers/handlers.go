package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shortener/internal/storage"
	"shortener/internal/utils"
)

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}

type Handler struct {
	Storage *storage.MemStorage
}

func (h *Handler) CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	var request RequestJSON

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Некоректный URL", http.StatusBadRequest)
		return
	}

	id := utils.GenerateID(6)

	response := ResponseJSON{
		Result: fmt.Sprintf("http://lacalhost:8080/%s", id),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

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
