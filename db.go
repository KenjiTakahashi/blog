package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Project struct {
	Id          int64
	Name        string `sql:"not null;unique"`
	Description string `sql:"size:NULL"`
	Site        string
	Active      bool
}

type Post struct {
	Id        int64
	Short     string `sql:"not null;unique"`
	Title     string `sql:"size:NULL"`
	Content   string `sql:"size:NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Tags      []Tag `gorm:"many2many:posts_tags;"`
}

type Tag struct {
	Id    int64
	Name  string `sql:"not null;unique"`
	Posts []Post `gorm:"many2many:posts_tags;"`
}

type Asset struct {
	Id      int64
	Name    string `sql:"not null"`
	Type    string `sql:"not null"`
	Kind    string `sql:"not null"`
	Content []byte
}

var db gorm.DB

func init() {
	var err error
	db, err = gorm.Open("postgres", "user=postgres password=postgres dbname=blog sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Project{}, &Post{}, &Tag{}, &Asset{})
	db.Model(&Asset{}).AddUniqueIndex("idx_asset_name_type_kind", "name", "type", "kind")
}
