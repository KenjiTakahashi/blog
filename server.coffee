flatiron = require 'flatiron'
director = require 'director'
connect = require 'connect'
datepicker = require './utils/datepicker'

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
    app.render 'index',
        month: datepicker.month(),
        year: datepicker.year(),
        previous: datepicker.previous(),
        current: datepicker.current(),
        next: datepicker.next(),
        (err, data) ->
            res.html err, data

routes =
    '/':
        get: placeholder

app.router = new director.http.Router(routes).configure async: true

app.start process.env.app_port || 8080
