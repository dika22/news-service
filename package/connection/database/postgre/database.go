package postgre

import (
	"database/sql"
	"fmt"
	"log"
	"news-service/package/config"

	_ "github.com/lib/pq"
)

func NewDatabase(c *config.Database) *sql.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta",
		c.DBHost,
		c.DBUsername,
		c.DBPassword,
		c.DBDBName,
		c.DBPort,
	)

	// Open a database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}

	// Check if the connection is alive
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	fmt.Println("Successfully connected to the database!")
	return db
}