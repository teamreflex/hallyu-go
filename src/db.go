package bot

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	RemoteID uint `gorm:"index"`
  Title string
	Handle string
}

func GetDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})

	if err != nil {
		fmt.Printf("Failed to connect to database: %s\n", err)
		os.Exit(1)
	}

	return db
}