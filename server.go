package main

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/KenjiTakahashi/blog/db"
	"github.com/rs/xid"
)

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
		return "", errors.New("Unexpected Path Ending")
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

func Log(req *http.Request, msg string, args ...interface{}) {
	reqID := req.Context().Value("ReqID")
	if reqID == nil {
		reqID = "--------------------"
	}
	statusCode := req.Context().Value("StatusCode")
	if statusCode == nil {
		statusCode = http.StatusOK
	}

	format := ":: %s :: %03d"
	if msg != "" {
		format += " :: " + msg
	}
	log.Printf(format, append([]interface{}{reqID, statusCode}, args...)...)
}

func LogIfErr(req *http.Request, err error) {
	if err != nil {
		Log(req, "POST SEND ERROR : %s", err)
	}
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
	var buf bytes.Buffer
	if err = c.Execute(&buf, args); err != nil {
		SetStatusError(rw, req, err)
		return
	}
	_, err = buf.WriteTo(rw)
	LogIfErr(req, err)
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

	post, err := db.GetPost("science", id)
	if err == db.ErrNotFound {
		SetStatusCode(rw, req, http.StatusNotFound)
		return
	}
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
	if feed, err := getFeed(); err != nil {
		SetStatusError(rw, req, err)
	} else {
		LogIfErr(req, feed.WriteRss(rw))
	}
})

var HFeedAtom = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	if feed, err := getFeed(); err != nil {
		SetStatusError(rw, req, err)
	} else {
		LogIfErr(req, feed.WriteAtom(rw))
	}
})

var HFeed = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	head, err := PathShiftChecked(req)
	if err != nil {
		SetStatusCode(rw, req, http.StatusNotFound)
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

var HAsset = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	kind := PathShift(req)
	id, err := PathShiftChecked(req)
	if err != nil {
		SetStatusCode(rw, req, http.StatusNotFound)
		return
	}

	asset, err := db.GetAsset(kind, id)
	if err == db.ErrNotFound {
		SetStatusCode(rw, req, http.StatusNotFound)
		return
	}
	if err != nil {
		SetStatusError(rw, req, err)
		return
	}
	http.ServeContent(rw, req, "", time.Time{}, bytes.NewReader(asset))
})

var HRoot = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	head := PathShift(req)

	if head == "" {
		tmplExec(rw, req, "", nil)
		return
	}
	switch head {
	case "assets":
		HAsset.ServeHTTP(rw, req)
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
		reqID := xid.New().String()
		fullURL := req.URL.Path

		*req = *req.WithContext(context.WithValue(req.Context(), "ReqID", reqID))

		log.Printf(":: %s :: req :: %s", reqID, fullURL)

		next.ServeHTTP(rw, req)

		statusCode := req.Context().Value("StatusCode")
		if statusCode == nil {
			statusCode = http.StatusOK
		}

		err := req.Context().Value("Error")
		if err != nil {
			Log(req, "ERROR : %s", err)
		} else {
			Log(req, "")
		}
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
