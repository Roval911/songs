package hanlers

import (
	"github.com/sirupsen/logrus"
	"songs/internal/config"
	"songs/internal/storages"
)

type Handler struct {
	storage storages.Storages
	logger  *logrus.Logger
	config  *config.Config
}

func NewHandler(storage storages.Storages, logger *logrus.Logger, cfg *config.Config) *Handler {
	return &Handler{
		storage: storage,
		logger:  logger,
		config:  cfg,
	}
}
