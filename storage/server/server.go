package server

import (
	"encoding/json"
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

type translationObject struct {
	Lang  string `json:"lang"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func encode(t []*i18n.Translation) []byte {
	objs := make([]*translationObject, 0, len(t))
	for _, item := range t {
		objs = append(objs, &translationObject{
			Lang:  item.Lang.String(),
			Key:   item.Key,
			Value: item.Value,
		})
	}
	data, _ := json.Marshal(objs)
	return data
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	translations, err := server.storage.GetAll()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(encode(translations))
}
