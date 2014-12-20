package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func H404(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	etmpl.Execute(w, 404)
	log.Printf("404 :: %s :: %s", req.RemoteAddr, req.URL)
}

func H500(w http.ResponseWriter, req *http.Request, rcv interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	etmpl.Execute(w, 500)
	log.Printf("500 :: %s :: %s :: %s", req.RemoteAddr, req.URL, rcv)
}

func tmplExec(w http.ResponseWriter, req *http.Request, subtmpl string, args interface{}) {
	c, err := tmpl.Clone()
	if err != nil {
		H500(w, req, err)
		return
	}

	if subtmpl != "" {
		_, err = c.New("i").Parse(subtmpl)
		if err != nil {
			H500(w, req, err)
			return
		}
	}
	err = c.Execute(w, args)
	if err != nil {
		log.Println(err)
	}
	log.Printf("200 :: %s :: %s", req.RemoteAddr, req.URL)
}

func HRoot(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", req.RemoteAddr, req.URL)
	tmplExec(w, req, "", nil)
}

func HAsset(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Printf("req :: %s :: %s", req.RemoteAddr, req.URL)
	var asset Asset
	if db.Find(&asset, "name = ? and kind = ?", ps.ByName("id"), ps.ByName("kind")).RecordNotFound() {
		H404(w, req)
		return
	}
	http.ServeContent(w, req, "", time.Time{}, bytes.NewReader(asset.Content))
	log.Printf("200 :: %s :: %s", req.RemoteAddr, req.URL)
}

func HPosts(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", req.RemoteAddr, req.URL)
	var posts []Post
	db.Select("title, short, created_at").Order("created_at desc").Find(&posts)
	tmplExec(w, req, r, posts)
}

func HPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Printf("req :: %s :: %s", req.RemoteAddr, req.URL)
	var post Post
	if db.First(&post, "short = ?", ps.ByName("id")).RecordNotFound() {
		H404(w, req)
		return
	}
	tmplExec(w, req, p, post)
}

func HProjects(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", req.RemoteAddr, req.URL)
	var projects []Project
	db.Order("id desc").Find(&projects)
	tmplExec(w, req, t, projects)
}

func main() {
	router := &httprouter.Router{
		NotFound:     H404,
		PanicHandler: H500,
	}
	router.GET("/", HRoot)
	router.GET("/assets/:kind/:id", HAsset)
	router.GET("/posts", HPosts)
	router.GET("/posts/:id", HPost)
	router.GET("/projects", HProjects)

	log.Fatal(http.ListenAndServe(":9000", router))
}
