package db

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Project struct {
	ID          int64
	Name        string `sql:"not null;unique"`
	Description string `sql:"size:NULL"`
	Site        string
	Active      bool
}

type Post struct {
	ID        int64
	Short     string `sql:"not null;unique"`
	Title     string `sql:"size:NULL"`
	Content   string `sql:"size:NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Tags      []Tag `gorm:"many2many:posts_tags;"`
}

type Tag struct {
	ID    int64
	Name  string `sql:"not null;unique"`
	Posts []Post `gorm:"many2many:posts_tags;"`
}

type Asset struct {
	ID      int64
	Name    string `sql:"not null"`
	Type    string `sql:"not null"`
	Kind    string `sql:"not null"`
	Content []byte
}

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("postgres", "host=database user=root port=26257 dbname=blog sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	DB.AutoMigrate(&Project{}, &Post{}, &Tag{}, &Asset{})
	DB.Model(&Asset{}).AddUniqueIndex("idx_asset_name_type_kind", "name", "type", "kind")
}
