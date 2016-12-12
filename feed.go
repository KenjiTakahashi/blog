package main

import (
	"github.com/KenjiTakahashi/blog/db"
	"github.com/gorilla/feeds"
)

var feed = &feeds.Feed{
	Title:       "kenji.sx",
	Link:        &feeds.Link{Href: "http://kenji.sx"},
	Description: "Karol Woźniak aka Kenji Takahashi :: place",
	Author:      &feeds.Author{"Karol Woźniak", "wozniakk@gmail.com"},
}

func feedFeed() {
	var posts []db.Post
	db.DB.Order("created_at desc").Limit(10).Find(&posts)

	feed.Items = make([]*feeds.Item, len(posts))
	for i, post := range posts {
		feed.Items[i] = &feeds.Item{
			Title:   post.Title,
			Link:    &feeds.Link{Href: "http://kenji.sx/posts/" + post.Short},
			Created: post.CreatedAt,
		}
	}
}
