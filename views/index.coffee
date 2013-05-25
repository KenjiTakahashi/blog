doctype 5
html ->
    head ->
        meta charset: 'utf-8'
        title 'Karol Woźniak aka Kenji Takahashi'
        text "#{css 'code'}#{css 'index'}"
    body ->
        div id: "o", ->
            div class: "m"
            div id: "i", ->
                div id: "head", ->
                    a href: "/posts", -> 'posts'
                    text ' :: '
                    a href: "/projects", -> 'projects'
                    text ' :: '
                    a href: "/", -> 'Karol Woźniak'
                    text ' aka '
                    a href: "/", -> 'Kenji Takahashi'
                if @items? and @items.length > 0
                    ul id: "items", ->
                        for item in @items
                            li ->
                                if @date?
                                    text "#{@date item.date} "
                                    a href: "/posts/#{item.short}", ->
                                        text "#{item.title}"
                                else
                                    if item.active
                                        a href: "#{item.site}", ->
                                            text "#{item.name}"
                                    else
                                        a
                                            href: "#{item.site}",
                                            class: "dead",
                                            -> "#{item.name}"
                                    text " #{item.desc}"
                else if @post?
                    div id: "post", ->
                        p id: "title", -> "#{@post.title}"
                        text "#{@post.content}"
            div class: "m"
        div id: "foot", ->
            a target: "_blank", href: "https://github.com/KenjiTakahashi", -> 'github'
            text ' :: '
            a href:"/feed", -> 'rss'
            text ' :: Kenji Takahashi © 2013'
        text "#{js 'anal'}"
