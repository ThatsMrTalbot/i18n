package server

import (
	"net/http"

	"github.com/ThatsMrTalbot/i18n"
)

// Middleware for server, return true to continue
type ResponseMiddleware func(w http.ResponseWriter, r *http.Request) bool

type Server struct {
	storage    i18n.Storage
	middleware []ResponseMiddleware
}

func NewServer(storage i18n.Storage, middleware ...ResponseMiddleware) *Server {
	return &Server{
		storage:    storage,
		middleware: middleware,
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, middleware := range server.middleware {
		if !middleware(w, r) {
			return
		}
	}

	translations, err := server.storage.GetAll()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	supported, err := server.storage.SupportedLanguages()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	def, err := server.storage.DefaultLanguage()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(encode(translations, supported, def))
}
