flatiron = require 'flatiron'
director = require 'director'
connect = require 'connect'
datepicker = require './utils/datepicker'

hljs = require 'highlight.js'
hljs.tabReplace = '    '
marked = require 'marked'
marked.setOptions gfm: true, highlight: (code, lang) ->
    return hljs.highlight(lang, code).value

db = require './db'
posts = db.posts
projects = db.projects

app = flatiron.app
app.use flatiron.plugins.http, before: [
    (req, res) ->
        res.redirect = (path) ->
            res.writeHead 302, 'Location': path
            res.end()
        res.emit 'next'
    require('stylus').middleware
        src: "#{__dirname}/assets"
        dest: "#{__dirname}/public"
    require './utils/dispatcher'
    connect.static "#{__dirname}/public"
]
app.use require './utils/jade'

placeholder = (req, res, post) ->
    query = req.query
    tag = query.t
    if tag?
        tag = parseInt tag
    page = query.p
    if page?
        page = parseInt page
    post.content = marked post.content
    post.tags = ("<a href='/tags/#{t}'>#{t}</a>" for t in post.tags).join(', ')
    posts.tags tag, (err, tags, has_prev_tag, has_next_tag) ->
        prev_tag = [has_prev_tag, tag? and tag - 1]
        next_tag = [has_next_tag, tag? and tag + 1 or 1]
        if err
            res.html err, null
        else if tags == null
            res._html.internal null
        else
            posts.titles page, (err, titles, has_prev_page, has_next_page) ->
                if err
                    res.html err, null
                else if titles == null
                    res._html.internal null
                else
                    prev_title = [has_prev_page, page? and page - 1]
                    next_title = [has_next_page, page? and page + 1 or 1]
                    projects.all (err, projects) ->
                        app.render 'index',
                            url: req.url.split('?', 1)[0],
                            month: datepicker.month(),
                            year: datepicker.year(),
                            previous: datepicker.previous(),
                            current: datepicker.current(),
                            next: datepicker.next(),
                            latest: post,
                            titles: titles,
                            prev_title: prev_title,
                            next_title: next_title,
                            projects: projects,
                            tags: tags,
                            prev_tag: prev_tag,
                            next_tag: next_tag,
                            (err, data) ->
                                res.html err, data

routes =
    '/':
        get: ->
            [req, res] = [@req, @res]
            posts.latest (err, data) ->
                placeholder req, res, data
    '/posts/:id':
        get: (id) ->
            [req, res] = [@req, @res]
            posts.one id, (err, data) ->
                placeholder req, res, data

app.router = new director.http.Router(routes).configure async: true

app.start process.env.app_port || 8080
