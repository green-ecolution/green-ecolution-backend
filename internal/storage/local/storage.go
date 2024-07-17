package local

import (
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/local/info"
)

func NewRepository(cfg *config.Config) (*storage.Repository, error) {
	infoRepo, err := info.NewInfoRepository(cfg)
	if err != nil {
		return nil, err
	}

	return &storage.Repository{
		Info: infoRepo,
	}, nil
}
