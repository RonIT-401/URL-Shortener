package handlers

import (
	"fmt"
	"net/http"
	"shortener/internal/storage"
	"shortener/internal/utils"
)

type Handler struct {
	Storage *storage.MemStorage
}

func (h *Handler) CreateShortUrl(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Нужен url", http.StatusBadRequest)
		return
	}

	id := utils.GenerateID(6)
	
	h.Storage.Save(id, url)
	fmt.Fprintf(w, "Ваш ID: %s", id)
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