package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/KenjiTakahashi/blog/db"
)

func getRA(req *http.Request) string {
	ip := req.Header.Get("X-Real-IP")
	if ip == "" {
		ip = req.RemoteAddr
	}
	return fmt.Sprintf("%-21s", ip)
}

type H404t struct{}

func (h H404t) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	etmpl.Execute(w, 404)
	log.Printf("404 :: %s :: %s", getRA(req), req.URL)
}

var H404 = H404t{}

func H500(w http.ResponseWriter, req *http.Request, rcv interface{}) {
	w.WriteHeader(http.StatusInternalServerError)
	etmpl.Execute(w, 500)
	log.Printf("500 :: %s :: %s :: %s", getRA(req), req.URL, rcv)
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
	log.Printf("200 :: %s :: %s", getRA(req), req.URL)
}

func HRoot(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", getRA(req), req.URL)
	tmplExec(w, req, "", nil)
}

func HAsset(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Printf("req :: %s :: %s", getRA(req), req.URL)
	asset, err := db.Get("asset:%s:%s", ps.ByName("id"), ps.ByName("kind"))
	if err != nil {
		log.Println(err, ps)
	}
	http.ServeContent(w, req, "", time.Time{}, strings.NewReader(asset))
	log.Printf("200 :: %s :: %s", getRA(req), req.URL)
}

func HPosts(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", getRA(req), req.URL)
	posts, err := db.GetPosts(0)
	if err != nil {
		log.Println(err)
	}
	tmplExec(w, req, r, posts)
}

func HPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Printf("req :: %s :: %s", getRA(req), req.URL)
	postStr, err := db.Get("post:science/%s", ps.ByName("id"))
	if err == db.ErrNotFound {
		H404.ServeHTTP(w, req)
		return
	}
	if err != nil {
		log.Println(err)
	}
	var post db.Post
	err = json.Unmarshal([]byte(postStr), &post)
	if err != nil {
		log.Println(err)
	}
	tmplExec(w, req, p, post)
}

func HProjects(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", getRA(req), req.URL)
	projects, err := db.GetProjects()
	if err != nil {
		log.Println(err)
	}
	tmplExec(w, req, t, projects)
}

func HFeedRss(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", getRA(req), req.URL)
	feedFeed()
	if err := feed.WriteRss(w); err != nil {
		log.Println(err)
	}
}

func HFeedAtom(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	log.Printf("req :: %s :: %s", getRA(req), req.URL)
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
