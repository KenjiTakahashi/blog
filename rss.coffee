Rss = require 'juan-rss'

class RSS
    constructor: (@db) ->
        @generator = new Rss {
            title: 'kenji.sx'
            description: 'Karol "Kenji Takahashi" Woźniak\'s personal website'
            feed_url: 'http://kenji.sx/feed'
            site_url: 'http://kenji.sx'
            author: 'Karol "Kenji Takahashi" Woźniak'
        }

    generate: (callback) ->
        @db.last20 (err, data) =>
            if err or not data
                callback err, null
            else
                for item in data
                    @generator.item {
                        title: item.title
                        url: "http://kenji.sx/posts/#{item.short}"
                        date: item.date
                    }
                callback null, @generator.xml()

module.exports = RSS
