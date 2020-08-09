package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/KenjiTakahashi/blog/db"
)

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func ShiftPath(p string) (string, string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func getRA(req *http.Request) string {
	ip := req.Header.Get("X-Real-IP")
	if ip == "" {
		ip = req.RemoteAddr
	}
	return fmt.Sprintf("%-21s", ip)
}

var H404 = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	etmpl.Execute(w, 404)
	// FIXME: This URL might be "truncated" already
	log.Printf("404 :: %s :: %s", getRA(req), req.URL)
})

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

var HProjects = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	projects, err := db.GetProjects()
	if err != nil {
		H500(rw, req, err)
		return
	}
	tmplExec(rw, req, t, projects)
})

var HPost = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	var id string
	// FIXME: Check that there is no more URL
	id, req.URL.Path = ShiftPath(req.URL.Path)

	postStr, err := db.Get("post:science/%s", id)
	if err == db.ErrNotFound {
		H404.ServeHTTP(rw, req)
		return
	}
	if err != nil {
		H500(rw, req, err)
		return
	}
	var post db.Post
	err = json.Unmarshal([]byte(postStr), &post)
	if err != nil {
		H500(rw, req, err)
		return
	}
	tmplExec(rw, req, p, post)
})

var HPosts = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	head, _ := ShiftPath(req.URL.Path)

	if head == "" {
		posts, err := db.GetPosts(0)
		if err != nil {
			H500(rw, req, err)
			return
		}
		tmplExec(rw, req, r, posts)
		return
	}

	HPost.ServeHTTP(rw, req)
})

var HFeedRss = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	feedFeed()
	if err := feed.WriteRss(rw); err != nil {
		log.Println(err)
	}
})

var HFeedAtom = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	feedFeed()
	if err := feed.WriteAtom(rw); err != nil {
		log.Println(err)
	}
})

var HFeed = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "atom":
		HFeedAtom.ServeHTTP(rw, req)
	case "rss":
		HFeedRss.ServeHTTP(rw, req)
	default:
		H404.ServeHTTP(rw, req)
	}
})

var HAssets = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	var kind string
	kind, req.URL.Path = ShiftPath(req.URL.Path)
	var id string
	id, req.URL.Path = ShiftPath(req.URL.Path)

	asset, err := db.Get("asset:%s:%s", kind, id)
	if err == db.ErrNotFound {
		H404.ServeHTTP(rw, req)
		return
	}
	if err != nil {
		H500(rw, req, err)
		return
	}
	http.ServeContent(rw, req, "", time.Time{}, strings.NewReader(asset))

	log.Printf("200 :: %s :: %s", getRA(req), req.URL)
})

var HRoot = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	log.Printf("req :: %s :: %s", getRA(req), req.URL)

	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	log.Println("!!!!!", head)
	log.Println("!!!!!", req.URL.Path)

	if head == "" {
		tmplExec(rw, req, "", nil)
		return
	}
	switch head {
	case "assets":
		HAssets.ServeHTTP(rw, req)
	case "feed":
		HFeed.ServeHTTP(rw, req)
	case "posts":
		HPosts.ServeHTTP(rw, req)
	case "projects":
		HProjects.ServeHTTP(rw, req)
	default:
		H404.ServeHTTP(rw, req)
	}
})

func main() {
	log.Fatal(http.ListenAndServe(":9100", HRoot))
}
