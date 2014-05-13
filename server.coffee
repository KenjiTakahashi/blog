express = require 'express'
coffeecup = require 'coffeecup'
DBAssets = require './assets'
dispatcher = require './dispatcher'
assets = require 'connect-assets'
mdify = require './markdownify'

db = require './db'
posts = db.posts
images = db.images
projects = db.projects
raws = db.raws
Rss = require './rss'
rss = new Rss posts

datepicker = require './datepicker'

dbassets = new DBAssets db.assets

app = express()
app.engine 'coffee', coffeecup.__express
app.set 'view engine', 'coffee'
app.set 'views', "#{__dirname}/views"
app.use dispatcher
app.use assets buildDir: no
app.use dbassets.middleware

app.get '/', (req, res) ->
    res.render 'index'

app.get '/posts', (req, res) ->
    posts.titles res.dis res.html, (res, data) ->
        res.render 'index',
            items: data,
            date: datepicker.get_string,
            res.html

app.get '/posts/:short', (req, res) ->
    posts.one req.params.short, res.dis res.html, (res, data) ->
        data.content = mdify data.content
        res.render 'index', post: data, res.html

app.get '/images/:id', (req, res) ->
    images.one req.params.id, res.dis res.img, (res, data) ->
        res.img null, data.content, data.type

app.get '/raw/:id', (req, res) ->
    raws.one req.params.id, res.dis res.html, (res, data) ->
        res.html null, data.content

app.get '/projects', (req, res) ->
    projects.all res.dis res.html, (res, data) ->
        res.render 'index', items: data, res.html

app.get '/feed', (req, res) ->
    rss.generate res.xml

app.listen process.env.app_port || 8080, process.env.app_ip || "127.0.0.1"
