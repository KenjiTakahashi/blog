package db

import (
	"sync"
)

type Posts struct {
	List []*Post
	Map  map[string]*Post
}

func (p *Posts) Set(key string, post *Post) {
	post.Short = key
	p.List = append(p.List, post)
	p.Map[key] = post
}

func (p *Posts) GetMany(n int) []*Post {
	if n == 0 {
		return p.List
	}
	if n < 0 {
		n = 0
	}
	if n > len(p.List) {
		n = len(p.List)
	}
	return p.List[:n]
}

func (p *Posts) GetOne(id string) (*Post, error) {
	post, exists := p.Map[id]
	if !exists {
		return nil, ErrNotFound
	}
	return post, nil
}

type cache struct {
	Assets   map[string]map[string][]byte
	Posts    Posts
	Projects []Project
}

func newCache() cache {
	return cache{
		Assets: map[string]map[string][]byte{
			"image":  {},
			"raw":    {},
			"script": {},
		},
		Posts: Posts{
			Map: map[string]*Post{},
		},
	}
}

var defcLock sync.RWMutex
var defc = newCache()
