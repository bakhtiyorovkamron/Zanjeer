package storage

import (
	"github.com/Projects/Zanjeer/config"
	"github.com/Projects/Zanjeer/pkg/db"
	"github.com/Projects/Zanjeer/pkg/logger"
	"github.com/Projects/Zanjeer/storage/postgres"
)

type StorageI interface {
	Postgres() postgres.PostgresI
}
type StoragePg struct {
	postgres postgres.PostgresI
}

func (s *StoragePg) Postgres() postgres.PostgresI {
	return s.postgres
}
func New(db *db.Postgres, log *logger.Logger, cfg config.Config) StorageI {
	return &StoragePg{
		postgres: postgres.New(db, log, cfg),
	}
}
