express = require 'express'
director = require 'director'
connect = require 'connect'
assets = require 'connect-assets'
jsprimer = require 'connect-assets-jsprimer'

mdify = require './utils/markdownify'

db = require './db'
posts = db.posts
projects = db.projects
raws = db.raws

Datepicker = require './utils/datepicker'
datepicker = new Datepicker posts

app = express()
app.use require './utils/dispatcher'
app.use connect.static "#{__dirname}/public"
app.use assets src: "#{__dirname}/assets", buildDir: false
jsprimer.loadFiles assets

placeholder = (req, res, post) ->
    query = req.query
    tag = query.t
    if tag?
        tag = parseInt tag
    page = query.p
    if page?
        page = parseInt page
    url = req.url.split '?'
    urls = [url[0], "", "", "", "", "", url[1] and "?#{url[1]}" or ""]
    for n, v of query
        if n != 't'
            urls[1] += "&#{n}=#{v}"
        if n != 'p'
            urls[2] += "&#{n}=#{v}"
        if n != 'tag' and n != 'date'
            urls[3] += "&#{n}=#{v}"
        if n != 'd'
            urls[4] += "&#{n}=#{v}"
        if n != 'date' and n != 'tag'
            urls[5] += "&#{n}=#{v}"
    post.content = mdify post.content
    post.month = datepicker.get_month_name post.date.getMonth()
    post.tags = ("<a href='#{urls[0]}?tag=#{t}#{urls[3]}'>#{t}</a>" for t in post.tags).join(', ')
    posts.tags tag, (err, tags, has_prev_tag, has_next_tag) ->
        prev_tag = [has_prev_tag, tag? and tag - 1]
        next_tag = [has_next_tag, tag? and tag + 1 or 1]
        if err
            res.html err, null
        else if tags == null
            res.html 'INTERNAL', null
        else
            posts.titles query.date, query.tag, page,
            (err, titles, has_prev_page, has_next_page) ->
                if err
                    res.html err, null
                else if titles == null
                    res.html 'INTERNAL', null
                else
                    prev_title = [has_prev_page, page? and page - 1]
                    next_title = [has_next_page, page? and page + 1 or 1]
                    projects.all (err, projects) ->
                        datepicker.get query.d, (err, dp) ->
                            if query.date?
                                dp.selected = query.date
                            app.render 'index.jade',
                                urls: urls,
                                dp: dp,
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
            posts.latest (err, data) =>
                placeholder @req, @res, data
    '/posts/:id':
        get: (id) ->
            posts.one id, (err, data) =>
                placeholder @req, @res, data
    '/raw/:id':
        get: (id) ->
            raws.one id, @res.html
router = new director.http.Router(routes).configure async: true
app.use (req, res, next) ->
    router.dispatch req, res, (err) ->
        if err == undefined || err
            next()

app.listen process.env.app_port || 8080
