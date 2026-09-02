package handler

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/ElJhonnypro/url-shortener-go/internal/logger"
	"github.com/ElJhonnypro/url-shortener-go/internal/service"
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
func (h *URLHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
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
