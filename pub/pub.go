package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/KenjiTakahashi/blog/db"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/cli"
)

func bail(err error) {
	if err != nil {
		panic(err)
	}
}

func getbool(txt string) bool {
	return txt == "true" || txt == "yes" || txt == "1"
}

func dbe(db *gorm.DB) int {
	if db.Error != nil {
		log.Println("Error creating/updating/deleting record: ", db.Error)
		return 1
	}
	return 0
}

func dbc(obj interface{}) int {
	return dbe(db.DB.Create(obj))
}

func dbu(obj interface{}) int {
	return dbe(db.DB.Save(obj))
}

func dbr(obj interface{}) int {
	return dbe(db.DB.Delete(obj))
}

type pC struct{}

func (c *pC) Run(args []string) int {
	if len(args) != 1 {
		log.Println("Invalid number of arguments")
		return 1
	}

	post, err := ioutil.ReadFile(args[0])
	bail(err)
	lines := bytes.SplitN(post, []byte{'\n'}, 4)

	short := filepath.Base(args[0])

	time, err := time.Parse("2 Jan 2006", string(lines[0]))
	bail(err)

	tagsB := [][]byte{}
	if len(lines[2]) > 0 {
		tagsB = bytes.Split(lines[2], []byte{','})
	}
	tags := make([]db.Tag, len(tagsB))
	for i, tag := range tagsB {
		tags[i] = db.Tag{Name: string(tag)}
	}

	dbpost := &db.Post{
		Short:     short[:len(short)-3],
		Title:     string(lines[1]),
		Content:   string(lines[3]),
		Tags:      tags,
		CreatedAt: time,
	}
	if errno := dbc(dbpost); errno != 0 {
		return errno
	}

	return dbe(db.DB.Model(dbpost).Update("CreatedAt", time))
}
func (c *pC) Help() string {
	return c.Synopsis()
}
func (c *pC) Synopsis() string {
	return "p <file> - publish post"
}

type tC struct{}

func (c *tC) Run(args []string) int {
	if len(args) != 1 {
		log.Println("Invalid number of arguments")
		return 1
	}

	fi, err := os.Open(args[0])
	bail(err)
	defer fi.Close()
	cr := csv.NewReader(fi)
	cr.Comma = ' '
	projects, err := cr.ReadAll()
	bail(err)

	tx := db.DB.Begin()

	if errno := dbr(&db.Project{}); errno != 0 {
		tx.Rollback()
		return errno
	}

	for _, project := range projects {
		errno := dbc(&db.Project{
			Name:        project[0],
			Description: project[1],
			Site:        project[2],
			Active:      getbool(project[3]),
		})
		if errno != 0 {
			tx.Rollback()
			return errno
		}
	}

	return dbe(tx.Commit())
}
func (c *tC) Help() string {
	return c.Synopsis()
}
func (c *tC) Synopsis() string {
	return "t <file> - add project"
}

type aC struct{}

func (c *aC) Run(args []string) int {
	if len(args) != 2 {
		log.Println("Invalid number of arguments")
		return 1
	}

	asset, err := ioutil.ReadFile(args[1])
	bail(err)
	typ := filepath.Ext(args[1])
	name := strings.TrimSuffix(filepath.Base(args[1]), typ)

	return dbc(&db.Asset{
		Name:    name,
		Type:    typ[1:],
		Kind:    args[0],
		Content: asset,
	})
}
func (c *aC) Help() string {
	return c.Synopsis()
}
func (c *aC) Synopsis() string {
	return "a <kind> <file> - add asset code"
}

func main() {
	ui := cli.NewCLI("pub", "0.2")
	ui.Args = os.Args[1:]
	ui.Commands = map[string]cli.CommandFactory{
		"p": func() (cli.Command, error) {
			return &pC{}, nil
		},
		"t": func() (cli.Command, error) {
			return &tC{}, nil
		},
		"a": func() (cli.Command, error) {
			return &aC{}, nil
		},
	}

	code, err := ui.Run()
	bail(err)

	os.Exit(code)
}
