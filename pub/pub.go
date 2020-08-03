package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/tidwall/buntdb"

	"github.com/KenjiTakahashi/blog/db"
)

type E struct {
	N int
	M error
}

func (e E) Error() string {
	return fmt.Sprintf("ERR %03d: %s", e.N, e.M)
}

func e(err error) int {
	e, ok := err.(E)
	if ok {
		if e.M != nil {
			log.Print(err)
			return e.N
		}
	} else if err != nil {
		log.Print(err)
		return 1000
	}
	return 0
}

func bail(err E) {
	if err.M != nil {
		os.Exit(e(err))
	}
}

func getbool(txt string) bool {
	return txt == "true" || txt == "yes" || txt == "1"
}

type pC struct{}

func (c *pC) Run(args []string) int {
	if len(args) != 1 {
		log.Println("Invalid number of arguments")
		return 9999
	}

	post, err := ioutil.ReadFile(args[0])
	bail(E{7, err})
	lines := bytes.SplitN(post, []byte{'\n'}, 4)

	short := filepath.Base(args[0])
	short = short[:len(short)-3]

	time, err := time.Parse("2 Jan 2006", string(lines[0]))
	bail(E{8, err})

	tagsB := [][]byte{}
	if len(lines[2]) > 0 {
		tagsB = bytes.Split(lines[2], []byte{','})
	}

	typ := "science"
	if len(tagsB) == 1 {
		typ = string(tagsB[0])
	}

	err = db.BDB.Update(func(tx *buntdb.Tx) error {
		p := db.Post{
			Title:     string(lines[1]),
			Content:   string(lines[3]),
			CreatedAt: time,
		}
		value, err := json.Marshal(&p)
		if err != nil {
			return E{3, err}
		}
		_, _, err = tx.Set(fmt.Sprintf("post:%s/%s", typ, short), string(value), nil)
		if err != nil {
			return E{4, err}
		}
		return nil
	})
	return e(err)
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
		return 9999
	}

	fi, err := os.Open(args[0])
	bail(E{9, err})
	defer fi.Close()
	cr := csv.NewReader(fi)
	cr.Comma = ' '
	projects, err := cr.ReadAll()
	bail(E{10, err})

	err = db.BDB.Update(func(tx *buntdb.Tx) error {
		for i, project := range projects {
			p := db.Project{
				Name:        project[0],
				Description: project[1],
				Site:        project[2],
				Active:      getbool(project[3]),
			}
			value, err := json.Marshal(&p)
			e(E{1, err})
			_, _, err = tx.Set(fmt.Sprintf("project:%03d", i), string(value), nil)
			e(E{2, err})
		}
		return nil
	})
	return e(E{5, err})
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
		return 9999
	}

	asset, err := ioutil.ReadFile(args[1])
	bail(E{11, err})
	ext := filepath.Ext(args[1])
	name := strings.TrimSuffix(filepath.Base(args[1]), ext)

	err = db.BDB.Update(func(tx *buntdb.Tx) error {
		_, _, err = tx.Set(fmt.Sprintf("asset:%s:%s", name, args[0]), string(asset), nil)
		if err != nil {
			return E{6, err}
		}
		return nil
	})
	return e(err)
}
func (c *aC) Help() string {
	return c.Synopsis()
}
func (c *aC) Synopsis() string {
	return "a <kind> <file> - add asset code"
}

func main() {
	ui := cli.NewCLI("pub", "0.3")
	ui.Args = os.Args[2:]
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
	bail(E{12, err})

	bail(E{13, db.BDB.Shrink()})

	os.Exit(code)
}
