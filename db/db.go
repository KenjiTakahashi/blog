package db

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/karrick/godirwalk"
	"golang.org/x/sys/unix"
)

var db = ""

var ErrNotFound = fmt.Errorf("NF")

type Project struct {
	Name        string
	Description string
	Site        string
	Active      bool
}

func GetProjects() ([]Project, error) {
	defcLock.RLock()
	defer defcLock.RUnlock()

	return defc.Projects, nil
}

func GetAsset(kind, id string) ([]byte, error) {
	defcLock.RLock()
	defer defcLock.RUnlock()

	assets, exist := defc.Assets[kind]
	if !exist {
		return nil, ErrNotFound
	}
	asset, exists := assets[id]
	if !exists {
		return nil, ErrNotFound
	}
	return asset, nil
}

type Post struct {
	Short     string
	CreatedAt time.Time
	Title     string
	Kind      string
	Content   []byte
}

func GetPost(kind, id string) (*Post, error) {
	defcLock.RLock()
	defer defcLock.RUnlock()

	return defc.Posts.GetOne(id)
}

func GetPosts(limit int) ([]*Post, error) {
	defcLock.RLock()
	defer defcLock.RUnlock()

	return defc.Posts.GetMany(limit), nil
}

func readProjects(filepath string) ([]Project, error) {
	fi, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	cr := csv.NewReader(fi)
	cr.Comma = ' '
	projects, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}

	projectsOut := make([]Project, len(projects))
	for i, project := range projects {
		projectsOut[len(projects)-i-1] = Project{
			Name:        project[0],
			Description: project[1],
			Site:        project[2],
			Active:      project[3] == "yes",
		}
	}
	return projectsOut, nil
}

func readPost(filepath string) (*Post, error) {
	post, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	lines := bytes.SplitN(post, []byte{'\n'}, 4)

	time, err := time.Parse("2 Jan 2006", string(lines[0]))
	if err != nil {
		return nil, err
	}

	var tagsB [][]byte
	if len(lines[2]) > 0 {
		tagsB = bytes.Split(lines[2], []byte{','})
	}
	kind := "science"
	if len(tagsB) == 1 {
		kind = string(tagsB[0])
	}
	return &Post{
		CreatedAt: time,
		Title:     string(lines[1]),
		Kind:      kind,
		Content:   lines[3],
	}, nil
}

var assetKinds = map[string]string{
	".png":  "image",
	".html": "raw",
	".js":   "script",
}

var inotfd int

func fillCache() error {
	cache := newCache()
	err := godirwalk.Walk(db, &godirwalk.Options{
		Callback: func(relpath string, de *godirwalk.Dirent) error {
			name := de.Name()
			if de.IsDir() {
				_, err := unix.InotifyAddWatch(
					inotfd, relpath,
					unix.IN_MODIFY|unix.IN_CREATE|unix.IN_DELETE|unix.IN_MOVED_TO,
				)
				if err != nil {
					log.Printf("E006 : %s", err)
				}
				return nil
			}
			if name == "projects" {
				projects, err := readProjects(relpath)
				if err != nil {
					log.Printf("E001 : %s", err)
					return nil
				}
				cache.Projects = projects
				return nil
			}
			ext := path.Ext(name)
			base := strings.TrimSuffix(name, ext)
			if ext == ".md" {
				post, err := readPost(relpath)
				if err != nil {
					log.Printf("E002 : %s", err)
					return nil
				}
				cache.Posts.Set(base, post)
				return nil
			}
			kind, exists := assetKinds[ext]
			if !exists {
				return nil
			}
			asset, err := ioutil.ReadFile(relpath)
			if err != nil {
				log.Printf("E003 : %s", err)
				return nil
			}
			cache.Assets[kind][base] = asset
			return nil
		},
		ErrorCallback: func(_ string, _ error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true,
	})
	if err != nil {
		return err
	}
	sort.Slice(cache.Posts.List, func(i, j int) bool {
		return cache.Posts.List[i].CreatedAt.After(cache.Posts.List[j].CreatedAt)
	})

	defcLock.Lock()
	defc = cache
	defcLock.Unlock()
	return nil
}

func init() {
	db = os.Args[1]

	go func() {
		for {
			fd, err := unix.InotifyInit()
			if err != nil {
				log.Printf("E005 : %s", err)
				continue
			}
			inotfd = fd

			if err = fillCache(); err != nil {
				log.Fatalf("E004 : %s", err)
			}

			var buffer [64 * (unix.SizeofInotifyEvent + unix.PathMax + 1)]byte
			if _, err = unix.Read(inotfd, buffer[:]); err != nil {
				log.Fatalf("E007 : %s", err)
			}
			if err = unix.Close(inotfd); err != nil {
				log.Fatalf("E008 : %s", err)
			}
		}
	}()
}
