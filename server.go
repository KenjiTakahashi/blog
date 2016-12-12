package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/KenjiTakahashi/blog/db"
	"github.com/julienschmidt/httprouter"
)

type H404t struct{}

func (h H404t) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	etmpl.Execute(w, 404)
	log.Printf("404 :: %s :: %s", req.RemoteAddr, req.URL)
}

var H404 = H404t{}

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
	var asset db.Asset
	if db.DB.Find(&asset, "name = ? and kind = ?", ps.ByName("id"), ps.ByName("kind")).RecordNotFound() {
		H404.ServeHTTP(w, req)
		return
	}
	http.ServeContent(w, req, "", time.Time{}, bytes.NewReader(asset.Content))
	log.Printf("200 :: %s :: %s", req.RemoteAddr, req.URL)
}

func HPosts(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", req.RemoteAddr, req.URL)
	var posts []db.Post
	db.DB.Select("title, short, created_at").Order("created_at desc").Find(&posts)
	tmplExec(w, req, r, posts)
}

func HPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Printf("req :: %s :: %s", req.RemoteAddr, req.URL)
	var post db.Post
	if db.DB.First(&post, "short = ?", ps.ByName("id")).RecordNotFound() {
		H404.ServeHTTP(w, req)
		return
	}
	tmplExec(w, req, p, post)
}

func HProjects(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", req.RemoteAddr, req.URL)
	var projects []db.Project
	db.DB.Order("id desc").Find(&projects)
	tmplExec(w, req, t, projects)
}

func HFeedRss(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", req.RemoteAddr, req.URL)
	feedFeed()
	if err := feed.WriteRss(w); err != nil {
		log.Println(err)
	}
}

func HFeedAtom(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", req.RemoteAddr, req.URL)
	feedFeed()
	if err := feed.WriteAtom(w); err != nil {
		log.Println(err)
	}
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
	router.GET("/feed/rss", HFeedRss)
	router.GET("/feed/atom", HFeedAtom)

	log.Fatal(http.ListenAndServe(":9100", router))
}
