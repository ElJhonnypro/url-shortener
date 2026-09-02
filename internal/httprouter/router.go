package httprouter

import (
	"net/http"

	"github.com/ElJhonnypro/url-shortener-go/internal/handler"
)

func NewRouter(handler *handler.URLHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /newurl", handler.Create)
	mux.HandleFunc("GET /code/{code}", handler.Get)

	return mux
}
