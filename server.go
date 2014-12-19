package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func tmplExec(w http.ResponseWriter, subtmpl string, args interface{}) {
	c, err := tmpl.Clone()
	if err != nil {
		log.Println(err)
		http.Error(w, "500", 500)
		return
	}

	if subtmpl != "" {
		_, err = c.New("i").Parse(subtmpl)
		if err != nil {
			log.Println(err)
			http.Error(w, "500", 500)
			return
		}
	}
	err = c.Execute(w, args)
	if err != nil {
		log.Println(err)
	}
}

func HRoot(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	tmplExec(w, "", nil)
}

func HAsset(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var asset Asset
	if db.Find(&asset, "name = ? and kind = ?", ps.ByName("id"), ps.ByName("kind")).RecordNotFound() {
		http.Error(w, "404", 404)
		return
	}
	http.ServeContent(w, req, "", time.Time{}, bytes.NewReader(asset.Content))
}

func HPosts(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var posts []Post
	db.Select("title, short, created_at").Order("created_at desc").Find(&posts)
	tmplExec(w, r, posts)
}

func HPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var post Post
	db.First(&post, "short = ?", ps.ByName("id"))
	tmplExec(w, p, post)
}

func HProjects(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var projects []Project
	db.Order("id desc").Find(&projects)
	tmplExec(w, t, projects)
}

func main() {
	router := httprouter.New()
	router.GET("/", HRoot)
	router.GET("/assets/:kind/:id", HAsset)
	router.GET("/posts", HPosts)
	router.GET("/posts/:id", HPost)
	router.GET("/projects", HProjects)

	log.Fatal(http.ListenAndServe(":9000", router))
}
