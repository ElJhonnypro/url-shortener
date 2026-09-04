package service

import (
	"github.com/ElJhonnypro/url-shortener/internal/logger"
	"github.com/ElJhonnypro/url-shortener/internal/repository"
)

type URLService struct {
	repository *repository.URLRepository
	log        *logger.Logger
}

func NewURLService(repository *repository.URLRepository, log *logger.Logger) *URLService {
	return &URLService{
		repository: repository,
		log:        log,
	}
}

func (s *URLService) Create(url string, code string) error {
	return s.repository.Create(url, code)
}

func (s *URLService) Get(code string) (string, error) {
	return s.repository.Get(code)
}
