package migrate

import (
	"github.com/dika22/news-service/package/structs"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

type Migrate struct {
	db *gorm.DB
}

func (h *Migrate) Migrate(c *cli.Context) error {
	return h.db.AutoMigrate(
		&structs.Articles{},
	)
}

func NewMigrate(db *gorm.DB) []*cli.Command {
	h := Migrate{
		db: db,
	}

	return []*cli.Command{
		{
			Name:   "migrate-db",
			Usage:  "Migrate database",
			Action: h.Migrate,
		},
	}
}
