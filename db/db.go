package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tidwall/buntdb"
)

var BDB *buntdb.DB
var ErrNotFound = buntdb.ErrNotFound

func bail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Get(format string, args ...interface{}) (string, error) {
	var result string
	err := BDB.View(func(tx *buntdb.Tx) error {
		var err error
		result, err = tx.Get(fmt.Sprintf(format, args...))
		return err
	})
	return result, err
}

type cloner interface {
	clone() cloner
}

type augmenter interface {
	augment(string)
}

type Project struct {
	Name        string
	Description string
	Site        string
	Active      bool
}

func (p *Project) clone() cloner {
	np := *p
	return &np
}

func GetProjects() ([]interface{}, error) {
	return getJSONList((*buntdb.Tx).DescendKeys, "project:*", &Project{}, 0)
}

type Post struct {
	Short     string `json:"omitempty"`
	Title     string
	Content   string
	CreatedAt time.Time
}

func (p *Post) clone() cloner {
	np := *p
	return &np
}

func (p *Post) augment(key string) {
	p.Short = strings.Split(key, "/")[1]
}

func GetPosts(limit int) ([]interface{}, error) {
	return getJSONList((*buntdb.Tx).Ascend, "post:CreatedAt", &Post{}, limit)
}

type bIterFunc func(*buntdb.Tx, string, func(string, string) bool) error

func getJSONList(f bIterFunc, key string, out cloner, limit int) ([]interface{}, error) {
	var outs []interface{}
	var err error
	i := 0
	berr := BDB.View(func(tx *buntdb.Tx) error {
		return f(tx, key, func(key, value string) bool {
			if limit > 0 && i >= limit {
				return false
			}
			if err = json.Unmarshal([]byte(value), out); err != nil {
				return false
			}
			if a, ok := out.(augmenter); ok {
				a.augment(key)
			}
			outs = append(outs, out)
			out = out.clone()
			i++
			return true
		})
	})
	if err != nil {
		return nil, err
	}
	return outs, berr
}

func init() {
	if len(os.Args) < 2 {
		bail(fmt.Errorf("needs DB path argument"))
	}
	var err error
	BDB, err = buntdb.Open(os.Args[1])
	bail(BDB.SetConfig(buntdb.Config{SyncPolicy: buntdb.Always}))
	bail(BDB.CreateIndex("post:CreatedAt", "post:*", buntdb.Desc(buntdb.IndexJSON("CreatedAt"))))
	bail(err)
}
