package database

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cast"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/dika22/news-service/package/config"
)

type Database struct {
	c config.Database
}

const (
	WebDB       = "wdb"
	SqliteDB    = "sqlite"
	SqliteDBWeb = "sqlite-web"
)

func NewDatabase(connType string, c *config.Database) *gorm.DB {
	db := &gorm.DB{}
	gormCfg := &gorm.Config{
		SkipDefaultTransaction: true,
		// Logger:                 logger.Default.LogMode(logger.Info),
	}
	var dialector gorm.Dialector
	var err error
	switch connType {
	case SqliteDB, SqliteDBWeb:
		dir, _ := os.Getwd()

		match, _ := regexp.Match(`/internal/domain/.+/(delivery)`, []byte(dir))
		if match {
			dir += "/../../../.."
		}

		match, _ = regexp.Match(`/pkg/validator`, []byte(dir))
		if match {
			dir += "/../../"
		}
		match, _ = regexp.Match(`/cmd/api/middleware`, []byte(dir))
		if match {
			dir += "../../../"
		}

		match, _ = regexp.Match(`/internal/domain/.+/(repository)/.+`, []byte(dir))
		if match {
			dir += "/../../../../.."
		}

		match, _ = regexp.Match(`/internal/domain/.+/(usecase)`, []byte(dir))
		if match {
			dir += "../../../../.."
		}

		dir = strings.ReplaceAll(dir, "//", "/")
		file := dir + "/sqlite.db"
		if connType == SqliteDBWeb {
			file = dir + "/sqlite-web.db"
		}
		dialector = sqlite.Open(file)
	case WebDB:
		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta",
		    c.DBHost,	
		    c.DBUsername,
			c.DBPassword,
			c.DBDBName,
			c.DBPort,
		)
		dialector = postgres.New(postgres.Config{
			DSN: dsn,
		})
	default:
		panic("unknown connection type")
	}

	db, err = gorm.Open(dialector, gormCfg)
	if err != nil {
		panic(err)
	}
	dbConn, err := db.DB()
	if err != nil {
		panic(err)
	}
	dbConn.SetMaxIdleConns(cast.ToInt(c.IdleConns))
	dbConn.SetMaxOpenConns(cast.ToInt(c.MaxConns))

	return db
}
