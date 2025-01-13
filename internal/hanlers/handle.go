package hanlers

import (
	"github.com/sirupsen/logrus"
	"songs/internal/storages"
)

type Handler struct {
	storage storages.Storages
	logger  *logrus.Logger
}

func NewHandler(storage storages.Storages, logger *logrus.Logger) *Handler {
	return &Handler{
		storage: storage,
		logger:  logger,
	}
}
