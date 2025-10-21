package handler

import (
	"encoding/json"
	"net/http"

	"url-shortener/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	router  *mux.Router
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	h := &Handler{
		router:  mux.NewRouter(),
		service: service,
	}
	h.registerRouters()
	return h
}

func (h *Handler) registerRouters() {
	h.router.HandleFunc("/health", h.HealthHandler).Methods("GET")
	h.router.HandleFunc("/shorten", h.ShortenHandler).Methods("POST")
	h.router.HandleFunc("/{shortCode}", h.RedirectHandler).Methods("GET")
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req service.CreateShortURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid JSON format",
		})
		return
	}

	// 2. 调用Service层处理业务逻辑
	resp, err := h.service.Shorten.CreateShortURL(r.Context(), &req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// 3. 返回成功响应
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	// 1. 调用Service层获取长链接
	longURL, err := h.service.Shorten.GetLongURL(r.Context(), shortCode)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Short URL not found",
		})
		return
	}

	// 2. 执行重定向
	http.Redirect(w, r, longURL, http.StatusFound)
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
