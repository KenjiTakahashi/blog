package main

import (
	"context"
	"encoding/json"
	"errors"
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
func PathSplit1String(p string) (string, string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func PathSplit1(req *http.Request) (string, string) {
	return PathSplit1String(req.URL.Path)
}

func PathShift(req *http.Request) string {
	var head string
	head, req.URL.Path = PathSplit1(req)
	return head
}

func PathShiftChecked(req *http.Request) (string, error) {
	head, tail := PathSplit1(req)
	if tail != "/" {
		return "", errors.New("AAAA")
	}
	req.URL.Path = tail
	return head, nil
}

func SetStatusCode(rw http.ResponseWriter, req *http.Request, code int) {
	*req = *req.WithContext(context.WithValue(req.Context(), "StatusCode", code))
	rw.WriteHeader(code)
}

func SetStatusError(rw http.ResponseWriter, req *http.Request, err error) {
	*req = *req.WithContext(context.WithValue(req.Context(), "Error", err))
	SetStatusCode(rw, req, http.StatusInternalServerError)
}

var H404 = func(rw http.ResponseWriter) error {
	return etmpl.Execute(rw, http.StatusNotFound)
}

var H500 = func(rw http.ResponseWriter) error {
	return etmpl.Execute(rw, http.StatusInternalServerError)
}

func tmplExec(rw http.ResponseWriter, req *http.Request, subtmpl string, args interface{}) {
	c, err := tmpl.Clone()
	if err != nil {
		SetStatusError(rw, req, err)
		return
	}

	if subtmpl != "" {
		_, err = c.New("i").Parse(subtmpl)
		if err != nil {
			SetStatusError(rw, req, err)
			return
		}
	}
	if err = c.Execute(rw, args); err != nil {
		SetStatusError(rw, req, err)
	}
}

var HProjects = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	projects, err := db.GetProjects()
	if err != nil {
		SetStatusError(rw, req, err)
		return
	}
	tmplExec(rw, req, t, projects)
})

var HPost = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	id, err := PathShiftChecked(req)
	if err != nil {
		return
	}

	postStr, err := db.Get("post:science/%s", id)
	if err == db.ErrNotFound {
		SetStatusCode(rw, req, http.StatusNotFound)
		return
	}
	if err != nil {
		SetStatusError(rw, req, err)
		return
	}
	var post db.Post
	err = json.Unmarshal([]byte(postStr), &post)
	if err != nil {
		SetStatusError(rw, req, err)
		return
	}
	tmplExec(rw, req, p, post)
})

var HPosts = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	head, _ := PathSplit1(req)

	if head == "" {
		posts, err := db.GetPosts(0)
		if err != nil {
			SetStatusError(rw, req, err)
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
	head, err := PathShiftChecked(req)
	if err != nil {
		return
	}

	switch head {
	case "atom":
		HFeedAtom.ServeHTTP(rw, req)
	case "rss":
		HFeedRss.ServeHTTP(rw, req)
	default:
		SetStatusCode(rw, req, http.StatusNotFound)
	}
})

var HAssets = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	kind := PathShift(req)
	id, err := PathShiftChecked(req)
	if err != nil {
		return
	}

	asset, err := db.Get("asset:%s:%s", kind, id)
	if err == db.ErrNotFound {
		SetStatusCode(rw, req, http.StatusNotFound)
		return
	}
	if err != nil {
		SetStatusError(rw, req, err)
		return
	}
	http.ServeContent(rw, req, "", time.Time{}, strings.NewReader(asset))
})

var HRoot = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	head := PathShift(req)

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
		SetStatusCode(rw, req, http.StatusNotFound)
	}
})

func NewMLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ip := req.Header.Get("X-Real-IP")
		if ip == "" {
			ip = req.RemoteAddr
		}
		fullURL := req.URL.Path

		log.Printf("req :: %-21s :: %s", ip, fullURL)

		next.ServeHTTP(rw, req)

		statusCode := req.Context().Value("StatusCode")
		if statusCode == nil {
			statusCode = http.StatusOK
		}

		format := "%03d :: %-21s :: %s"
		args := []interface{}{statusCode, ip, fullURL}
		err := req.Context().Value("Error")
		if err != nil {
			format += " :: ERROR : %s"
			args = append(args, err)
		}
		log.Printf(format, args...)
	})
}

func NewErrH(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(rw, req)

		statusCode := req.Context().Value("StatusCode")
		// TODO: Allow handlers being http.Handler?
		switch statusCode {
		case http.StatusNotFound:
			H404(rw)
		case http.StatusInternalServerError:
			H500(rw)
		}
	})
}

func main() {
	log.Fatal(http.ListenAndServe(":9100", NewErrH(NewMLog(HRoot))))
}
