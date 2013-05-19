express = require 'express'
coffeecup = require 'coffeecup'
#Assets = require './utils/assets'
Mixpanel = require 'mixpanel'
mixpanel = Mixpanel.init "213db0f74e843fe533c2a173757c1d0a"
dispatcher = require './dispatcher.coffee'
assets = require 'connect-assets'
mdify = require './markdownify'

db = require './db'
posts = db.posts
projects = db.projects
raws = db.raws
Rss = require './rss'
rss = new Rss posts

datepicker = require './datepicker'

#assets = new Assets "#{__dirname}/assets", db.assets, run

app = express()
app.engine 'coffee', coffeecup.__express
app.set 'view engine', 'coffee'
app.set 'views', "#{__dirname}/views"
app.use dispatcher
app.use assets buildDir: no
#app.use assets.middleware

app.get '/', (req, res) ->
    res.render 'index'

app.get '/posts', (req, res) ->
    posts.titles (err, data) ->
        if err
            res.html err, null
        else
            res.render 'index',
                items: data,
                date: datepicker.get_string,
                res.html

app.get '/posts/:short', (req, res) ->
    posts.one req.params.short, (err, data) ->
        if err
            res.html err, null
        else
            data.content = mdify data.content
            res.render 'index', post: data, res.html

app.get '/projects', (req, res) ->
    projects.all (err, data) ->
        if err
            res.html err, null
        else
            res.render 'index', items: data, res.html

app.get '/feed', (req, res) ->
    rss.generate res.xml

app.listen process.env.app_port || 8080, process.env.app_ip || "127.0.0.1"
