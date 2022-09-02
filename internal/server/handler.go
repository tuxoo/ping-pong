package server

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Print("hello")
}
