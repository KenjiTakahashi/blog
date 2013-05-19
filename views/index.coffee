doctype 5
html ->
    head ->
        meta charset: 'utf-8'
        title 'Karol Woźniak aka Kenji Takahashi'
        text "#{css 'code'}#{css 'index'}"
    body ->
        div id: "head", ->
            a href: "/posts", -> 'blog'
            text ' :: '
            a href: "/projects", -> 'projects'
            text ' :: '
            a href: "/", -> 'Karol Woźniak'
            text ' aka '
            a href: "/", -> 'Kenji Takahashi'
            text ' :: '
            a href:"/feed", -> 'rss'
        if @items? and @items.length > 0
            ul id: "items", ->
                for item in @items
                    li ->
                        if @date?
                            text "#{@date item.date} "
                            a href: "/posts/#{item.short}", -> "#{item.title}"
                        else
                            if item.active
                                a href: "#{item.site}", -> "#{item.name}"
                            else
                                a href: "#{item.site}", class: "dead", -> "#{item.name}"
                            text " #{item.desc}"
        else if @post?
            div id: "post", ->
                p id: "title", -> "#{@post.title}"
                text "#{@post.content}"
        div id: "foot", -> 'Kenji Takahashi © 2013'
        text "#{js 'app'}"
