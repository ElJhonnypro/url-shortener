package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ElJhonnypro/url-shortener-go/internal/handler"
	"github.com/ElJhonnypro/url-shortener-go/internal/httprouter"
	"github.com/ElJhonnypro/url-shortener-go/internal/logger"
	"github.com/ElJhonnypro/url-shortener-go/internal/middleware"
	"github.com/ElJhonnypro/url-shortener-go/internal/repository"
	"github.com/ElJhonnypro/url-shortener-go/internal/service"
)

func main() {

	console := logger.NewConsoleOutput()

	file, err := logger.NewFileOutput("logs/app.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	log := logger.NewLogger(
		console,
		file,
	)

	repository := repository.NewURLRepository(log)

	if _, err := repository.Connect(); err != nil {
		log.Error("Database initialization failed: " + err.Error())
		return
	}

	defer repository.Close()
	service := service.NewURLService(repository, log)
	handler := handler.NewURLHandler(service, log)

	createLimiter := middleware.NewRateLimiter(
		10,
		time.Minute,
	)

	redirectLimiter := middleware.NewRateLimiter(
		60,
		time.Minute,
	)

	router := httprouter.NewRouter(
		handler,
		createLimiter,
		redirectLimiter,
	)

	cors := middleware.CORS(router)

	server := http.Server{
		Addr:    ":8080",
		Handler: cors,
	}

	server.RegisterOnShutdown(func() {
		log.Info("Server turning off.")
		log.Warn("BYE.")
	})

	log.Info("Server started at 8080")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer stop()

	go func() {
		<-ctx.Done()

		log.Info("Shutting down...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error("Server shutdown failed: " + err.Error())
		}
	}()

	err = server.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		log.Error("Server failed: " + err.Error())
	}
}
