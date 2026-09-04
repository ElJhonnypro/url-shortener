package handler

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"

	"github.com/ElJhonnypro/url-shortener/internal/logger"
	"github.com/ElJhonnypro/url-shortener/internal/service"
)

type URLHandler struct {
	service *service.URLService
	log     *logger.Logger
}

func NewURLHandler(service *service.URLService, log *logger.Logger) *URLHandler {
	return &URLHandler{
		service: service,
		log:     log,
	}
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func generateCode(length int) (string, error) {
	code := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err // Properly bubble up the error
		}
		code[i] = charset[num.Int64()]
	}

	return string(code), nil
}
func isValidURL(rawURL string) bool {
	if len(rawURL) < 5 || len(rawURL) > 2048 {
		return false
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	if u.Scheme == "" || u.Host == "" {
		return false
	}

	switch u.Scheme {
	case "http", "https":
		return true
	}
	return false
}

func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusUnsupportedMediaType)
		return
	}

	// Validate the URL
	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusNotAcceptable)
		return
	}

	// Validate the URL format (basic validation)
	if !isValidURL(req.URL) {
		http.Error(w, "Invalid URL format", http.StatusNotAcceptable)
		return
	}

	Code, err := generateCode(5)
	if err != nil {
		h.log.Error("Failed to generate code: " + err.Error())
		http.Error(w, "Failed to generate code", http.StatusInternalServerError)
		return
	}
	err = h.service.Create(req.URL, Code)
	if err != nil {
		h.log.Error("Failed to create URL: " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(
		map[string]any{"code": Code},
	)

}

func (h *URLHandler) Get(w http.ResponseWriter, r *http.Request) {
	Code := r.URL.Path[len("/code/"):]
	url, err := h.service.Get(Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)

}

func (h *URLHandler) Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
