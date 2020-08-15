package main

import (
	"github.com/gorilla/feeds"

	"github.com/KenjiTakahashi/blog/db"
)

func getFeed() (*feeds.Feed, error) {
	posts, err := db.GetPosts(10)
	if err != nil {
		return nil, err
	}

	feed := &feeds.Feed{
		Title:       "kenji.sx",
		Link:        &feeds.Link{Href: "http://kenji.sx"},
		Description: "Karol Woźniak aka Kenji Takahashi :: place",
		Author:      &feeds.Author{"Karol Woźniak", "wozniakk@gmail.com"},
		Items:       make([]*feeds.Item, len(posts)),
	}
	for i, post := range posts {
		feed.Items[i] = &feeds.Item{
			Title:   post.Title,
			Link:    &feeds.Link{Href: "http://kenji.sx/posts/" + post.Short},
			Created: post.CreatedAt,
		}
	}
	return feed, nil
}
