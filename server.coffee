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

placeholder = (req, res) ->
    if not res?
        res = @res
    posts.latest (err, latest) ->
        latest.content = marked latest.content
        latest.tags = ("<a href='/tags/#{t}'>#{t}</a>" for t in latest.tags).join(', ')
        posts.tags (err, tags) ->
            posts.titles (err, titles) ->
                projects.all (err, projects) ->
                    app.render 'index',
                        month: datepicker.month(),
                        year: datepicker.year(),
                        previous: datepicker.previous(),
                        current: datepicker.current(),
                        next: datepicker.next(),
                        latest: latest,
                        titles: titles,
                        projects: projects,
                        tags: tags,
                        (err, data) ->
                            res.html err, data

routes =
    '/':
        get: placeholder

app.router = new director.http.Router(routes).configure async: true

app.start process.env.app_port || 8080
