package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func H404(w http.ResponseWriter, req *http.Request) {
	log.Printf("code :: 404 :: %s :: %s", req.URL, req.RemoteAddr)
	w.WriteHeader(http.StatusNotFound)
	etmpl.Execute(w, 404)
}

func H500(w http.ResponseWriter, req *http.Request, rcv interface{}) {
	log.Printf("code :: 500 :: %s :: %s", req.URL, req.RemoteAddr)
	if rcv != nil {
		log.Println(rcv)
	}
	w.WriteHeader(http.StatusInternalServerError)
	etmpl.Execute(w, 500)
}

func tmplExec(w http.ResponseWriter, req *http.Request, subtmpl string, args interface{}) {
	c, err := tmpl.Clone()
	if err != nil {
		H500(w, req, nil)
		return
	}

	if subtmpl != "" {
		_, err = c.New("i").Parse(subtmpl)
		if err != nil {
			H500(w, req, nil)
			return
		}
	}
	err = c.Execute(w, args)
	if err != nil {
		log.Println(err)
	}
}

func HRoot(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	tmplExec(w, req, "", nil)
}

func HAsset(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var asset Asset
	if db.Find(&asset, "name = ? and kind = ?", ps.ByName("id"), ps.ByName("kind")).RecordNotFound() {
		H404(w, req)
		return
	}
	http.ServeContent(w, req, "", time.Time{}, bytes.NewReader(asset.Content))
}

func HPosts(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var posts []Post
	db.Select("title, short, created_at").Order("created_at desc").Find(&posts)
	tmplExec(w, req, r, posts)
}

func HPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var post Post
	if db.First(&post, "short = ?", ps.ByName("id")).RecordNotFound() {
		H404(w, req)
		return
	}
	tmplExec(w, req, p, post)
}

func HProjects(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var projects []Project
	db.Order("id desc").Find(&projects)
	tmplExec(w, req, t, projects)
}

func main() {
	router := &httprouter.Router{
		NotFound: H404,
		PanicHandler: H500,
	}
	router.GET("/", HRoot)
	router.GET("/assets/:kind/:id", HAsset)
	router.GET("/posts", HPosts)
	router.GET("/posts/:id", HPost)
	router.GET("/projects", HProjects)

	log.Fatal(http.ListenAndServe(":9000", router))
}
