package httprouter

import (
	"net/http"

	"github.com/ElJhonnypro/url-shortener-go/internal/handler"
	"github.com/ElJhonnypro/url-shortener-go/internal/middleware"
)

func NewRouter(
	handler *handler.URLHandler,
	createLimiter *middleware.RateLimiter,
	redirectLimiter *middleware.RateLimiter,
) http.Handler {

	mux := http.NewServeMux()

	mux.Handle(
		"POST /newurl",
		createLimiter.Middleware(http.HandlerFunc(handler.Create)),
	)

	mux.Handle(
		"GET /code/{code}",
		redirectLimiter.Middleware(http.HandlerFunc(handler.Get)),
	)

	mux.Handle(
		"GET /",
		redirectLimiter.Middleware(http.HandlerFunc(handler.Status)),
	)

	mux.Handle(
		"GET /status",
		redirectLimiter.Middleware(http.HandlerFunc(handler.Status)),
	)

	return mux
}
