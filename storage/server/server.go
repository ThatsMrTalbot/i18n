package server

import (
	"net/http"

	"github.com/ThatsMrTalbot/i18n"
)

type Server struct {
	storage i18n.Storage
}

func NewServer(storage i18n.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
