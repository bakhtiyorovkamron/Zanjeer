package postgres

import (
	"github.com/Projects/Zanjeer/config"
	"github.com/Projects/Zanjeer/pkg/db"
	"github.com/Projects/Zanjeer/pkg/logger"
)

type postgresRepo struct {
	Db  *db.Postgres
	Log *logger.Logger
	Cfg config.Config
}

func New(db *db.Postgres, log *logger.Logger, cfg config.Config) PostgresI {
	return &postgresRepo{
		Db:  db,
		Log: log,
		Cfg: cfg,
	}
}
