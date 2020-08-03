package main

import (
	"github.com/gorilla/feeds"

	"github.com/KenjiTakahashi/blog/db"
)

var feed = &feeds.Feed{
	Title:       "kenji.sx",
	Link:        &feeds.Link{Href: "http://kenji.sx"},
	Description: "Karol Woźniak aka Kenji Takahashi :: place",
	Author:      &feeds.Author{"Karol Woźniak", "wozniakk@gmail.com"},
}

func feedFeed() error {
	posts, err := db.GetPosts(10)
	if err != nil {
		return err
	}

	feed.Items = make([]*feeds.Item, len(posts))
	for i, post := range posts {
		post := post.(*db.Post)
		feed.Items[i] = &feeds.Item{
			Title:   post.Title,
			Link:    &feeds.Link{Href: "http://kenji.sx/posts/" + post.Short},
			Created: post.CreatedAt,
		}
	}
	return nil
}
