package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateInvalidJSON(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodPost,
		"/newurl",
		strings.NewReader(`{"url":`),
	)

	rec := httptest.NewRecorder()

	handler := &URLHandler{}

	handler.Create(rec, req)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusUnsupportedMediaType,
			rec.Code,
		)
	}
}
func TestIsValidURL(t *testing.T) {

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "valid https",
			url:  "https://google.com",
			want: true,
		},
		{
			name: "valid http",
			url:  "http://example.com",
			want: true,
		},
		{
			name: "missing scheme",
			url:  "google.com",
			want: false,
		},
		{
			name: "invalid scheme",
			url:  "ftp://example.com",
			want: false,
		},
		{
			name: "random text",
			url:  "hello",
			want: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			got := isValidURL(tt.url)

			if got != tt.want {
				t.Errorf(
					"expected %v, got %v",
					tt.want,
					got,
				)
			}
		})
	}
}
