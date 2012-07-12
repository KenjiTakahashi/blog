flatiron = require 'flatiron'
app = flatiron.app
director = require 'director'

app.use flatiron.plugins.http

routes =
    '/':
        get: ->
            @res.writeHead 200, {'Content-Type': 'text/plain'}
            @res.end "Hello World!"

app.router = new director.http.Router(routes).configure {async: true}

app.start 19689
